package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine.git/pipeline"
	"reflect"
)

func NewSequencerTask() pipeline.Task {
	return &SequencerTask{
	}
}

type SequencerTask struct {
}

func (f *SequencerTask) Run(ctx context.Context,  payload pipeline.Payload) error {
	fmt.Println("task: ", reflect.TypeOf(*f).Name(), payload.GetCurrentHeight())
	return nil
}
