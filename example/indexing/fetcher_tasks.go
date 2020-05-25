package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine.git/pipeline"
	"reflect"
)

func NewFetcherTask() pipeline.Task {
	return &FetcherTask{
	}
}

type FetcherTask struct {
}

func (f *FetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := (p).(*payload)
	fmt.Println("task: ", reflect.TypeOf(*f).Name(), payload.currentHeight)
	return nil
}
