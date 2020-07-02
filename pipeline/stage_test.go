package pipeline_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/figment-networks/indexing-engine/pipeline"
	mock "github.com/figment-networks/indexing-engine/pipeline/mock"
)

func TestStage_Running(t *testing.T) {
	t.Run("Run() runs stage runner", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		stageRunnerMock := mock.NewMockStageRunner(ctrl)
		stageRunnerMock.EXPECT().Run(ctx, payloadMock, gomock.Any()).Return(nil)

		s := pipeline.NewStage("test", stageRunnerMock)

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

		s := pipeline.NewStage("test", stageRunnerMock)

		err := s.Run(ctx, payloadMock, nil)
		if err == nil {
			t.Errorf("exp: %f, got: nil", err)
		}
	})

	t.Run("Run() runs only whitelisted tasks", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		options := &pipeline.Options{
			TaskWhitelist: []pipeline.TaskName{"whitelistTask"},
		}

		task1 := mock.NewMockTask(ctrl)
		task1.EXPECT().GetName().Return("otherTask").Times(1)
		task1.EXPECT().Run(gomock.Any(), gomock.Any()).Times(0)

		task2 := mock.NewMockTask(ctrl)
		task2.EXPECT().GetName().Return("whitelistTask").Times(1)
		task2.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		s := pipeline.NewStage("test", pipeline.SyncRunner(task1, task2))

		err := s.Run(ctx, payloadMock, options)
		if err != nil {
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

		task1.EXPECT().GetName().Return("task1")
		task2.EXPECT().GetName().Return("task2")

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

		task1.EXPECT().GetName().Return("task1")

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

		task1.EXPECT().GetName().Return("task1")
		task2.EXPECT().GetName().Return("task2")

		task1.EXPECT().Run(ctx, payloadMock).Return(nil)
		task2.EXPECT().Run(ctx, payloadMock).Return(errors.New("test error"))

		var taskValidator pipeline.TaskValidator = func(n string) bool { return true }
		sr := pipeline.SyncRunner(task1, task2)

		err := sr.Run(ctx, payloadMock, taskValidator)
		if err == nil {
			t.Errorf("should return error")
		}
	})

	t.Run("returns success when RetryTask fails then succeeds", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		retryTask := mock.NewMockTask(ctrl)
		task2 := mock.NewMockTask(ctrl)

		retryTask.EXPECT().GetName().Return("retryTask").Times(1)
		task2.EXPECT().GetName().Return("task2").Times(1)

		gomock.InOrder(
			retryTask.EXPECT().Run(ctx, payloadMock).Return(errors.New("test error")),
			retryTask.EXPECT().Run(ctx, payloadMock).Return(nil),
			task2.EXPECT().Run(ctx, payloadMock).Return(nil),
		)

		var taskValidator pipeline.TaskValidator = func(n string) bool { return true }
		sr := pipeline.SyncRunner(pipeline.RetryingTask(retryTask, func(err error) bool {
			return true
		}, 3), task2)

		err := sr.Run(ctx, payloadMock, taskValidator)
		if err != nil {
			t.Errorf("should not return error")
		}
	})

	t.Run("returns error when RetryTask fails all attempts", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		retryTask := mock.NewMockTask(ctrl)
		task2 := mock.NewMockTask(ctrl)

		retryTask.EXPECT().GetName().Return("retryTask").Times(1)
		retryTask.EXPECT().Run(ctx, payloadMock).Return(errors.New("test error")).Times(3)

		task2.EXPECT().GetName().Return("task2").Times(0)
		task2.EXPECT().Run(ctx, payloadMock).Return(nil).Times(0)

		var taskValidator pipeline.TaskValidator = func(n string) bool { return true }
		sr := pipeline.SyncRunner(pipeline.RetryingTask(retryTask, func(err error) bool {
			return true
		}, 3), task2)

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

		task1.EXPECT().GetName().Return("task1")
		task2.EXPECT().GetName().Return("task2")

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

		task1.EXPECT().GetName().Return("task1")
		task2.EXPECT().GetName().Return("task2")

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

		task1.EXPECT().GetName().Return("task1")
		task2.EXPECT().GetName().Return("task2")

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

		task1.EXPECT().GetName().Return("task1")
		task2.EXPECT().GetName().Return("task2")

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

func TestStage_RetryingStageRunner(t *testing.T) {
	t.Run("tasks return success", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		task1 := mock.NewMockTask(ctrl)
		task2 := mock.NewMockTask(ctrl)

		task1.EXPECT().GetName().Return("task1")
		task2.EXPECT().GetName().Return("task2")

		gomock.InOrder(
			task1.EXPECT().Run(ctx, payloadMock).Return(nil),
			task2.EXPECT().Run(ctx, payloadMock).Return(nil),
		)

		var taskValidator pipeline.TaskValidator = func(n string) bool { return true }

		sr := pipeline.RetryingStageRunner(
			pipeline.SyncRunner(task1, task2),
			func(err error) bool { return true },
			3)

		err := sr.Run(ctx, payloadMock, taskValidator)
		if err != nil {
			t.Errorf("should not return error")
		}
	})

	t.Run("task returns error then success", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		task1 := mock.NewMockTask(ctrl)
		task1.EXPECT().GetName().Return("task1").Times(2)

		gomock.InOrder(
			task1.EXPECT().Run(ctx, payloadMock).Return(errors.New("test error")).Times(1),
			task1.EXPECT().Run(ctx, payloadMock).Return(nil).Times(1),
		)

		var taskValidator pipeline.TaskValidator = func(n string) bool { return true }
		sr := pipeline.RetryingStageRunner(
			pipeline.SyncRunner(task1),
			func(err error) bool { return true },
			3)

		err := sr.Run(ctx, payloadMock, taskValidator)
		if err != nil {
			t.Errorf("should not return error")
		}
	})

	t.Run("task returns error every try", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadMock := mock.NewMockPayload(ctrl)

		task1 := mock.NewMockTask(ctrl)
		task1.EXPECT().GetName().Return("task1").Times(3)
		task1.EXPECT().Run(ctx, payloadMock).Return(errors.New("test error")).Times(3)

		var taskValidator pipeline.TaskValidator = func(n string) bool { return true }
		sr := pipeline.RetryingStageRunner(
			pipeline.SyncRunner(task1),
			func(err error) bool { return true },
			3)

		err := sr.Run(ctx, payloadMock, taskValidator)
		if err == nil {
			t.Errorf("should return error")
		}
	})
}
