package metrics

type Observer interface {
	Observe(float64)
}
type HistogramBucketOptions struct {
	Type    string
	Buckets []string
}
type HistogramOptions struct {
	Namespace string
	Subsystem string
	Name      string
	Desc      string
	Tags      []string
	Buckets   HistogramBucketOptions
}

func NewHistogramWithTags(opts HistogramOptions) (*GroupTagHistogram, error) {
	return DefaultMetrics.NewHistogramWithTags(opts)
}

func MustNewHistogramWithTags(opts HistogramOptions) *GroupTagHistogram {
	g, err := DefaultMetrics.NewHistogramWithTags(opts)
	if err != nil {
		panic("Metric Error: " + err.Error())
	}
	return g
}

type TagObserver interface {
	WithTags(map[string]string) (Observer, error)
	WithLabels([]string) Observer
}

type GroupTagHistogram struct {
	tagobservers []TagObserver
	options      HistogramOptions
}

func (gth *GroupTagHistogram) AddHistogram(o TagObserver) {
	gth.tagobservers = append(gth.tagobservers, o)
}

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

func (gth *GroupTagHistogram) WithLabels(labels []string) *GroupObserver {
	gh := &GroupObserver{}
	for _, tc := range gth.tagobservers {
		c := tc.WithLabels(labels)
		gh.AddObserver(c)
	}

	return gh
}

type GroupObserver struct {
	observers []Observer
}

func (g *GroupObserver) AddObserver(o Observer) {
	g.observers = append(g.observers, o)
}

func (g *GroupObserver) Observe(f float64) {
	for _, c := range g.observers {
		c.Observe(f)
	}
}
