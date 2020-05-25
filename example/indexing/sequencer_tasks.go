package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
	"reflect"
)

func NewSequencerTask() pipeline.Task {
	return &SequencerTask{
	}
}

type SequencerTask struct {
}

func (f *SequencerTask) Run(ctx context.Context,  p pipeline.Payload) error {
	payload := p.(*payload)
	fmt.Println("task: ", reflect.TypeOf(*f).Name(), payload.currentHeight)
	return nil
}
