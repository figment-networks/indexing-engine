package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	FetcherTaskName = "FetcherExample"
)

func NewFetcherTask() pipeline.Task {
	return &FetcherTask{
	}
}

type FetcherTask struct {
}

func (t *FetcherTask) GetName() string {
	return FetcherTaskName
}

func (t *FetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := (p).(*payload)
	fmt.Println("task: ", t.GetName(), payload.CurrentHeight)
	return nil
}
