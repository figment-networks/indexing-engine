package pipeline

import "context"

// PayloadFactory is implemented by objects which know how to create payload for every height
type PayloadFactory interface {
	// Gets new payload
	GetPayload(int64) Payload
}

// Payload is implemented by values that can be sent through a pipeline.
type Payload interface {
	// MarkAsProcessed is invoked by the pipeline when the payload
	// reaches the end of execution for current height
	MarkAsProcessed()
}

// Source is executed before processing of individual heights.
// It is responsible for getting start and end height.
type Source interface {
	// Next gets next height
	Next(context.Context, Payload) bool

	// Current returns current height
	Current() int64

	// Err return error if any
	Err() error
}

// Sink is executed as a last stage in the pipeline
type Sink interface {
	// Consume consumes payload
	Consume(context.Context, Payload) error
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

	// GetName gets name of task
	GetName() string
}

