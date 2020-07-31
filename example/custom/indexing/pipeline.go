package indexing

import (
	"context"
	"fmt"

	"github.com/figment-networks/indexing-engine/pipeline"
)

func StartPipeline() error {
	// creates new custom pipeline. Will run stages in the order of which they were declared
	p := pipeline.NewCustom(NewPayloadFactory())

	// Adds fetcher stage
	p.AddStage(
		pipeline.NewStageWithTasks(pipeline.StageFetcher, NewFetcherTask()),
	)

	// Adds parser stage
	p.AddStage(
		pipeline.NewAsyncStageWithTasks(pipeline.StageParser, NewParserTask(), NewParserTask2()),
	)

	// Adds syncer stage
	p.AddStage(
		pipeline.NewAsyncStageWithTasks(pipeline.StageSyncer, NewSyncerTask()),
	)
	// Wraps entire stage with retry mechanism
	p.RetryStage(
		pipeline.StageSyncer,
		func(err error) bool {
			// Make error always transient for simplicity
			return true
		},
		3,
	)

	// Adds sequencer and aggregator stages which will run concurrently
	p.AddConcurrentStages(
		pipeline.NewStageWithTasks(pipeline.StageSequencer, NewSequencerTask()),
		pipeline.NewStageWithTasks(pipeline.StageAggregator, NewAggregatorTask()),
	)

	// Demonstrates how to use custom func as a stage runner
	cleanupFunc := pipeline.StageRunnerFunc(func(ctx context.Context, p pipeline.Payload, f pipeline.TaskValidator) error {
		payload := (p).(*payload)
		fmt.Println("task: ", "CleanupExample (Custom Func)", payload.CurrentHeight)
		return nil
	})

	// Adds cleanup stage with custom stagerunner
	p.AddStage(
		pipeline.NewCustomStage(pipeline.StageCleanup, cleanupFunc),
	)

	ctx := context.Background()

	options := &pipeline.Options{}
	if err := p.Start(ctx, NewSource(), NewSink(), options); err != nil {
		return err
	}
	return nil
}
