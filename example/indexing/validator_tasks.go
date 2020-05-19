package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine.git/pipeline"
	"reflect"
)

func NewValidatorTask() pipeline.Task {
	return &ValidatorTask{
	}
}

type ValidatorTask struct {
}

func (f *ValidatorTask) Run(ctx context.Context,  payload pipeline.Payload) error {
	fmt.Println("task: ", reflect.TypeOf(*f).Name(), payload.GetCurrentHeight())
	return nil
}

