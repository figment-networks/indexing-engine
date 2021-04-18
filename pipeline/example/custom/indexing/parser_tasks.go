package indexing

import (
	"context"
	"fmt"

	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	ParserTaskName1 = "ParserExample (Task 1)"
	ParserTaskName2 = "ParserExample (Task 2)"
)

func NewParserTask() pipeline.Task {
	return &ParserTask{
		name: ParserTaskName1,
	}
}

func NewParserTask2() pipeline.Task {
	return &ParserTask{
		name: ParserTaskName2,
	}
}

type ParserTask struct {
	name string
}

func (t *ParserTask) GetName() string {
	return t.name
}

func (t *ParserTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := (p).(*payload)
	fmt.Println("task: ", t.GetName(), payload.CurrentHeight)
	return nil
}
