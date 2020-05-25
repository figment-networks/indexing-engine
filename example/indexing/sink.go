package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
)

func NewSink() pipeline.Sink {
	return &sink{}
}

type sink struct {}

func (s sink) Consume(ctx context.Context, p pipeline.Payload) error {
	fmt.Println("sink")
	return nil
}