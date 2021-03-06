package pipeline

import "context"

func NewSource() Source {
	return &source{
		currentHeight: 10,
		startHeight:   10,
		endHeight:     11,
	}
}

type source struct {
	startHeight   int64
	currentHeight int64
	endHeight     int64
	err           error
}

func (s *source) Next(ctx context.Context, p Payload) bool {
	if s.err == nil && s.currentHeight < s.endHeight {
		s.currentHeight = s.currentHeight + 1
		return true
	}
	return false
}

func (s *source) Current() int64 {
	return s.currentHeight
}

func (s *source) Err() error {
	return s.err
}

func (s *source) Skip(stageName StageName) bool {
	return false
}
