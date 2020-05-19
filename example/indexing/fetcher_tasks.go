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

func (f *FetcherTask) Run(ctx context.Context, payload pipeline.Payload) error {
	fmt.Println("task: ", reflect.TypeOf(*f).Name(), payload.GetCurrentHeight())
	return nil
}
