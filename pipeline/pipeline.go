package pipeline

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
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
	AddStageBefore(existingStageName StageName, name StageName, stageRunner StageRunner)
	AddStageAfter(existingStageName StageName, name StageName, stageRunner StageRunner)
	Start(ctx context.Context, source Source, sink Sink, options *Options) error
	Run(ctx context.Context, height int64, options *Options) (Payload, error)
}

// DefaultPipeline is implemented by types that want to use set stages
type DefaultPipeline interface {
	Pipeline

	SetStageRunner(stageName StageName, stageRunner StageRunner)
}

// CustomPipeline is implemented by types that want to add custom stages
type CustomPipeline interface {
	Pipeline

	AddStage(stageName StageName, stageRunner StageRunner)
	AddConcurrentStages(stages ...*stage)
}

// NewDefault creates a new DefaultPipeline with default stages set in default run order
func NewDefault(payloadFactor PayloadFactory) DefaultPipeline {
	p := new(payloadFactor)

	emptyRunner := StageRunnerFunc(func(context.Context, Payload, TaskValidator) error {
		return nil
	})

	p.AddStage(StageSetup, emptyRunner)
	p.AddStage(StageSyncer, emptyRunner)
	p.AddStage(StageFetcher, emptyRunner)
	p.AddStage(StageParser, emptyRunner)
	p.AddStage(StageValidator, emptyRunner)
	p.AddConcurrentStages(NewStage(StageSequencer, emptyRunner), NewStage(StageAggregator, emptyRunner))
	p.AddStage(StagePersistor, emptyRunner)
	p.AddStage(StageCleanup, emptyRunner)

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

// SetOptions sets pipeline options
func (p *pipeline) SetOptions(o *Options) {
	p.options = o
}

// SetLogger sets logger
func (p *pipeline) SetLogger(l Logger) {
	logger = l
}

// SetStageRunner sets up stagerunner for given stage
func (p *pipeline) SetStageRunner(stageName StageName, stageRunner StageRunner) {
	for _, stages := range p.stages {
		for _, s := range stages {
			if s.Name == stageName {
				s.runner = stageRunner
			}
		}
	}
}

// AddStage adds stage to pipeline
func (p *pipeline) AddStage(stageName StageName, stageRunner StageRunner) {
	p.stages = append(p.stages, []*stage{NewStage(stageName, stageRunner)})
}

// AddConcurrentStages adds multiple stages that will run concurrently in the pipeline
func (p *pipeline) AddConcurrentStages(stages ...*stage) {
	p.stages = append(p.stages, stages)
}

// AddStageBefore adds custom stage before existing stage
func (p *pipeline) AddStageBefore(existingStageName StageName, name StageName, stageRunner StageRunner) {
	p.beforeStage[existingStageName] = append(p.beforeStage[existingStageName], NewStage(name, stageRunner))
}

// AddStageBefore adds custom stage after existing stage
func (p *pipeline) AddStageAfter(existingStageName StageName, name StageName, stageRunner StageRunner) {
	p.afterStage[existingStageName] = append(p.afterStage[existingStageName], NewStage(name, stageRunner))
}

// Start starts the pipeline
func (p *pipeline) Start(ctx context.Context, source Source, sink Sink, options *Options) error {
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

// Run run one-off pipeline iteration for given height
func (p *pipeline) Run(ctx context.Context, height int64, options *Options) (Payload, error) {
	pCtx, _ := p.setupCtx(ctx)

	p.options = options

	payload := p.payloadFactory.GetPayload(height)

	if err := p.runStages(pCtx, payload); err != nil {
		return nil, err
	}

	payload.MarkAsProcessed()

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
func (p *pipeline) runStages(ctx context.Context, payload Payload) error {
	for _, stages := range p.stages {
		if len(stages) == 1 {
			if err := p.runStage(ctx, stages[0], payload); err != nil {
				return err
			}
		} else if len(stages) > 1 {
			if err := p.runStagesConcurrently(ctx, payload, stages); err != nil {
				return err
			}
		} else {
			logInfo("no stages to run")
		}
	}

	return nil
}

// runStagesConcurrently runs indexing stages concurrently
func (p *pipeline) runStagesConcurrently(ctx context.Context, payload Payload, stages []*stage) error {
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
			if err := p.runStage(ctx, stage, payload); err != nil {
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
func (p *pipeline) runStage(ctx context.Context, stage *stage, payload Payload) error {
	if stage == nil {
		return ErrMissingStage
	}

	if p.canRunStage(stage.Name) {
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
func (p *pipeline) canRunStage(stageName StageName) bool {
	if p.options != nil && len(p.options.StagesBlacklist) > 0 {
		for _, s := range p.options.StagesBlacklist {
			if s == stageName {
				return false
			}
		}
	}

	return true
}
