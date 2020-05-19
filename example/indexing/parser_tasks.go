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

func (f *ParserTask) Run(ctx context.Context,  payload pipeline.Payload) error {
	fmt.Println("task: ", reflect.TypeOf(*f).Name(), payload.GetCurrentHeight())
	return nil
}
