package metrics

import "time"

// Timer This one is taken from prometheus impementation to get rid out of dependancy
type Timer struct {
	begin    time.Time
	observer Observer
}

// NewTimer is a constructor for Timer, counting from time.Now()
func NewTimer(o Observer) *Timer {
	return &Timer{
		begin:    time.Now(),
		observer: o,
	}
}

// ObserveDuration is commiting time duration that passed from it's creation
func (t *Timer) ObserveDuration() time.Duration {
	d := time.Since(t.begin)
	if t.observer != nil {
		t.observer.Observe(d.Seconds())
	}
	return d
}
