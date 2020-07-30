package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	SequencerTaskName = "SequencerExample"
)

func NewSequencerTask() pipeline.Task {
	return &SequencerTask{
	}
}

type SequencerTask struct {
}

func (t *SequencerTask) GetName() string {
	return SequencerTaskName
}

func (t *SequencerTask) Run(ctx context.Context,  p pipeline.Payload) error {
	payload := p.(*payload)
	fmt.Println("task: ", t.GetName(), payload.CurrentHeight)
	return nil
}
