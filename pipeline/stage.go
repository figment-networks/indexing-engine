package pipeline

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"reflect"
	"strings"
	"sync"
)

var (
	_ Stage = (*stage)(nil)
)

func NewStage(name StageName, runner StageRunner, stageType StageType) *stage {
	return &stage{
		Name:      name,
		runner:    runner,
		stageType: stageType,
	}
}

type StageType int64

const (
	StageTypeChore StageType = iota
	StageTypeSyncing
	StageTypeIndexing

	// StageTypeCustom it is a stage added dynamically by the user
	StageTypeCustom
)

type stage struct {
	Name      StageName
	runner    StageRunner
	stageType StageType
}

// Run runs the stage runner assigned to stage
func (s *stage) Run(ctx context.Context, payload Payload, options *Options) error {
	return s.runner.Run(ctx, payload, func(taskName string) bool {
		return s.canRunTask(taskName, options)
	})
}

// canRunTask determines if task can be ran
func (s *stage) canRunTask(taskName string, options *Options) bool {
	if s.stageType == StageTypeIndexing && options != nil && len(options.IndexingTasksWhitelist) > 0 {
		for _, t := range options.IndexingTasksWhitelist {
			if strings.Contains(taskName, t) {
				return true
			}
		}
		return false
	}

	if s.stageType == StageTypeIndexing && options != nil && len(options.IndexingTasksBlacklist) > 0 {
		for _, t := range options.IndexingTasksBlacklist {
			if strings.Contains(taskName, t) {
				return false
			}
		}
		return true
	}

	return true
}

// SyncRunner runs tasks one by one
func SyncRunner(tasks ...Task) StageRunner {
	return syncRunner{tasks: tasks}
}

type syncRunner struct {
	tasks []Task
}

// Run runs syncRunner
func (r syncRunner) Run(ctx context.Context, payload Payload, canRunTask TaskValidator) error {
	for _, t := range r.tasks {
		if canRunTask(reflect.TypeOf(t).String()) {
			err := t.Run(ctx, payload)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AsyncRunner runs tasks concurrently
func AsyncRunner(tasks ...Task) StageRunner {
	return asyncRunner{tasks: tasks}
}

type asyncRunner struct {
	tasks []Task
}

// Run runs AsyncRunner
func (ar asyncRunner) Run(ctx context.Context, payload Payload, canRunTask TaskValidator) error {
	var wg sync.WaitGroup
	var errs error
	errCh := make(chan error, len(ar.tasks))
	for _, task := range ar.tasks {
		if canRunTask(reflect.TypeOf(task).String()) {
			wg.Add(1)
			go func(task Task, ctx context.Context, payload Payload) {
				if err := task.Run(ctx, payload); err != nil {
					errCh <- err
				}
				wg.Done()
			}(task, ctx, payload)
		}
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

// RetryingStageRunner implement retry mechanism for StageRunner
func RetryingStageRunner(sr StageRunner, isTransient func(error) bool, maxRetries int) StageRunner {
	return StageRunnerFunc(func(ctx context.Context, p Payload, f TaskValidator) error {
		var err error
		for i := 0; i < maxRetries; i++ {
			if err = sr.Run(ctx, p, f); err != nil {
				if !isTransient(err) {
					return err
				}
			} else {
				break
			}
		}
		return err
	})
}

// RetryingStageRunner implement retry mechanism for Task
func RetryingTask(st Task, isTransient func(error) bool, maxRetries int) Task {
	return TaskFunc(func(ctx context.Context, p Payload) error {
		var err error
		for i := 0; i < maxRetries; i++ {
			if err = st.Run(ctx, p); err != nil {
				if !isTransient(err) {
					return err
				}
			} else {
				break
			}
		}
		return err
	})
}
