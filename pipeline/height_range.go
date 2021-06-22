package pipeline

import "errors"

var (
	// ErrInvalidLatestHeight is returned when the latest height is a negative number
	ErrInvalidLatestHeight = errors.New("lastest height is invalid")

	// ErrInvalidInitialHeight is returned when the initial height is a negative number
	ErrInvalidInitialHeight = errors.New("initial height is invalid")

	// ErrInvalidBatchSize is returned when the batch size is a negative number
	ErrInvalidBatchSize = errors.New("batch size is invalid")

	// ErrNothingToProcess is returned when there are no heights to process
	ErrNothingToProcess = errors.New("nothing to process")
)

// HeightRange represents a range of heights to be processed
type HeightRange struct {
	// The most recent height
	LatestHeight int64

	// The last processed height
	LastHeight int64

	// The starting height when no height has been processed yet
	InitialHeight int64

	// The number of heights processed in one run
	BatchSize int64
}

// StartHeight calculates the first height to process
func (hr *HeightRange) StartHeight() int64 {
	if hr.LastHeight < 0 {
		return hr.InitialHeight
	}

	return hr.LastHeight + 1
}

// EndHeight calculates the last height to process
func (hr *HeightRange) EndHeight() int64 {
	if hr.BatchSize == 0 {
		return hr.LatestHeight
	}

	batchEnd := hr.StartHeight() + hr.BatchSize - 1

	if batchEnd > hr.LatestHeight {
		return hr.LatestHeight
	}
	return batchEnd
}

// Length calculates the number of heights to process
func (hr *HeightRange) Length() int64 {
	return hr.EndHeight() - hr.StartHeight() + 1
}

// Validate checks if the height range is valid
func (hr *HeightRange) Validate(checkLength bool) error {
	if hr.LatestHeight < 0 {
		return ErrInvalidLatestHeight
	}

	if hr.InitialHeight < 0 {
		return ErrInvalidInitialHeight
	}

	if hr.BatchSize < 0 {
		return ErrInvalidBatchSize
	}

	if checkLength && hr.Length() <= 0 {
		return ErrNothingToProcess
	}

	return nil
}
