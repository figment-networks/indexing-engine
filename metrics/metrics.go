package metrics

import (
	"fmt"
)

var DetaultMetrics *Metrics

func init() {
	DetaultMetrics = NewMetrics()
}

type MetricsEngine interface {
	Name() string

	NewCounterWithTags(opts Options) (TagCounter, error)
	NewGaugeWithTags(opts Options) (TagGauge, error)
	NewHistogramWithTags(opts HistogramOptions) (TagObserver, error)
}

type Options struct {
	Namespace string
	Subsystem string
	Name      string
	Desc      string
	Tags      []string
}

type Metrics struct {
	engines map[string]MetricsEngine
}

func NewMetrics() *Metrics {
	return &Metrics{engines: make(map[string]MetricsEngine)}
}

func (m *Metrics) AddEngine(eng MetricsEngine) error {
	_, ok := m.engines[eng.Name()]
	if ok {
		return fmt.Errorf("Engine %s is already registred")
	}
	m.engines[eng.Name()] = eng
	return nil
}

func (m *Metrics) NewCounterWithTags(opts Options) (*GroupTagCounter, error) {
	gc := &GroupTagCounter{}
	for _, e := range m.engines {
		c, err := e.NewCounterWithTags(opts)
		if err != nil {
			return nil, err
		}
		gc.AddCounter(c)

	}
	return gc, nil
}

func (m *Metrics) NewGaugeWithTags(opts Options) (*GroupTagGauge, error) {
	gc := &GroupTagGauge{}
	for _, e := range m.engines {
		c, err := e.NewGaugeWithTags(opts)
		if err != nil {
			return nil, err
		}
		gc.AddGauge(c)

	}
	return gc, nil
}

func (m *Metrics) NewHistogramWithTags(opts HistogramOptions) (*GroupTagHistogram, error) {
	gc := &GroupTagHistogram{}
	for _, e := range m.engines {
		c, err := e.NewHistogramWithTags(opts)
		if err != nil {
			return nil, err
		}
		gc.AddHistogram(c)

	}
	return gc, nil
}
