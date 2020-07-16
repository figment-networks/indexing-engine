package indexing

import (
	"context"
	"fmt"

	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	ValidatorTaskName1 = "ValidatorExample (Task 1)"
	ValidatorTaskName2 = "ValidatorExample (Task 2)"
)

func NewValidatorTask() pipeline.Task {
	return &ValidatorTask{
		name: ValidatorTaskName1,
	}
}

func NewValidatorTask2() pipeline.Task {
	return &ValidatorTask{
		name: ValidatorTaskName2,
	}
}

type ValidatorTask struct {
	name string
}

func (t *ValidatorTask) GetName() string {
	return t.name
}

func (t *ValidatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := (p).(*payload)
	fmt.Println("task: ", t.GetName(), payload.CurrentHeight)
	return nil
}
