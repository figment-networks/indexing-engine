package worker

import (
	"math"
	"time"
)

const _backoffFactor = 2

// Backoff implements an exponential backoff algorithm
type Backoff struct {
	attempts uint64
}

// Attempt increments the number of attempts
func (b *Backoff) Attempt() {
	b.attempts++
}

// Reset resets the number of attempts
func (b *Backoff) Reset() {
	b.attempts = 0
}

// Delay calculates the backoff time
func (b Backoff) Delay() time.Duration {
	if b.attempts == 0 {
		return time.Duration(0)
	}

	backoff := math.Pow(_backoffFactor, float64(b.attempts))

	return time.Duration(backoff) * time.Second
}
