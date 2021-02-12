package pipeline

import (
	"context"
	"strings"
	"sync"

	"github.com/hashicorp/go-multierror"

	"github.com/figment-networks/indexing-engine/metrics"
)

var (
	_ Stage = (*stage)(nil)
)

// NewStageWithTasks creates a stage with tasks that will run one by one
func NewStageWithTasks(name StageName, tasks ...Task) *stage {
	return &stage{
		Name:   name,
		runner: syncRunner{tasks},
	}
}

// NewAsyncStageWithTasks creates a stage with tasks that will run concurrently
func NewAsyncStageWithTasks(name StageName, tasks ...Task) *stage {
	return &stage{
		Name:   name,
		runner: asyncRunner{tasks},
	}
}

// NewCustomStage creates a stage with custom stagerunner
func NewCustomStage(name StageName, runner stageRunner) *stage {
	return &stage{
		Name:   name,
		runner: runner,
	}
}

type stage struct {
	Name   StageName
	runner stageRunner
}

// Run runs the stage runner assigned to stage
func (s *stage) Run(ctx context.Context, payload Payload, options *Options) error {
	observer := stageDurationMetric.WithLabels(string(s.Name))

	timer := metrics.NewTimer(observer)
	defer timer.ObserveDuration()

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

// runTask executes a pipeline task
func runTask(ctx context.Context, task Task, payload Payload) error {
	observer := taskDurationMetric.WithLabels(task.GetName())

	timer := metrics.NewTimer(observer)
	defer timer.ObserveDuration()

	return task.Run(ctx, payload)
}

type syncRunner struct {
	tasks []Task
}

// Run runs syncRunner
func (r syncRunner) Run(ctx context.Context, payload Payload, canRunTask TaskValidator) error {
	for _, task := range r.tasks {
		if canRunTask(task.GetName()) {
			err := runTask(ctx, task, payload)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
				if err := runTask(ctx, task, payload); err != nil {
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

// retryingStageRunner implement retry mechanism for stageRunner
func retryingStageRunner(sr stageRunner, isTransient func(error) bool, maxRetries int) stageRunner {
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
		if err = runTask(ctx, r.task, p); err != nil {
			if !r.isTransient(err) {
				return err
			}
		} else {
			break
		}
	}
	return err
}

// RetryingTask implements retry mechanism for Task
func RetryingTask(st Task, isTransient func(error) bool, maxRetries int) Task {
	return &retryTask{
		name:        st.GetName(),
		task:        st,
		isTransient: isTransient,
		maxRetries:  maxRetries,
	}
}
