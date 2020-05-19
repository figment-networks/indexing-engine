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

func (f *SyncerTask) Run(ctx context.Context,  payload pipeline.Payload) error {
	fmt.Println("task: ", reflect.TypeOf(*f).Name(), payload.GetCurrentHeight())
	return nil
}
