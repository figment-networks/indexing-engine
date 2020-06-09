package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	ParserTaskName = "ParserExample"
)

func NewParserTask() pipeline.Task {
	return &ParserTask{
	}
}

type ParserTask struct {
}

func (t *ParserTask) GetName() string {
	return ParserTaskName
}

func (t *ParserTask) Run(ctx context.Context,  p pipeline.Payload) error {
	payload := (p).(*payload)
	fmt.Println("task: ", t.GetName(), payload.CurrentHeight)
	return nil
}
