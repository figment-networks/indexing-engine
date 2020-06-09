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

var (
	ErrVersionsDirNotSet = errors.New("versions directory not set")
)

type StageName string

type Options struct {
	// TaskWhitelist holds name of indexing tasks which will be executed
	TaskWhitelist VersionTasks
}

// Pipeline implements a modular, multi-stage pipeline
type Pipeline struct {
	payloadFactory PayloadFactory
	options        *Options

	stages map[StageName]*stage

	beforeStage map[StageName][]*stage
	afterStage  map[StageName][]*stage
}

func New(payloadFactor PayloadFactory) *Pipeline {
	return &Pipeline{
		payloadFactory: payloadFactor,

		stages:      make(map[StageName]*stage),
		beforeStage: make(map[StageName][]*stage),
		afterStage:  make(map[StageName][]*stage),
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
func (p *Pipeline) Start(ctx context.Context, source Source, sink Sink, options *Options) error {
	pCtx, _ := p.setupCtx(ctx)

	p.options = options

	var pipelineErr error
	var recentPayload Payload
	for ok := true; ok; ok = source.Next(ctx, recentPayload) {
		payload := p.payloadFactory.GetPayload(source.Current())

		pipelineErr = p.runStages(pCtx, payload)
		if pipelineErr != nil {
			// We don't want to run pipeline for rest of heights since we don't want to have gaps in records
			break
		}

		if err := sink.Consume(pCtx, payload); err != nil {
			pipelineErr = err
			// Stop execution when sink errors out
			break
		}

		payload.MarkAsProcessed()

		recentPayload = payload
	}

	if err := source.Err(); err != nil {
		pipelineErr = multierror.Append(pipelineErr, err)
	}

	return pipelineErr
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
	if err := p.runStage(StageSyncer, ctx, payload); err != nil {
		return err
	}

	if err := p.runStage(StageFetcher, ctx, payload); err != nil {
		return err
	}

	if err := p.runStage(StageParser, ctx, payload); err != nil {
		return err
	}

	if err := p.runStage(StageValidator, ctx, payload); err != nil {
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
	return true
}
