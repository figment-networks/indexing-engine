package metrics

// Gauge is an interface that generalizes a group of metrics
// that are used to report a curent level/number
type Gauge interface {
	Counter

	Set(float64)

	Dec()
	Sub(float64)
}

// NewGaugeWithTags gauge type constructor
func NewGaugeWithTags(opts Options) (*GroupTagGauge, error) {
	return DefaultMetrics.NewGaugeWithTags(opts)
}

// MustNewGaugeWithTags gauge constructor with embedded error
func MustNewGaugeWithTags(opts Options) *GroupTagGauge {
	g, err := DefaultMetrics.NewGaugeWithTags(opts)
	if err != nil {
		panic("Metric Error: " + err.Error())
	}
	return g
}

// TagGauge interface for appending tags to gauge metrics
type TagGauge interface {
	WithTags(map[string]string) (Gauge, error)
	WithLabels(...string) Gauge
}

// GroupTagGauge a group of different Gauges from different metric systems
type GroupTagGauge struct {
	taggauges []TagGauge
	options   Options
}

// AddGauge is appending new gauge to set
func (gtg *GroupTagGauge) AddGauge(g TagGauge) {
	gtg.taggauges = append(gtg.taggauges, g)
}

// WithTags makes a group with given tags (label-value pairs)
func (gtg *GroupTagGauge) WithTags(tags map[string]string) (*GroupGauge, error) {
	gc := &GroupGauge{}
	for _, tc := range gtg.taggauges {
		c, err := tc.WithTags(tags)
		if err != nil {
			return nil, err
		}
		gc.AddGauge(c)
	}

	return gc, nil
}

// WithLabels makes a group with given labels
func (gtg *GroupTagGauge) WithLabels(labels ...string) *GroupGauge {
	gc := &GroupGauge{}
	for _, tc := range gtg.taggauges {
		c := tc.WithLabels(labels...)
		gc.AddGauge(c)
	}

	return gc
}

// GroupGauge is a set of Gauges that within the same tags
type GroupGauge struct {
	gauges []Gauge
}

// AddGauge adds gauge
func (gg *GroupGauge) AddGauge(g Gauge) {
	gg.gauges = append(gg.gauges, g)
}

// Inc increases gauge
func (gg *GroupGauge) Inc() {
	for _, c := range gg.gauges {
		c.Inc()
	}
}

// Dec decreses gauge
func (gg *GroupGauge) Dec() {
	for _, c := range gg.gauges {
		c.Dec()
	}
}

// Add adds to gauge
func (gg *GroupGauge) Add(a float64) {
	for _, c := range gg.gauges {
		c.Add(a)
	}
}

// Sub substract from gauge
func (gg *GroupGauge) Sub(a float64) {
	for _, c := range gg.gauges {
		c.Sub(a)
	}
}

// Set sets gauge to given value
func (gg *GroupGauge) Set(a float64) {
	for _, c := range gg.gauges {
		c.Set(a)
	}
}
