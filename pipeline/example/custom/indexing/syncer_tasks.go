package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	SyncerTaskName = "SyncerExample"
)

func NewSyncerTask() pipeline.Task {
	return &SyncerTask{}
}

type SyncerTask struct {
}

func (t *SyncerTask) GetName() string {
	return SyncerTaskName
}

func (t *SyncerTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := (p).(*payload)
	fmt.Println("task: ", t.GetName(), payload.CurrentHeight)
	return nil
}
