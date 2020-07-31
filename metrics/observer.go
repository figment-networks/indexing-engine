package metrics

// Observer is an interface that generalizes a group of metrics
// that are receiving a set of numeric values in time
type Observer interface {
	Observe(float64)
}

// HistogramBucketOptions an options to fine-tune histogram buckets
type HistogramBucketOptions struct {
	Type    string
	Buckets []string
}

// HistogramOptions a set of options for histogram metric
type HistogramOptions struct {
	Namespace string
	Subsystem string
	Name      string
	Desc      string
	Tags      []string
	Buckets   HistogramBucketOptions
}

// NewHistogramWithTags Histogram constructor for DefaultMetrics global
func NewHistogramWithTags(opts HistogramOptions) (*GroupTagHistogram, error) {
	return DefaultMetrics.NewHistogramWithTags(opts)
}

// MustNewHistogramWithTags constructor with embedded panic
func MustNewHistogramWithTags(opts HistogramOptions) *GroupTagHistogram {
	g, err := DefaultMetrics.NewHistogramWithTags(opts)
	if err != nil {
		panic("Metric Error: " + err.Error())
	}
	return g
}

// TagObserver interface for appending tags to metrics
type TagObserver interface {
	WithTags(map[string]string) (Observer, error)
	WithLabels(...string) Observer
}

// GroupTagHistogram a group of different Histograms from different metric systems
type GroupTagHistogram struct {
	tagobservers []TagObserver
	options      HistogramOptions
}

// AddHistogram is appending new histogram to set
func (gth *GroupTagHistogram) AddHistogram(o TagObserver) {
	gth.tagobservers = append(gth.tagobservers, o)
}

// WithTags makes a group with given tags (label-value pairs)
func (gth *GroupTagHistogram) WithTags(tags map[string]string) (*GroupObserver, error) {
	gh := &GroupObserver{}
	for _, tc := range gth.tagobservers {
		c, err := tc.WithTags(tags)
		if err != nil {
			return nil, err
		}
		gh.AddObserver(c)
	}

	return gh, nil
}

// WithLabels make a group with given labels
func (gth *GroupTagHistogram) WithLabels(labels ...string) *GroupObserver {
	gh := &GroupObserver{}
	for _, tc := range gth.tagobservers {
		c := tc.WithLabels(labels...)
		gh.AddObserver(c)
	}

	return gh
}

// GroupObserver is a set of Observers that within the same tags
type GroupObserver struct {
	observers []Observer
}

// AddObserver adds observer
func (g *GroupObserver) AddObserver(o Observer) {
	g.observers = append(g.observers, o)
}

// Observe calls the metrics ultimately
func (g *GroupObserver) Observe(f float64) {
	for _, c := range g.observers {
		c.Observe(f)
	}
}
