package pipeline_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/indexing-engine/pipeline/mock"
)

var (
	_ pipeline.Source = (*sourceMock)(nil)
)

var (
	allStages = [...]pipeline.StageName{
		pipeline.StageSetup,
		pipeline.StageSyncer,
		pipeline.StageFetcher,
		pipeline.StageParser,
		pipeline.StageValidator,
		pipeline.StageSequencer,
		pipeline.StageAggregator,
		pipeline.StagePersistor,
		pipeline.StageCleanup,
	}
)

type payloadMock struct{}

func (p *payloadMock) MarkAsProcessed() {}

type sourceMock struct {
	startHeight   int64
	endHeight     int64
	currentHeight int64
	skip          bool
}

func (s *sourceMock) Next(context.Context, pipeline.Payload) bool {
	if s.currentHeight == s.endHeight {
		return false
	}
	s.currentHeight = s.currentHeight + 1
	return true
}

func (s *sourceMock) Current() int64 {
	return s.currentHeight
}

func (s *sourceMock) Err() error {
	return nil
}

func (s *sourceMock) Skip(stageName pipeline.StageName) bool {
	return s.skip
}

func TestPipeline_SetStages(t *testing.T) {
	t.Run("all set stages are executed", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(2)

		p := pipeline.NewDefault(payloadFactoryMock)

		setupTaskMock := mock.NewMockTask(ctrl)
		fetcherTaskMock := mock.NewMockTask(ctrl)
		parserTaskMock := mock.NewMockTask(ctrl)
		validatorTaskMock := mock.NewMockTask(ctrl)
		syncerTaskMock := mock.NewMockTask(ctrl)
		sequencerTaskMock := mock.NewMockTask(ctrl)
		aggregatorTaskMock := mock.NewMockTask(ctrl)
		cleanupTaskMock := mock.NewMockTask(ctrl)

		setupTaskMock.EXPECT().GetName().Return("setupTask").Times(4)
		fetcherTaskMock.EXPECT().GetName().Return("fetcherTask").Times(4)
		parserTaskMock.EXPECT().GetName().Return("parserTask").Times(4)
		validatorTaskMock.EXPECT().GetName().Return("validatorTask").Times(4)
		syncerTaskMock.EXPECT().GetName().Return("syncerTask").Times(4)
		sequencerTaskMock.EXPECT().GetName().Return("sequencerTask").Times(4)
		aggregatorTaskMock.EXPECT().GetName().Return("aggregatorTask").Times(4)
		cleanupTaskMock.EXPECT().GetName().Return("cleanupTask").Times(4)

		p.SetTasks(pipeline.StageSetup, setupTaskMock)
		p.SetTasks(pipeline.StageFetcher, fetcherTaskMock)
		p.SetTasks(pipeline.StageParser, parserTaskMock)
		p.SetTasks(pipeline.StageValidator, validatorTaskMock)
		p.SetTasks(pipeline.StageSyncer, syncerTaskMock)
		p.SetTasks(pipeline.StageSequencer, sequencerTaskMock)
		p.SetTasks(pipeline.StageAggregator, aggregatorTaskMock)
		p.SetTasks(pipeline.StageCleanup, cleanupTaskMock)

		setupTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		fetcherTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		parserTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		validatorTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		syncerTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		sequencerTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		aggregatorTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		cleanupTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		sinkMock := mock.NewMockSink(ctrl)
		sinkMock.EXPECT().Consume(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		options := &pipeline.Options{}

		if err := p.Start(ctx, &sourceMock{1, 2, 1, false}, sinkMock, options); err != nil {
			t.Errorf("should not return error")
		}
	})
}

func TestPipeline_Start(t *testing.T) {
	stageErr := errors.New("err")

	t.Run("pipeline runs stages in default order when running all stages", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewDefault(payloadFactoryMock)

		setupTaskMock := mock.NewMockTask(ctrl)
		fetcherTaskMock := mock.NewMockTask(ctrl)
		parserTaskMock := mock.NewMockTask(ctrl)
		validatorTaskMock := mock.NewMockTask(ctrl)
		syncerTaskMock := mock.NewMockTask(ctrl)
		sequencerTaskMock := mock.NewMockTask(ctrl)
		aggregatorTaskMock := mock.NewMockTask(ctrl)
		persistorTaskMock := mock.NewMockTask(ctrl)
		cleanupTaskMock := mock.NewMockTask(ctrl)

		setupTaskMock.EXPECT().GetName().Return("setupTask").Times(2)
		fetcherTaskMock.EXPECT().GetName().Return("fetcherTask").Times(2)
		parserTaskMock.EXPECT().GetName().Return("parserTask").Times(2)
		validatorTaskMock.EXPECT().GetName().Return("validatorTask").Times(2)
		syncerTaskMock.EXPECT().GetName().Return("syncerTask").Times(2)
		sequencerTaskMock.EXPECT().GetName().Return("sequencerTask").Times(2)
		aggregatorTaskMock.EXPECT().GetName().Return("aggregatorTask").Times(2)
		persistorTaskMock.EXPECT().GetName().Return("aggregatorTask").Times(2)
		cleanupTaskMock.EXPECT().GetName().Return("cleanupTask").Times(2)

		p.SetTasks(pipeline.StageSetup, setupTaskMock)
		p.SetTasks(pipeline.StageFetcher, fetcherTaskMock)
		p.SetTasks(pipeline.StageParser, parserTaskMock)
		p.SetTasks(pipeline.StageValidator, validatorTaskMock)
		p.SetTasks(pipeline.StageSyncer, syncerTaskMock)
		p.SetTasks(pipeline.StageSequencer, sequencerTaskMock)
		p.SetTasks(pipeline.StageAggregator, aggregatorTaskMock)
		p.SetTasks(pipeline.StagePersistor, persistorTaskMock)
		p.SetTasks(pipeline.StageCleanup, cleanupTaskMock)
		sinkMock := mock.NewMockSink(ctrl)

		runSetup := setupTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil)
		runSyncer := syncerTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		runFetcher := fetcherTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		runParser := parserTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		runValidator := validatorTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		runSequencer := sequencerTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		runAggregator := aggregatorTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		runPersistor := persistorTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		runCleanup := cleanupTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		runSetup.Times(1)
		runSyncer.After(runSetup)
		runFetcher.After(runSyncer)
		runParser.After(runFetcher)
		runValidator.After(runParser)

		runSequencer.After(runValidator)
		runAggregator.After(runValidator)

		runPersistor.After(runSequencer)
		runPersistor.After(runAggregator)
		runCleanup.After(runPersistor)

		sinkMock.EXPECT().Consume(gomock.Any(), gomock.Any()).Return(nil).Times(1).After(runCleanup)
		options := &pipeline.Options{}

		if err := p.Start(ctx, &sourceMock{1, 1, 1, false}, sinkMock, options); err != nil {
			t.Errorf("did not expect error")
		}
	})

	t.Run("pipeline skip running stages", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewDefault(payloadFactoryMock)
		cleanupTaskMock := mock.NewMockTask(ctrl)
		sinkMock := mock.NewMockSink(ctrl)
		runCleanup := cleanupTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(0)

		sinkMock.EXPECT().Consume(gomock.Any(), gomock.Any()).Return(nil).Times(1).After(runCleanup)
		options := &pipeline.Options{}

		if err := p.Start(ctx, &sourceMock{1, 1, 1, true}, sinkMock, options); err != nil {
			t.Errorf("did not expect error")
		}
	})

	t.Run("pipeline runs stages in default order when running some stages", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewDefault(payloadFactoryMock)

		setupTaskMock := mock.NewMockTask(ctrl)
		parserTaskMock := mock.NewMockTask(ctrl)
		syncerTaskMock := mock.NewMockTask(ctrl)
		cleanupTaskMock := mock.NewMockTask(ctrl)

		setupTaskMock.EXPECT().GetName().Return("setupTask").Times(2)
		parserTaskMock.EXPECT().GetName().Return("parserTask").Times(2)
		syncerTaskMock.EXPECT().GetName().Return("syncerTask").Times(2)
		cleanupTaskMock.EXPECT().GetName().Return("cleanupTask").Times(2)

		p.SetTasks(pipeline.StageSetup, setupTaskMock)
		p.SetTasks(pipeline.StageParser, parserTaskMock)
		p.SetTasks(pipeline.StageSyncer, syncerTaskMock)
		p.SetTasks(pipeline.StageCleanup, cleanupTaskMock)
		sinkMock := mock.NewMockSink(ctrl)

		gomock.InOrder(
			setupTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1),
			syncerTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1),
			parserTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1),
			cleanupTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1),
			sinkMock.EXPECT().Consume(gomock.Any(), gomock.Any()).Return(nil).Times(1),
		)

		options := &pipeline.Options{}

		if err := p.Start(ctx, &sourceMock{1, 1, 1, false}, sinkMock, options); err != nil {
			t.Errorf("did not expect error")
		}
	})

	t.Run("pipeline returns error if syncing stage errors", func(t *testing.T) {
		for _, stageWithErr := range [...]pipeline.StageName{
			pipeline.StageSetup,
			pipeline.StageSyncer,
			pipeline.StageFetcher,
			pipeline.StageParser,
			pipeline.StageValidator,
			pipeline.StagePersistor,
			pipeline.StageCleanup,
		} {

			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
			payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

			p := pipeline.NewDefault(payloadFactoryMock)

			shouldRun := true
			for _, stage := range allStages {
				var returnVal error
				mockTask := mock.NewMockTask(ctrl)

				if !shouldRun {
					mockTask.EXPECT().GetName().Return("mockTask").Times(0)
					mockTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(0)
					p.SetTasks(stage, mockTask)
					continue
				}

				if stage == stageWithErr {
					returnVal = stageErr
					shouldRun = false
				}

				mockTask.EXPECT().GetName().Return("mockTask").Times(2)
				mockTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(returnVal).Times(1)
				p.SetTasks(stage, mockTask)
			}

			sinkMock := mock.NewMockSink(ctrl)
			sinkMock.EXPECT().Consume(gomock.Any(), gomock.Any()).Return(nil).Times(0)

			options := &pipeline.Options{}

			if err := p.Start(ctx, &sourceMock{1, 2, 1, false}, sinkMock, options); err != stageErr {
				t.Errorf("expected error")
			}
		}
	})

	t.Run("pipeline returns error if async stage errors", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewDefault(payloadFactoryMock)

		aggregatorTask := mock.NewMockTask(ctrl)
		aggregatorTask.EXPECT().GetName().Return("aggregatorTask").Times(2)
		aggregatorTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(stageErr).Times(1)
		p.SetTasks(pipeline.StageAggregator, aggregatorTask)

		sequencerTask := mock.NewMockTask(ctrl)
		sequencerTask.EXPECT().GetName().Return("sequencerTask").Times(2)
		sequencerTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		p.SetTasks(pipeline.StageSequencer, sequencerTask)

		cleanupTask := mock.NewMockTask(ctrl)
		cleanupTask.EXPECT().GetName().Return("cleanupTask").Times(0)
		cleanupTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(0)
		p.SetTasks(pipeline.StageCleanup, cleanupTask)

		sinkMock := mock.NewMockSink(ctrl)
		sinkMock.EXPECT().Consume(gomock.Any(), gomock.Any()).Return(nil).Times(0)

		options := &pipeline.Options{}

		if err := p.Start(ctx, &sourceMock{1, 2, 1, false}, sinkMock, options); err == nil {
			t.Errorf("expected error")
		}
	})
}

func TestPipeline_NewCustom(t *testing.T) {
	t.Run("custom pipeline runs stages in custom order", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewCustom(payloadFactoryMock)

		runCalls := []*gomock.Call{}

		for _, stageName := range []pipeline.StageName{
			pipeline.StageParser, pipeline.StageAggregator, pipeline.StageSetup,
		} {
			mockTask := mock.NewMockTask(ctrl)
			mockTask.EXPECT().GetName().Return("mockTask").Times(2)

			call := mockTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			runCalls = append(runCalls, call)

			p.AddStage(pipeline.NewStageWithTasks(stageName, mockTask))
		}

		sinkMock := mock.NewMockSink(ctrl)
		runCalls = append(runCalls,
			sinkMock.EXPECT().Consume(gomock.Any(), gomock.Any()).Return(nil).Times(1))

		gomock.InOrder(runCalls...)

		if err := p.Start(ctx, &sourceMock{1, 1, 1, false}, sinkMock, nil); err != nil {
			t.Errorf("did not expect error")
		}
	})

	t.Run("custom pipeline runs tasks in concurrent stages", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewCustom(payloadFactoryMock)

		mockTasks := make([]pipeline.Task, 3)

		for i := 0; i < 3; i++ {
			mockTask := mock.NewMockTask(ctrl)
			mockTask.EXPECT().GetName().Return("mockTask").Times(2)
			mockTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			mockTasks[i] = mockTask
		}

		p.AddConcurrentStages(
			pipeline.NewStageWithTasks(pipeline.StageCleanup, mockTasks[0]),
			pipeline.NewStageWithTasks(pipeline.StageParser, mockTasks[1]),
			pipeline.NewStageWithTasks(pipeline.StageSetup, mockTasks[2]),
		)

		sinkMock := mock.NewMockSink(ctrl)
		sinkMock.EXPECT().Consume(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		if err := p.Start(ctx, &sourceMock{1, 1, 1, false}, sinkMock, nil); err != nil {
			t.Errorf("did not expect error")
		}
	})
}

func TestPipeline_AddStageBefore(t *testing.T) {
	t.Run("new stage is executed before given stage", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewDefault(payloadFactoryMock)

		stages := []struct {
			name         pipeline.StageName
			existingName pipeline.StageName
		}{
			{"beforeSetup", pipeline.StageSetup},
			{"beforeFetcher", pipeline.StageFetcher},
			{"beforeParser", pipeline.StageParser},
			{"beforeCleanup", pipeline.StageCleanup},
		}

		for _, stage := range stages {
			existingStageTask := mock.NewMockTask(ctrl)
			existingStageTask.EXPECT().GetName().Return("mockTask").Times(2)
			p.SetTasks(stage.existingName, existingStageTask)

			beforeTask := mock.NewMockTask(ctrl)
			beforeTask.EXPECT().GetName().Return("mockTask").Times(2)
			p.AddStageBefore(stage.existingName, pipeline.NewStageWithTasks(stage.name, beforeTask))

			gomock.InOrder(
				beforeTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1),
				existingStageTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1),
			)
		}

		options := &pipeline.Options{}

		if _, err := p.Run(ctx, 1, options); err != nil {
			t.Errorf("should not return error")
		}
	})

	t.Run("pipeline returns err", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		stageErr := errors.New("err")
		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewDefault(payloadFactoryMock)

		beforeTask := mock.NewMockTask(ctrl)
		beforeTask.EXPECT().GetName().Return("mockTask").Times(2)
		beforeTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(stageErr).Times(1)
		p.AddStageBefore(pipeline.StageFetcher, pipeline.NewStageWithTasks("beforeFetcher", beforeTask))

		existingStageTask := mock.NewMockTask(ctrl)
		existingStageTask.EXPECT().GetName().Return("mockTask").Times(0)
		existingStageTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(0)
		p.SetTasks(pipeline.StageFetcher, existingStageTask)

		options := &pipeline.Options{}

		if _, err := p.Run(ctx, 1, options); err != stageErr {
			t.Errorf("should return error")
		}
	})
}

func TestPipeline_AddStageAfter(t *testing.T) {
	t.Run("new stage is executed after existing stage", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewDefault(payloadFactoryMock)

		stages := []struct {
			name         pipeline.StageName
			existingName pipeline.StageName
		}{
			{"afterSetup", pipeline.StageSetup},
			{"afterFetcher", pipeline.StageFetcher},
			{"afterParser", pipeline.StageParser},
			{"afterCleanup", pipeline.StageCleanup},
		}

		for _, stage := range stages {
			existingStageTask := mock.NewMockTask(ctrl)
			existingStageTask.EXPECT().GetName().Return("mockTask").Times(2)
			p.SetTasks(stage.existingName, existingStageTask)

			afterTask := mock.NewMockTask(ctrl)
			afterTask.EXPECT().GetName().Return("mockTask").Times(2)
			p.AddStageAfter(stage.existingName, pipeline.NewStageWithTasks(stage.name, afterTask))

			gomock.InOrder(
				existingStageTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1),
				afterTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1),
			)
		}

		options := &pipeline.Options{}

		if _, err := p.Run(ctx, 1, options); err != nil {
			t.Errorf("should not return error")
		}
	})

	t.Run("pipeline returns err", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		stageErr := errors.New("err")
		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewDefault(payloadFactoryMock)

		afterTask := mock.NewMockTask(ctrl)
		afterTask.EXPECT().GetName().Return("mockTask").Times(2)
		afterTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(stageErr).Times(1)
		p.AddStageAfter(pipeline.StageFetcher, pipeline.NewStageWithTasks("afterFetcher", afterTask))

		existingStageTask := mock.NewMockTask(ctrl)
		existingStageTask.EXPECT().GetName().Return("mockTask").Times(2)
		existingStageTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		p.SetTasks(pipeline.StageFetcher, existingStageTask)

		options := &pipeline.Options{}

		if _, err := p.Run(ctx, 1, options); err != stageErr {
			t.Errorf("should return error")
		}
	})
}

func TestPipeline_StagesBlacklist(t *testing.T) {
	t.Run("blacklisted stage should not run", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewDefault(payloadFactoryMock)

		rand.Seed(time.Now().Unix())
		blacklistedStage := allStages[rand.Intn(len(allStages))]
		options := &pipeline.Options{
			StagesBlacklist: []pipeline.StageName{blacklistedStage},
		}

		for _, stage := range allStages {
			var calls int
			if stage != blacklistedStage {
				calls = 1
			}
			mockTask := mock.NewMockTask(ctrl)
			mockTask.EXPECT().GetName().Return("mockTask").Times(calls * 2)
			mockTask.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(calls)
			p.SetTasks(stage, mockTask)
		}

		if _, err := p.Run(ctx, 1, options); err != nil {
			t.Errorf("should not return error")
		}
	})
}

func TestPipeline_RetryStage(t *testing.T) {
	t.Run("tasks return success", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewCustom(payloadFactoryMock)

		task1 := mock.NewMockTask(ctrl)
		task2 := mock.NewMockTask(ctrl)

		task1.EXPECT().GetName().Return("task1").Times(2)
		task2.EXPECT().GetName().Return("task2").Times(2)

		gomock.InOrder(
			task1.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1),
			task2.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1),
		)

		p.AddStage(pipeline.NewStageWithTasks("test_stage", task1, task2))
		p.RetryStage("test_stage", func(err error) bool { return true },
			3)

		if _, err := p.Run(ctx, 1, nil); err != nil {
			t.Errorf("should not return error")
		}
	})

	t.Run("task returns error then success", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewCustom(payloadFactoryMock)

		task1 := mock.NewMockTask(ctrl)
		task1.EXPECT().GetName().Return("task1").Times(4)

		gomock.InOrder(
			task1.EXPECT().Run(gomock.Any(), gomock.Any()).Return(errors.New("test err")).Times(1),
			task1.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1),
		)

		p.AddStage(pipeline.NewStageWithTasks("test_stage", task1))
		p.RetryStage("test_stage", func(err error) bool { return true },
			3)

		if _, err := p.Run(ctx, 1, nil); err != nil {
			t.Errorf("should not return error")
		}
	})

	t.Run("task returns error every try", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(1)

		p := pipeline.NewCustom(payloadFactoryMock)

		task1 := mock.NewMockTask(ctrl)
		task1.EXPECT().GetName().Return("task1").Times(6)
		task1.EXPECT().Run(gomock.Any(), gomock.Any()).Return(errors.New("test error")).Times(3)

		p.AddStage(pipeline.NewStageWithTasks("test_stage", task1))
		p.RetryStage("test_stage", func(err error) bool { return true },
			3)

		if _, err := p.Run(ctx, 1, nil); err == nil {
			t.Errorf("should return error")
		}
	})
}
