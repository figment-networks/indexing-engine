package pipeline_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/figment-networks/indexing-engine.git/pipeline"
	mock "github.com/figment-networks/indexing-engine.git/pipeline/mock"
)

func TestStage_Running(t *testing.T) {
	t.Run("Run() runs stage runner", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		stageRunnerMock := mock.NewMockStageRunner(ctrl)
		stageRunnerMock.EXPECT().Run(ctx, payloadMock, gomock.Any()).Return(nil)

		s := pipeline.NewStage("test", stageRunnerMock, pipeline.StageTypeCustom)

		err := s.Run(ctx, payloadMock, nil)
		if err != nil {
			t.Errorf("exp: nil, got: %f", err)
		}
	})

	t.Run("Run() runs stage runner with error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		stageRunnerMock := mock.NewMockStageRunner(ctrl)
		stageRunnerMock.EXPECT().Run(ctx, payloadMock, gomock.Any()).Return(errors.New("test error"))

		s := pipeline.NewStage("test", stageRunnerMock, pipeline.StageTypeCustom)

		err := s.Run(ctx, payloadMock, nil)
		if err == nil {
			t.Errorf("exp: %f, got: nil", err)
		}
	})
}

func TestStage_SyncRunner(t *testing.T) {
	t.Run("both tasks return success", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		task1 := mock.NewMockTask(ctrl)
		task2 := mock.NewMockTask(ctrl)

		gomock.InOrder(
			task1.EXPECT().Run(ctx, payloadMock).Return(nil),
			task2.EXPECT().Run(ctx, payloadMock).Return(nil),
		)

		var taskValidator pipeline.TaskValidator = func(n string) bool { return true }
		sr := pipeline.SyncRunner(task1, task2)

		err := sr.Run(ctx, payloadMock, taskValidator)
		if err != nil {
			t.Errorf("should not return error")
		}
	})

	t.Run("first task return error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		task1 := mock.NewMockTask(ctrl)
		task2 := mock.NewMockTask(ctrl)

		task1.EXPECT().Run(ctx, payloadMock).Return(errors.New("test error"))
		task2.EXPECT().Run(ctx, payloadMock).Return(nil).Times(0)

		var taskValidator pipeline.TaskValidator = func(n string) bool { return true }
		sr := pipeline.SyncRunner(task1, task2)

		err := sr.Run(ctx, payloadMock, taskValidator)
		if err == nil {
			t.Errorf("should return error")
		}
	})

	t.Run("second task return error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		task1 := mock.NewMockTask(ctrl)
		task2 := mock.NewMockTask(ctrl)

		task1.EXPECT().Run(ctx, payloadMock).Return(nil)
		task2.EXPECT().Run(ctx, payloadMock).Return(errors.New("test error"))

		var taskValidator pipeline.TaskValidator = func(n string) bool { return true }
		sr := pipeline.SyncRunner(task1, task2)

		err := sr.Run(ctx, payloadMock, taskValidator)
		if err == nil {
			t.Errorf("should return error")
		}
	})
}

func TestStage_AsyncRunner(t *testing.T) {
	t.Run("both tasks return success", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		task1 := mock.NewMockTask(ctrl)
		task2 := mock.NewMockTask(ctrl)

		task1.EXPECT().Run(ctx, payloadMock).Return(nil)
		task2.EXPECT().Run(ctx, payloadMock).Return(nil)

		var taskValidator pipeline.TaskValidator = func(n string) bool { return true }
		sr := pipeline.AsyncRunner(task1, task2)

		err := sr.Run(ctx, payloadMock, taskValidator)
		if err != nil {
			t.Errorf("should not return error")
		}
	})

	t.Run("first task return error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		task1 := mock.NewMockTask(ctrl)
		task2 := mock.NewMockTask(ctrl)

		task1.EXPECT().Run(ctx, payloadMock).Return(errors.New("test error"))
		task2.EXPECT().Run(ctx, payloadMock).Return(nil)

		var taskValidator pipeline.TaskValidator = func(n string) bool { return true }
		sr := pipeline.AsyncRunner(task1, task2)

		err := sr.Run(ctx, payloadMock, taskValidator)
		if err == nil {
			t.Errorf("should return error")
		}
	})

	t.Run("second task return error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		task1 := mock.NewMockTask(ctrl)
		task2 := mock.NewMockTask(ctrl)

		task1.EXPECT().Run(ctx, payloadMock).Return(nil)
		task2.EXPECT().Run(ctx, payloadMock).Return(errors.New("test error"))

		var taskValidator pipeline.TaskValidator = func(n string) bool { return true }
		sr := pipeline.AsyncRunner(task1, task2)

		err := sr.Run(ctx, payloadMock, taskValidator)
		if err == nil {
			t.Errorf("should return error")
		}
	})

	t.Run("both tasks return error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		task1 := mock.NewMockTask(ctrl)
		task2 := mock.NewMockTask(ctrl)

		task1.EXPECT().Run(ctx, payloadMock).Return(errors.New("test error"))
		task2.EXPECT().Run(ctx, payloadMock).Return(errors.New("test error"))

		var taskValidator pipeline.TaskValidator = func(n string) bool { return true }
		sr := pipeline.AsyncRunner(task1, task2)

		err := sr.Run(ctx, payloadMock, taskValidator)
		if err == nil {
			t.Errorf("should return error")
		}
	})
}
