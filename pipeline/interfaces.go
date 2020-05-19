package pipeline

import "context"

// PayloadFactory is implemented by objects which know how to create payload for every height
type PayloadFactory interface {
	// Gets new payload
	GetPayload() Payload
}

// Payload is implemented by values that can be sent through a pipeline.
type Payload interface {
	// Set current height to be processed
	SetCurrentHeight(int64)

	// Gets current height
	GetCurrentHeight() int64

	// MarkAsProcessed is invoked by the pipeline when the payload
	// reaches the end of execution for current height
	MarkAsProcessed()
}

// Source is executed before processing of individual heights.
// It is responsible for getting start and end height.
type Source interface {
	// Runs stage of type Chore
	Run(context.Context) error

	// Gets start height
	GetStartHeight() int64

	// Gets end height
	GetEndHeight() int64
}

// Sink is executed after all height have been processed
type Sink interface {
	// Runs stage of type Chore
	Run(context.Context) error
}

// TaskValidator is a type for validating task by provided task name
type TaskValidator func(string) bool

// StageRunner is implemented by types that know how to run tasks
type StageRunner interface {
	// Run StageRunner
	Run(context.Context, Payload, TaskValidator) error
}

// StageRunnerFunc is an adapter to allow the use of plain functions as StageRunner
type StageRunnerFunc func(context.Context, Payload, TaskValidator) error

// Run calls f(ctx, p, f).
func (srf StageRunnerFunc) Run(ctx context.Context, p Payload, f TaskValidator) error {
	return srf(ctx, p, f)
}

// Stage is implemented by types which invoke StageRunner
type Stage interface {
	// Run Stage
	Run(context.Context, Payload, *Options) error
}

// Task is implemented by types that want to be executed inside of a stage
type Task interface {
	//Run Task
	Run(context.Context, Payload) error
}

// TaskFunc is an adapter to allow the use of plain functions as Task
type TaskFunc func(context.Context, Payload) error

// Process calls f(ctx, p).
func (f TaskFunc) Run(ctx context.Context, p Payload) error {
	return f(ctx, p)
}
