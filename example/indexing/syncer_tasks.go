package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine.git/pipeline"
	"reflect"
)

func NewSyncerTask() pipeline.Task {
	return &SyncerTask{
	}
}

type SyncerTask struct {
}

func (f *SyncerTask) Run(ctx context.Context,  p pipeline.Payload) error {
	payload := (p).(*payload)
	fmt.Println("task: ", reflect.TypeOf(*f).Name(), payload.currentHeight)
	return nil
}
