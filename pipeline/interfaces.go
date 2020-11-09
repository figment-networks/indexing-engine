package pipeline

import "context"

// PayloadFactory is implemented by objects which know how to create payloadMock for every height
type PayloadFactory interface {
	// Gets new payloadMock
	GetPayload(int64) Payload
}

// Payload is implemented by values that can be sent through a pipeline.
type Payload interface {
	// MarkAsProcessed is invoked by the pipeline when the payloadMock
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

	// Skip return bool to skip stage
	Skip(StageName) bool
}

// Sink is executed as a last stage in the pipeline
type Sink interface {
	// Consume consumes payloadMock
	Consume(context.Context, Payload) error
}

// TaskValidator is a type for validating task by provided task name
type TaskValidator func(string) bool

// stageRunner is implemented by types that know how to run tasks
type stageRunner interface {
	// Run stageRunner
	Run(context.Context, Payload, TaskValidator) error
}

// StageRunnerFunc is an adapter to allow the use of plain functions as stageRunner
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

// Logger is implemented by types that want to hook up to logging mechanism in engine
type Logger interface {
	// Info logs info message
	Info(string)

	// Debug logs debug message
	Debug(string)
}
