package metrics

type Counter interface {
	Inc()
	Add(float64)
}

func NewCounterWithTags(opts Options) (*GroupTagCounter, error) {
	return DefaultMetrics.NewCounterWithTags(opts)
}

func MustNewCounterWithTags(opts Options) *GroupTagCounter {
	return DefaultMetrics.MustNewCounterWithTags(opts)
}

type TagCounter interface {
	WithTags(map[string]string) (Counter, error)
	WithLabels([]string) Counter
}

type GroupTagCounter struct {
	tagcounters []TagCounter
	options     Options
}

func (gtc *GroupTagCounter) AddCounter(c TagCounter) {
	gtc.tagcounters = append(gtc.tagcounters, c)
}

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

func (gtc *GroupTagCounter) WithLabels(labels []string) *GroupCounter {
	gc := &GroupCounter{}
	for _, tc := range gtc.tagcounters {
		c := tc.WithLabels(labels)
		gc.AddCounter(c)
	}

	return gc
}

type GroupCounter struct {
	counters []Counter
}

func (gc *GroupCounter) AddCounter(c Counter) {
	gc.counters = append(gc.counters, c)
}

func (gc *GroupCounter) Inc() {
	for _, c := range gc.counters {
		c.Inc()
	}
}

func (gc *GroupCounter) Add(a float64) {
	for _, c := range gc.counters {
		c.Add(a)
	}
}
