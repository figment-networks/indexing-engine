package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
)

func StartPipeline() error {
	p := pipeline.New(NewPayloadFactory())

	// Set options to control what stages a what indexing tasks to execute
	p.SetOptions(&pipeline.Options{
		//StagesWhitelist: []pipeline.StageName{pipeline.StageFetcher},
		//IndexingTasksWhitelist: []string{"SequencerTask"},
	})

	// Set fetcher stage
	// Demonstrates use of retrying mechanism for tasks inside the stage
	p.SetFetcherStage(
		pipeline.AsyncRunner(
			pipeline.RetryingTask(NewFetcherTask(), func(err error) bool {
				// Make error always transient for simplicity
				return true
			}, 3),
		),
	)

	// Set parser stage
	p.SetParserStage(pipeline.SyncRunner(NewParserTask()))

	// Set validator stage
	p.SetValidatorStage(pipeline.SyncRunner(NewValidatorTask()))

	// Set syncer stage
	// Demonstrates use of retrying mechanism for entire stage
	p.SetSyncerStage(
		pipeline.RetryingStageRunner(pipeline.SyncRunner(NewSyncerTask()), func(err error) bool {
			// Make error always transient for simplicity
			return true
		}, 3),
	)

	// Set sequencer stage
	p.SetSequencerStage(pipeline.AsyncRunner(NewSequencerTask()))

	// Set aggregator stage
	p.SetAggregatorStage(pipeline.AsyncRunner(NewAggregatorTask()))

	// Add custom stage before existing one
	// Demonstrates how to use func as a stage runner without a need to use structs
	beforeFetcherFunc := pipeline.StageRunnerFunc(func(ctx context.Context, p pipeline.Payload, f pipeline.TaskValidator) error {
		payload := (p).(*payload)
		fmt.Println("task: ", "BeforeFetcher", payload.currentHeight)
		return nil
	})
	p.AddStageBefore(pipeline.StageFetcher, "BeforeFetcher", beforeFetcherFunc)

	// Add custom stage after existing one
	// Demonstrates how to use func as a stage runner without a need to use structs
	afterValidatorFunc := pipeline.StageRunnerFunc(func(ctx context.Context, p pipeline.Payload, f pipeline.TaskValidator) error {
		payload := (p).(*payload)
		fmt.Println("task: ", "AfterValidator", payload.currentHeight)
		return nil
	})
	p.AddStageAfter(pipeline.StageValidator, "AfterValidator", afterValidatorFunc)

	ctx := context.Background()

	err := p.Start(ctx, NewSource(), NewSink())
	if err != nil {
		return err
	}
	return nil
}
