package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	ValidatorTaskName = "ValidatorExample"
)

func NewValidatorTask() pipeline.Task {
	return &ValidatorTask{
	}
}

type ValidatorTask struct {
}

func (t *ValidatorTask) GetName() string {
	return ValidatorTaskName
}

func (t *ValidatorTask) Run(ctx context.Context,  p pipeline.Payload) error {
	payload := (p).(*payload)
	fmt.Println("task: ", t.GetName(), payload.CurrentHeight)
	return nil
}

