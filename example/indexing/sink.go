package indexing

import (
	"context"
	"errors"
	"fmt"
	"github.com/figment-networks/indexing-engine.git/pipeline"
)

func NewSink() pipeline.Sink {
	return &sink{}
}

type sink struct {}

func (s sink) Run(ctx context.Context) error {
	statRecorder, ok := ctx.Value(pipeline.CtxStats).(*pipeline.StatsRecorder)
	if !ok {
		return errors.New("statrecorder not recognized")
	}
	statRecorder.SetCompleted(true)

	// Crate report from stats
	fmt.Println("sink")
	fmt.Println("completed at: ", statRecorder.Duration)
	return nil
}