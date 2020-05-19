package pipeline

import (
	"context"
	"errors"
	"github.com/hashicorp/go-multierror"
	"sync"
)

const (
	// Context key used by StatsRecorder
	CtxStats = "stats"

	// Setup stage (Chore): performs setup tasks
	StageSetup StageName = "stage_setup"

	// Fetcher stage (Syncing): fetches data for indexing
	StageFetcher StageName = "stage_fetcher"

	// Parser stage (Syncing): parses and normalizes fetched data to a single structure
	StageParser StageName = "stage_parser"

	// Validator stage (Syncing): validates parsed data
	StageValidator StageName = "stage_validator"

	// Syncer stage (Syncing): saves data to datastore
	StageSyncer StageName = "stage_syncer"

	// Sequencer stage (Indexing): Creates sequences from synced data (syncable)
	StageSequencer StageName = "stage_sequencer"

	// Aggregator stage (Indexing): Creates aggregates from synced data (syncable)
	StageAggregator StageName = "stage_aggregator"

	// Cleanup stage (Chore): Cleans up after execution
	StageCleanup StageName = "stage_cleanup"
)

type StageName string

type Options struct {
	// StagesWhitelist holds name of stages that will be executed
	StagesWhitelist []StageName

	// StagesBlacklist holds name of stages that will NOT be executed
	StagesBlacklist []StageName

	// IndexingTasksWhitelist holds name of indexing tasks which will be executed
	IndexingTasksWhitelist []string

	// IndexingTasksBlacklist holds name of indexing tasks which will NOT be executed
	IndexingTasksBlacklist []string
}

// Pipeline implements a modular, multi-stage pipeline
type Pipeline struct {
	source         Source
	sink           Sink
	payloadFactory PayloadFactory

	options *Options

	stages map[StageName]*stage

	beforeStage map[StageName][]*stage
	afterStage  map[StageName][]*stage
}

func New(source Source, sink Sink, payloadFactor PayloadFactory) *Pipeline {
	return &Pipeline{
		source:         source,
		sink:           sink,
		payloadFactory: payloadFactor,
		stages:         make(map[StageName]*stage),
		beforeStage:    make(map[StageName][]*stage),
		afterStage:     make(map[StageName][]*stage),
	}
}

// SetOptions sets pipeline options
func (p *Pipeline) SetOptions(o *Options) {
	p.options = o
}

// SetSetupStage add setup stage to list of available stages
func (p *Pipeline) SetSetupStage(stageRunner StageRunner) {
	p.stages[StageSetup] = NewStage(StageSetup, stageRunner, StageTypeChore)
}

// SetFetcherStage add fetcher stage to list of available stages
func (p *Pipeline) SetFetcherStage(stageRunner StageRunner) {
	p.stages[StageFetcher] = NewStage(StageFetcher, stageRunner, StageTypeSyncing)
}

// SetParserStage add parser stage to list of available stages
func (p *Pipeline) SetParserStage(stageRunner StageRunner) {
	p.stages[StageParser] = NewStage(StageParser, stageRunner, StageTypeSyncing)
}

// SetValidatorStage add validator stage to list of available stages
func (p *Pipeline) SetValidatorStage(stageRunner StageRunner) {
	p.stages[StageValidator] = NewStage(StageValidator, stageRunner, StageTypeSyncing)
}

// SetSyncerStage add syncer stage to list of available stages
func (p *Pipeline) SetSyncerStage(stageRunner StageRunner) {
	p.stages[StageSyncer] = NewStage(StageSyncer, stageRunner, StageTypeSyncing)
}

// SetSequencerStage add sequencer stage to list of available stages
func (p *Pipeline) SetSequencerStage(stageRunner StageRunner) {
	p.stages[StageSequencer] = NewStage(StageSequencer, stageRunner, StageTypeIndexing)
}

// SetAggregatorStage add aggregator stage to list of available stages
func (p *Pipeline) SetAggregatorStage(stageRunner StageRunner) {
	p.stages[StageAggregator] = NewStage(StageAggregator, stageRunner, StageTypeIndexing)
}

// SetCleanupStage add cleanup stage to list of available stages
func (p *Pipeline) SetCleanupStage(stageRunner StageRunner) {
	p.stages[StageCleanup] = NewStage(StageCleanup, stageRunner, StageTypeChore)
}

// AddStageBefore adds custom stage before existing stage
func (p *Pipeline) AddStageBefore(existingStageName StageName, name StageName, stageRunner StageRunner) {
	p.beforeStage[existingStageName] = append(p.beforeStage[existingStageName], NewStage(name, stageRunner, StageTypeCustom))
}

// AddStageBefore adds custom stage after existing stage
func (p *Pipeline) AddStageAfter(existingStageName StageName, name StageName, stageRunner StageRunner) {
	p.afterStage[existingStageName] = append(p.afterStage[existingStageName], NewStage(name, stageRunner, StageTypeCustom))
}

// Start starts the pipeline
func (p *Pipeline) Start(ctx context.Context) error {
	pCtx, _ := p.setupCtx(ctx)

	// Source is responsible for getting start and end heights
	// If it fails we cannot proceed with execution of pipeline
	if err := p.source.Run(pCtx); err != nil {
		return err
	}

	startHeight, endHeight, err := p.getInterval()
	if err != nil {
		return err
	}

	var stagesErr error
	for i := *startHeight; i <= *endHeight; i++ {
		payload := p.payloadFactory.GetPayload()

		payload.SetCurrentHeight(i)

		stagesErr = p.runStages(pCtx, payload)
		if stagesErr != nil {
			// We don't want to run pipeline for rest of heights since we don't want to have gaps in records
			break
		}

		payload.MarkAsProcessed()
	}

	if err := p.sink.Run(pCtx); err != nil {
		return multierror.Append(stagesErr, err)
	}

	return stagesErr
}

// setupCtx sets up the context
func (p *Pipeline) setupCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	// Setup cancel
	pCtx, cancelFunc := context.WithCancel(ctx)

	// Setup stats recorder
	statRecorder := NewStatsRecorder()
	pCtx = context.WithValue(pCtx, CtxStats, statRecorder)

	return pCtx, cancelFunc
}

// getInterval gets start and end heights from the source
func (p *Pipeline) getInterval() (*int64, *int64, error) {
	startHeight := p.source.GetStartHeight()
	endHeight := p.source.GetEndHeight()

	if startHeight > endHeight {
		return nil, nil, errors.New("start height has to be smaller than end height")
	}
	return &startHeight, &endHeight, nil
}

// runStages runs all the stages
func (p *Pipeline) runStages(ctx context.Context, payload Payload) error {
	if err := p.runStage(StageSetup, ctx, payload); err != nil {
		return err
	}

	if err := p.runSyncingStages(ctx, payload); err != nil {
		return err
	}

	if err := p.runIndexingStages(ctx, payload); err != nil {
		return err
	}

	if err := p.runStage(StageCleanup, ctx, payload); err != nil {
		return err
	}
	return nil
}

// runSyncingStages runs syncing stages in sequence
func (p *Pipeline) runSyncingStages(ctx context.Context, payload Payload) error {
	if err := p.runStage(StageFetcher, ctx, payload); err != nil {
		return err
	}

	if err := p.runStage(StageParser, ctx, payload); err != nil {
		return err
	}

	if err := p.runStage(StageValidator, ctx, payload); err != nil {
		return err
	}

	if err := p.runStage(StageSyncer, ctx, payload); err != nil {
		return err
	}
	return nil
}

// runIndexingStages runs indexing stages concurrently
func (p *Pipeline) runIndexingStages(ctx context.Context, payload Payload) error {
	var errs error
	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	wg.Add(2)

	go func() {
		if err := p.runStage(StageSequencer, ctx, payload); err != nil {
			errCh <- err
		}
		wg.Done()
	}()

	go func() {
		if err := p.runStage(StageAggregator, ctx, payload); err != nil {
			errCh <- err
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		errs = multierror.Append(errs, err)
	}
	return errs
}

// runStage executes stage runner for given stage
func (p *Pipeline) runStage(stageName StageName, ctx context.Context, payload Payload) error {
	if p.canRunStage(stageName) {
		before := p.beforeStage[stageName]
		if len(before) > 0 {
			for _, s := range before {
				if err := s.Run(ctx, payload, p.options); err != nil {
					return err
				}
			}
		}

		if err := p.stages[stageName].Run(ctx, payload, p.options); err != nil {
			return err
		}

		after := p.afterStage[stageName]
		if len(after) > 0 {
			for _, s := range after {
				if err := s.Run(ctx, payload, p.options); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// canRunStage determines if stage can be ran
func (p *Pipeline) canRunStage(stageName StageName) bool {
	_, ok := p.stages[stageName]
	if !ok {
		return false
	}

	if p.options != nil && len(p.options.StagesWhitelist) > 0 {
		for _, s := range p.options.StagesWhitelist {
			if s == stageName {
				return true
			}
		}
		return false
	}

	if p.options != nil && len(p.options.StagesBlacklist) > 0 {
		for _, s := range p.options.StagesBlacklist {
			if s == stageName {
				return false
			}
		}
		return true
	}

	return true
}
