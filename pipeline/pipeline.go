package pipeline

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"

	"github.com/figment-networks/indexing-engine/metrics"
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
	ErrMissingStages = errors.New("provide stages to run concurrently")
	ErrMissingStage  = errors.New("no stage to run")
)

type StageName string
type TaskName string

type Options struct {
	// StagesBlacklist holds list of stages to turn off
	StagesBlacklist []StageName

	// TaskWhitelist holds name of indexing tasks which will be executed
	TaskWhitelist []TaskName
}

type Pipeline interface {
	SetLogger(l Logger)
	AddStageBefore(existingStageName StageName, stage *stage)
	AddStageAfter(existingStageName StageName, stage *stage)
	RetryStage(existingStageName StageName, isTransient func(error) bool, maxRetries int)
	Start(ctx context.Context, source Source, sink Sink, options *Options) error
	Run(ctx context.Context, height int64, options *Options) (Payload, error)
}

// DefaultPipeline is implemented by types that only want to configure existing stages in a pipeline
type DefaultPipeline interface {
	Pipeline

	SetTasks(stageName StageName, tasks ...Task)
	SetAsyncTasks(stageName StageName, tasks ...Task)
	SetCustomStage(stageName StageName, stageRunnerFunc stageRunner)
}

// CustomPipeline is implemented by types that want to create a pipeline by adding their own stages
type CustomPipeline interface {
	Pipeline

	AddStage(stage *stage)
	AddConcurrentStages(stages ...*stage)
}

// NewDefault creates a new DefaultPipeline with all default stages set in default run order
func NewDefault(payloadFactor PayloadFactory) DefaultPipeline {
	p := new(payloadFactor)

	emptyRunner := func(name StageName) StageRunnerFunc {
		return StageRunnerFunc(func(context.Context, Payload, TaskValidator) error {
			logInfo(fmt.Sprintf("stage name %s not set up", name))
			return nil
		})
	}

	p.AddStage(NewCustomStage(StageSetup, emptyRunner(StageSetup)))
	p.AddStage(NewCustomStage(StageSyncer, emptyRunner(StageSyncer)))
	p.AddStage(NewCustomStage(StageFetcher, emptyRunner(StageFetcher)))
	p.AddStage(NewCustomStage(StageParser, emptyRunner(StageParser)))
	p.AddStage(NewCustomStage(StageValidator, emptyRunner(StageValidator)))
	p.AddConcurrentStages(
		NewCustomStage(StageSequencer, emptyRunner(StageSequencer)),
		NewCustomStage(StageAggregator, emptyRunner(StageAggregator)),
	)
	p.AddStage(NewCustomStage(StagePersistor, emptyRunner(StagePersistor)))
	p.AddStage(NewCustomStage(StageCleanup, emptyRunner(StageCleanup)))

	return p
}

// NewCustom creates a new pipeline that satisfies CustomPipeline
func NewCustom(payloadFactor PayloadFactory) CustomPipeline {
	return new(payloadFactor)
}

// pipeline implements a modular, multi-stage pipeline
type pipeline struct {
	payloadFactory PayloadFactory
	options        *Options

	stages [][]*stage

	beforeStage map[StageName][]*stage
	afterStage  map[StageName][]*stage
}

func new(payloadFactor PayloadFactory) *pipeline {
	return &pipeline{
		payloadFactory: payloadFactor,
		stages:         [][]*stage{},

		beforeStage: make(map[StageName][]*stage),
		afterStage:  make(map[StageName][]*stage),
	}
}

// SetLogger sets logger
func (p *pipeline) SetLogger(l Logger) {
	logger = l
}

// SetAsyncTasks adds tasks which will run concurrently in a given stage
func (p *pipeline) SetAsyncTasks(stageName StageName, tasks ...Task) {
	p.setRunnerForStage(stageName, asyncRunner{tasks: tasks})
}

// SetTasks adds tasks which will run one by one in a given stage
func (p *pipeline) SetTasks(stageName StageName, tasks ...Task) {
	p.setRunnerForStage(stageName, syncRunner{tasks: tasks})
}

// SetCustomStage sets custom stage runner for a given stage
func (p *pipeline) SetCustomStage(stageName StageName, runner stageRunner) {
	p.setRunnerForStage(stageName, runner)
}

func (p *pipeline) setRunnerForStage(stageName StageName, runner stageRunner) {
	for _, stages := range p.stages {
		for _, s := range stages {
			if s.Name == stageName {
				s.runner = runner
				return
			}
		}
	}
	logInfo(fmt.Sprintf("cannot set stage runner for stage, stage '%v' not found on pipeline", stageName))
}

// AddConcurrentStages adds stages that will run concurrently in the pipeline
func (p *pipeline) AddConcurrentStages(stages ...*stage) {
	if len(stages) == 0 {
		return
	}
	p.stages = append(p.stages, stages)
}

// AddStage adds stage to pipeline
func (p *pipeline) AddStage(s *stage) {
	p.stages = append(p.stages, []*stage{s})
}

// AddStageBefore adds new stage before existing stage
func (p *pipeline) AddStageBefore(existingStageName StageName, s *stage) {
	p.beforeStage[existingStageName] = append(p.beforeStage[existingStageName], s)
}

// AddStageAfter adds new stage after existing stage
func (p *pipeline) AddStageAfter(existingStageName StageName, s *stage) {
	p.afterStage[existingStageName] = append(p.afterStage[existingStageName], s)
}

// RetryStage implements retry mechanism for entire stage
func (p *pipeline) RetryStage(existingStageName StageName, isTransient func(error) bool, maxRetries int) {
	for _, stages := range p.stages {
		for _, s := range stages {
			if s.Name == existingStageName {
				s.runner = retryingStageRunner(s.runner, isTransient, maxRetries)
			}
		}
	}
}

// Start starts the pipeline
func (p *pipeline) Start(ctx context.Context, source Source, sink Sink, options *Options) error {
	pCtx, _ := p.setupCtx(ctx)
	p.options = options

	heightCounter := heightsTotalMetric.WithLabels()
	durationObserver := heightDurationMetric.WithLabels()

	var pipelineErr error
	var recentPayload Payload
	for ok := true; ok; ok = source.Next(ctx, recentPayload) {
		payload := p.payloadFactory.GetPayload(source.Current())

		timer := metrics.NewTimer(durationObserver)

		pipelineErr = p.runStages(pCtx, payload, source)
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

		timer.ObserveDuration()
		heightCounter.Inc()

		recentPayload = payload
	}

	if err := source.Err(); err != nil {
		pipelineErr = multierror.Append(pipelineErr, err)
	}

	if pipelineErr != nil {
		errorsTotalMetric.WithLabels().Inc()
	}

	return pipelineErr
}

// Run run one-off pipeline iteration for given height
func (p *pipeline) Run(ctx context.Context, height int64, options *Options) (Payload, error) {
	pCtx, _ := p.setupCtx(ctx)

	p.options = options

	payload := p.payloadFactory.GetPayload(height)

	observer := heightDurationMetric.WithLabels()
	timer := metrics.NewTimer(observer)

	if err := p.runStages(pCtx, payload, NewSource()); err != nil {
		errorsTotalMetric.WithLabels().Inc()
		return nil, err
	}

	payload.MarkAsProcessed()

	timer.ObserveDuration()
	heightsTotalMetric.WithLabels().Inc()

	return payload, nil
}

// setupCtx sets up the context
func (p *pipeline) setupCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	// Setup cancel
	pCtx, cancelFunc := context.WithCancel(ctx)

	// Setup stats recorder
	statRecorder := NewStatsRecorder()
	pCtx = context.WithValue(pCtx, CtxStats, statRecorder)

	return pCtx, cancelFunc
}

// runStages runs all the stages
func (p *pipeline) runStages(ctx context.Context, payload Payload, source Source) error {
	for _, stages := range p.stages {
		if len(stages) == 1 {
			if err := p.runStage(ctx, stages[0], payload, source); err != nil {
				return err
			}
		} else if len(stages) > 1 {
			if err := p.runStagesConcurrently(ctx, payload, stages, source); err != nil {
				return err
			}
		} else {
			logInfo("no stages to run")
		}
	}

	return nil
}

// runStagesConcurrently runs indexing stages concurrently
func (p *pipeline) runStagesConcurrently(ctx context.Context, payload Payload, stages []*stage, source Source) error {
	stagesCount := len(stages)
	if stagesCount == 0 {
		return ErrMissingStages
	}

	var errs error
	var wg sync.WaitGroup
	errCh := make(chan error, stagesCount)
	wg.Add(stagesCount)

	for _, s := range stages {
		go func(stage *stage) {
			if err := p.runStage(ctx, stage, payload, source); err != nil {
				errCh <- err
			}
			wg.Done()
		}(s)
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
func (p *pipeline) runStage(ctx context.Context, stage *stage, payload Payload, source Source) error {
	if stage == nil {
		return ErrMissingStage
	}

	if p.canRunStage(stage.Name, source) {
		before := p.beforeStage[stage.Name]
		if len(before) > 0 {
			for _, s := range before {
				if err := s.Run(ctx, payload, p.options); err != nil {
					return err
				}
			}
		}

		if err := stage.Run(ctx, payload, p.options); err != nil {
			return err
		}

		after := p.afterStage[stage.Name]
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
func (p *pipeline) canRunStage(stageName StageName, source Source) bool {
	if p.options != nil && len(p.options.StagesBlacklist) > 0 {
		for _, s := range p.options.StagesBlacklist {
			if s == stageName {
				return false
			}
		}
	}
	return !source.Skip(stageName)
}
