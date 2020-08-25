package metrics

// Counter is an interface that generalizes a group of metrics
// that are calculate the number over time
type Counter interface {
	Inc()
	Add(float64)
}

// NewCounterWithTags is Counter constructor for DefaultMetrics global
func NewCounterWithTags(opts Options) (*GroupTagCounter, error) {
	return DefaultMetrics.NewCounterWithTags(opts)
}

// MustNewCounterWithTags constructor with embedded panic
func MustNewCounterWithTags(opts Options) *GroupTagCounter {
	return DefaultMetrics.MustNewCounterWithTags(opts)
}

// TagCounter interface for appending tags to metrics
type TagCounter interface {
	WithTags(map[string]string) (Counter, error)
	WithLabels(...string) Counter
}

// GroupTagCounter a group of different Counters from different metric systems
type GroupTagCounter struct {
	tagcounters []TagCounter
	options     Options

	registred map[uint64]*GroupCounter
}

// AddCounter is appending new counter to set
func (gtc *GroupTagCounter) AddCounter(c TagCounter) {
	gtc.tagcounters = append(gtc.tagcounters, c)
}

// WithTags makes a group with given tags (label-value pairs)
func (gtc *GroupTagCounter) WithTags(tags map[string]string) (*GroupCounter, error) {

	gc := &GroupCounter{}
	for _, tc := range gtc.tagcounters {
		c, err := tc.WithTags(tags)
		if err != nil {
			return nil, err
		}
		gc.AddCounter(c)
	}

	return gc, nil
}

// WithLabels make a group with given labels
func (gtc *GroupTagCounter) WithLabels(labels ...string) *GroupCounter {
	h := secureHash.GetHash(labels)
	gc, ok := gtc.registred[h]
	if ok {
		return gc
	}

	gc = &GroupCounter{}
	for _, tc := range gtc.tagcounters {
		c := tc.WithLabels(labels...)
		gc.AddCounter(c)
	}

	gtc.registred[h] = gc
	return gc
}

// GroupCounter is a set of Counters that within the same tags
type GroupCounter struct {
	counters []Counter
}

// AddCounter adds counter
func (gc *GroupCounter) AddCounter(c Counter) {
	gc.counters = append(gc.counters, c)
}

// Inc increases the counter
func (gc *GroupCounter) Inc() {
	for _, c := range gc.counters {
		c.Inc()
	}
}

// Add adds to counter
func (gc *GroupCounter) Add(a float64) {
	for _, c := range gc.counters {
		c.Add(a)
	}
}
