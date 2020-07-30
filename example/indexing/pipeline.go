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

	// Set fetcher stage
	// Demonstrates use of retrying mechanism for tasks inside the stage
	p.SetStageRunner(
		pipeline.StageFetcher,
		pipeline.AsyncRunner(
			pipeline.RetryingTask(NewFetcherTask(), func(err error) bool {
				// Make error always transient for simplicity
				return true
			}, 3),
		),
	)

	// Set parser stage
	p.SetStageRunner(pipeline.StageParser, pipeline.SyncRunner(NewParserTask()))

	// Set validator stage
	p.SetStageRunner(pipeline.StageValidator, pipeline.SyncRunner(NewValidatorTask()))

	// Set syncer stage
	// Demonstrates use of retrying mechanism for entire stage
	p.SetStageRunner(
		pipeline.StageSyncer,
		pipeline.RetryingStageRunner(pipeline.SyncRunner(NewSyncerTask()), func(err error) bool {
			// Make error always transient for simplicity
			return true
		}, 3),
	)

	// Set sequencer stage
	p.SetStageRunner(pipeline.StageSequencer, pipeline.AsyncRunner(NewSequencerTask()))

	// Set aggregator stage
	p.SetStageRunner(pipeline.StageAggregator, pipeline.AsyncRunner(NewAggregatorTask()))

	// Add custom stage before existing one
	// Demonstrates how to use func as a stage runner without a need to use structs
	beforeFetcherFunc := pipeline.StageRunnerFunc(func(ctx context.Context, p pipeline.Payload, f pipeline.TaskValidator) error {
		payload := (p).(*payload)
		fmt.Println("task: ", "BeforeFetcher", payload.CurrentHeight)
		return nil
	})
	p.AddStageBefore(pipeline.StageFetcher, "BeforeFetcher", beforeFetcherFunc)

	// Add custom stage after existing one
	// Demonstrates how to use func as a stage runner without a need to use structs
	afterValidatorFunc := pipeline.StageRunnerFunc(func(ctx context.Context, p pipeline.Payload, f pipeline.TaskValidator) error {
		payload := (p).(*payload)
		fmt.Println("task: ", "AfterValidator", payload.CurrentHeight)
		return nil
	})
	p.AddStageAfter(pipeline.StageValidator, "AfterValidator", afterValidatorFunc)

	ctx := context.Background()

	options := &pipeline.Options{}
	if err := p.Start(ctx, NewSource(), NewSink(), options); err != nil {
		return err
	}
	return nil
}
