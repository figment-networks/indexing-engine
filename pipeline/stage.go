package pipeline

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"strings"
	"sync"
)

var (
	_ Stage = (*stage)(nil)
)

func NewStage(name StageName, runner StageRunner) *stage {
	return &stage{
		Name:      name,
		runner:    runner,
	}
}

type stage struct {
	Name      StageName
	runner    StageRunner
}

// Run runs the stage runner assigned to stage
func (s *stage) Run(ctx context.Context, payload Payload, options *Options) error {
	return s.runner.Run(ctx, payload, func(taskName string) bool {
		return s.canRunTask(taskName, options)
	})
}

// canRunTask determines if task can be ran
func (s *stage) canRunTask(taskName string, options *Options) bool {
	if options != nil && len(options.TaskWhitelist) > 0 {
		for _, t := range options.TaskWhitelist {
			if strings.Contains(taskName, string(t)) {
				return true
			}
		}
		return false
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
	for _, task := range r.tasks {
		if canRunTask(task.GetName()) {
			err := task.Run(ctx, payload)
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
		if canRunTask(task.GetName()) {
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

// retryTask is a task with built-in retry mechanism
type retryTask struct {
	name        string
	task        Task
	isTransient func(error) bool
	maxRetries  int
}

// GetName get the name of retry task. It is the same as the original task name
func (r *retryTask) GetName() string {
	return r.name
}

// Run runs retry task
func (r *retryTask) Run(ctx context.Context, p Payload) error {
	var err error
	for i := 0; i < r.maxRetries; i++ {
		if err = r.task.Run(ctx, p); err != nil {
			if !r.isTransient(err) {
				return err
			}
		} else {
			break
		}
	}
	return err
}

// RetryingTask implement retry mechanism for Task
func RetryingTask(st Task, isTransient func(error) bool, maxRetries int) Task {
	return &retryTask{
		name:        st.GetName(),
		task:        st,
		isTransient: isTransient,
		maxRetries:  maxRetries,
	}
}
