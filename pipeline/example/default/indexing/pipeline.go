package indexing

import (
	"context"
	"fmt"

	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	versionsDir = "example/indexing/versions"
)

func StartPipeline() error {
	p := pipeline.NewDefault(NewPayloadFactory())

	// Set tasks in fetcher stage
	// Demonstrates use of retrying mechanism for tasks inside the stage
	p.SetTasks(
		pipeline.StageFetcher,
		pipeline.RetryingTask(NewFetcherTask(), func(err error) bool {
			// Make error always transient for simplicity
			return true
		}, 3),
	)

	// Set parser stage
	p.SetTasks(pipeline.StageParser, NewParserTask(), NewParserTask2())

	// Set validator stage
	p.SetAsyncTasks(pipeline.StageValidator, NewValidatorTask(), NewValidatorTask2())

	// Set syncer stage
	p.SetTasks(pipeline.StageSyncer, NewSyncerTask())

	// wraps entire stage with retry mechanism
	p.RetryStage(
		pipeline.StageSyncer,
		func(err error) bool {
			// Make error always transient for simplicity
			return true
		},
		3,
	)

	// Set sequencer stage
	p.SetTasks(pipeline.StageSequencer, NewSequencerTask())

	// Set aggregator stage
	p.SetTasks(pipeline.StageAggregator, NewAggregatorTask())

	// Add custom stage before existing one
	// Demonstrates how to use func as a stage runner without a need to use structs
	beforeFetcherFunc := pipeline.StageRunnerFunc(func(ctx context.Context, p pipeline.Payload, f pipeline.TaskValidator) error {
		payload := (p).(*payload)
		fmt.Println("task: ", "BeforeFetcher", payload.CurrentHeight)
		return nil
	})
	p.AddStageBefore(
		pipeline.StageFetcher,
		pipeline.NewCustomStage("BeforeFetcher", beforeFetcherFunc),
	)

	// Add custom stage after existing one
	// Demonstrates how to use func as a stage runner without a need to use structs
	afterValidatorFunc := pipeline.StageRunnerFunc(func(ctx context.Context, p pipeline.Payload, f pipeline.TaskValidator) error {
		payload := (p).(*payload)
		fmt.Println("task: ", "AfterValidator", payload.CurrentHeight)
		return nil
	})
	p.AddStageAfter(
		pipeline.StageValidator,
		pipeline.NewCustomStage("AfterValidator", afterValidatorFunc),
	)

	ctx := context.Background()

	options := &pipeline.Options{}
	if err := p.Start(ctx, NewSource(), NewSink(), options); err != nil {
		return err
	}
	return nil
}
