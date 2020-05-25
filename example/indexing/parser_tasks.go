package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine.git/pipeline"
	"reflect"
)

func NewParserTask() pipeline.Task {
	return &ParserTask{
	}
}

type ParserTask struct {
}

func (f *ParserTask) Run(ctx context.Context,  p pipeline.Payload) error {
	payload := (p).(*payload)
	fmt.Println("task: ", reflect.TypeOf(*f).Name(), payload.currentHeight)
	return nil
}
