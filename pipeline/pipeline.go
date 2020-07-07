package pipeline

import (
	"context"
	"errors"
	"fmt"
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

	// StagePersistor stage (Indexing): Persists data to datastore
	StagePersistor StageName = "stage_persistor"

	// Cleanup stage (Chore): Cleans up after execution
	StageCleanup StageName = "stage_cleanup"
)

var (
	ErrMissingStages           = errors.New("provide stages to run concurrently")
	ErrMisconfiguredDependency = errors.New("misconfigured stage dependency")

	defaultDependencies = map[StageName][]StageName{
		StageSetup:      []StageName{},
		StageSyncer:     []StageName{StageSetup},
		StageFetcher:    []StageName{StageSyncer},
		StageParser:     []StageName{StageFetcher},
		StageValidator:  []StageName{StageParser},
		StageSequencer:  []StageName{StageValidator},
		StageAggregator: []StageName{StageValidator},
		StagePersistor:  []StageName{StageAggregator, StageSequencer},
		StageCleanup:    []StageName{StagePersistor},
	}
)

type StageName string
type TaskName string

type Options struct {
	// StagesBlacklist holds list of stages to turn off
	StagesBlacklist []StageName

	// TaskWhitelist holds name of indexing tasks which will be executed
	TaskWhitelist []TaskName
}

// Pipeline implements a modular, multi-stage pipeline
type Pipeline struct {
	payloadFactory PayloadFactory
	options        *Options

	stages map[StageName]*stage

	stageDependencies map[StageName][]StageName
	beforeStage       map[StageName][]*stage
	afterStage        map[StageName][]*stage
}

// New creates a pipeline
func New(payloadFactor PayloadFactory) *Pipeline {
	return &Pipeline{
		payloadFactory: payloadFactor,

		stages:            make(map[StageName]*stage),
		stageDependencies: defaultDependencies,
		beforeStage:       make(map[StageName][]*stage),
		afterStage:        make(map[StageName][]*stage),
	}
}

func NewCustom(payloadFactor PayloadFactory) *Pipeline {
	return &Pipeline{
		payloadFactory: payloadFactor,

		stages:            make(map[StageName]*stage),
		stageDependencies: make(map[StageName][]StageName),
		beforeStage:       make(map[StageName][]*stage),
		afterStage:        make(map[StageName][]*stage),
	}
}

// SetOptions sets pipeline options
func (p *Pipeline) SetOptions(o *Options) {
	p.options = o
}

// SetLogger sets logger
func (p *Pipeline) SetLogger(l Logger) {
	logger = l
}

// SetStage sets up stage runner for given stage
func (p *Pipeline) SetStage(stageName StageName, stageRunner StageRunner) {
	p.stages[stageName] = NewStage(stageName, stageRunner)
}

// SetDependency sets dependencies that must run (if configured) before stage can be executed
func (p *Pipeline) SetDependency(stageName StageName, dependencies []StageName) {
	p.stageDependencies[stageName] = dependencies
}

// AddStageBefore adds custom stage before existing stage
func (p *Pipeline) AddStageBefore(existingStageName StageName, name StageName, stageRunner StageRunner) {
	p.beforeStage[existingStageName] = append(p.beforeStage[existingStageName], NewStage(name, stageRunner))
}

// AddStageBefore adds custom stage after existing stage
func (p *Pipeline) AddStageAfter(existingStageName StageName, name StageName, stageRunner StageRunner) {
	p.afterStage[existingStageName] = append(p.afterStage[existingStageName], NewStage(name, stageRunner))
}

// Start starts the pipeline
func (p *Pipeline) Start(ctx context.Context, source Source, sink Sink, options *Options) error {
	pCtx, _ := p.setupCtx(ctx)
	p.options = options

	if err := p.validateStages(); err != nil {
		return err
	}

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

// Run run one-off pipeline iteration for given height
func (p *Pipeline) Run(ctx context.Context, height int64, options *Options) (Payload, error) {
	if err := p.validateStages(); err != nil {
		return nil, err
	}

	pCtx, _ := p.setupCtx(ctx)

	p.options = options

	payload := p.payloadFactory.GetPayload(height)

	if err := p.runStages(pCtx, payload); err != nil {
		return nil, err
	}

	payload.MarkAsProcessed()

	return payload, nil
}

// validateStages verfies that each stage has a dependency configured,
// and that there are no dependency loops
func (p *Pipeline) validateStages() error {
	visited := make(map[StageName]struct{})

	for len(visited) < len(p.stages) {
		var stageRunCount int

		for _, s := range p.stages {
			if _, ok := visited[s.Name]; ok {
				continue
			}

			deps, ok := p.stageDependencies[s.Name]
			if !ok {
				// each stage must have an entry in the dependency map
				return ErrMisconfiguredDependency
			}

			canRun := true
			for _, d := range deps {
				if _, ok := p.stages[d]; !ok {
					// don't block stage from running if dependency is not set up:
					// default pipeline should still run if any stage is not configured
					continue
				}

				if _, ok := visited[d]; !ok {
					canRun = false
					break
				}
			}

			if canRun {
				stageRunCount++
				visited[s.Name] = struct{}{}
			}
		}

		if stageRunCount == 0 {
			return ErrMisconfiguredDependency
		}
	}

	return nil
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
	completedStages := make(map[StageName]struct{})

	for len(completedStages) < len(p.stageDependencies) {
		runNext := make([]StageName, 0)

		for stage, deps := range p.stageDependencies {
			if _, ok := completedStages[stage]; ok {
				continue
			}

			canRun := true
			for _, dep := range deps {
				if _, ok := completedStages[dep]; !ok {
					canRun = false
					break
				}
			}
			if canRun {
				runNext = append(runNext, stage)
			}
		}

		if len(runNext) == 0 {
			return errors.New("no stages to run")
		}

		var runErr error
		if len(runNext) == 1 {
			runErr = p.runStage(runNext[0], ctx, payload)
		} else if len(runNext) > 1 {
			runErr = p.runStagesConcurrently(ctx, payload, runNext)
		}

		if runErr != nil {
			return runErr
		}

		// mark stage as completed
		for _, stage := range runNext {
			completedStages[stage] = struct{}{}
		}
	}

	return nil
}

// runStagesConcurrently runs indexing stages concurrently
func (p *Pipeline) runStagesConcurrently(ctx context.Context, payload Payload, stages []StageName) error {
	stagesCount := len(stages)
	if stagesCount == 0 {
		return ErrMissingStages
	}

	var errs error
	var wg sync.WaitGroup
	errCh := make(chan error, stagesCount)
	wg.Add(stagesCount)

	for _, stageName := range stages {
		go func(stageName StageName) {
			if err := p.runStage(stageName, ctx, payload); err != nil {
				errCh <- err
			}
			wg.Done()
		}(stageName)
	}

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
		logInfo(fmt.Sprintf("stage name %s not set up", stageName))
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
