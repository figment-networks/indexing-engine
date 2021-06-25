package indexing

import (
	"context"
	"fmt"

	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	AggregatorTaskName = "AggregatorExample"
)

func NewAggregatorTask() pipeline.Task {
	return &AggregatorTask{}
}

type AggregatorTask struct {
}

func (t *AggregatorTask) GetName() string {
	return AggregatorTaskName
}

func (t *AggregatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := (p).(*payload)
	fmt.Println("task: ", t.GetName(), payload.CurrentHeight)
	return nil
}
