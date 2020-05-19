package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine.git/pipeline"
	"reflect"
)

func NewAggregatorTask() pipeline.Task {
	return &AggregatorTask{
	}
}

type AggregatorTask struct {
}

func (f *AggregatorTask) Run(ctx context.Context,  payload pipeline.Payload) error {
	fmt.Println("task: ", reflect.TypeOf(*f).Name(), payload.GetCurrentHeight())
	return nil
}
