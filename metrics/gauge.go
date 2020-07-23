package metrics

type Gauge interface {
	Counter

	Set(float64)

	Dec()
	Sub(float64)
}

func NewGaugeWithTags(opts Options) (*GroupTagGauge, error) {
	return DefaultMetrics.NewGaugeWithTags(opts)
}

func MustNewGaugeWithTags(opts Options) *GroupTagGauge {
	g, err := DefaultMetrics.NewGaugeWithTags(opts)
	if err != nil {
		panic("Metric Error: " + err.Error())
	}
	return g
}

type TagGauge interface {
	WithTags(map[string]string) (Gauge, error)
	WithLabels(...string) Gauge
}

type GroupTagGauge struct {
	taggauges []TagGauge
	options   Options
}

func (gtg *GroupTagGauge) AddGauge(g TagGauge) {
	gtg.taggauges = append(gtg.taggauges, g)
}

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

func (gtg *GroupTagGauge) WithLabels(labels ...string) *GroupGauge {
	gc := &GroupGauge{}
	for _, tc := range gtg.taggauges {
		c := tc.WithLabels(labels...)
		gc.AddGauge(c)
	}

	return gc
}

type GroupGauge struct {
	gauges []Gauge
}

func (gg *GroupGauge) Inc() {
	for _, c := range gg.gauges {
		c.Inc()
	}
}

func (gg *GroupGauge) Add(a float64) {
	for _, c := range gg.gauges {
		c.Add(a)
	}
}

func (gg *GroupGauge) Dec() {
	for _, c := range gg.gauges {
		c.Dec()
	}
}

func (gg *GroupGauge) Sub(a float64) {
	for _, c := range gg.gauges {
		c.Sub(a)
	}
}

func (gg *GroupGauge) Set(a float64) {
	for _, c := range gg.gauges {
		c.Set(a)
	}
}

func (gg *GroupGauge) AddGauge(g Gauge) {
	gg.gauges = append(gg.gauges, g)
}
