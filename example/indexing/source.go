package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine.git/pipeline"
)

func NewSource() pipeline.Source {
	return &source{}
}

type source struct {
	startHeight int64
	endHeight int64
}

func (s *source) GetStartHeight() int64 {
	return s.startHeight
}

func (s *source) GetEndHeight() int64 {
	return s.endHeight
}

func (s *source) Run(ctx context.Context) error {
	// Get start and end heights
	s.startHeight = 10
	s.endHeight = 12

	fmt.Println("source")

	return nil
}

