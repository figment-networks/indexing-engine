package metrics

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var DefaultMetrics *Metrics
var Debug bool

func init() {
	DefaultMetrics = NewMetrics()
	_, ok := os.LookupEnv("INDEXING_DEBUG_ENABLE")
	if ok {
		Debug = true
	}
}

func Hotload(name string) error {
	return DefaultMetrics.Hotload(name)
}

func AddEngine(eng MetricsEngine) error {
	return DefaultMetrics.AddEngine(eng)
}

func Handler() http.Handler {
	return DefaultMetrics.Handler()
}

type MetricsEngine interface {
	Name() string

	NewCounterWithTags(opts Options) (TagCounter, error)
	NewGaugeWithTags(opts Options) (TagGauge, error)
	NewHistogramWithTags(opts HistogramOptions) (TagObserver, error)
	Handler() http.Handler
}

type Options struct {
	Namespace string
	Subsystem string
	Name      string
	Desc      string
	Tags      []string
}

type Metrics struct {
	handler *mhandler
	mh      MetricsHandler

	engines map[string]MetricsEngine

	counters  []*GroupTagCounter
	gauges    []*GroupTagGauge
	observers []*GroupTagHistogram
}

func NewMetrics() *Metrics {
	engines := make(map[string]MetricsEngine)
	handle := &mhandler{}
	return &Metrics{handler: handle, mh: MetricsHandler{handler: handle}, engines: engines}
}

func (m *Metrics) AddEngine(eng MetricsEngine) error {
	_, ok := m.engines[eng.Name()]
	if ok {
		return fmt.Errorf("Engine %s is already registred")
	}
	m.engines[eng.Name()] = eng
	return nil
}

func (m *Metrics) Hotload(name string) error {
	eng, ok := m.engines[name]
	if !ok {
		return fmt.Errorf("There is no such engine loaded")
	}
	if Debug {
		log.Printf("Hotloading engine %s %+v", name, eng)
	}

	for _, c := range m.counters {
		counter, err := eng.NewCounterWithTags(c.options)
		if err != nil {
			return err
		}
		c.AddCounter(counter)
		if Debug {
			log.Printf("Added counter %+v", c.options)
		}
	}

	for _, c := range m.gauges {
		gauge, err := eng.NewGaugeWithTags(c.options)
		if err != nil {
			return err
		}
		c.AddGauge(gauge)
		if Debug {
			log.Printf("Added gauge %+v", c.options)
		}
	}

	for _, c := range m.observers {
		observer, err := eng.NewHistogramWithTags(c.options)
		if err != nil {
			return err
		}
		c.AddHistogram(observer)
		if Debug {
			log.Printf("Added histogram %+v", c.options)
		}
	}

	m.handler.Handler = eng.Handler()
	return nil
}

func (m *Metrics) NewCounterWithTags(opts Options) (*GroupTagCounter, error) {
	gc := &GroupTagCounter{
		options: opts,
	}
	for _, e := range m.engines {
		c, err := e.NewCounterWithTags(opts)
		if err != nil {
			return nil, err
		}
		gc.AddCounter(c)

	}

	m.counters = append(m.counters, gc)
	return gc, nil
}

func (m *Metrics) MustNewCounterWithTags(opts Options) *GroupTagCounter {
	c, err := m.NewCounterWithTags(opts)
	if err != nil {
		panic("Metric Error: " + err.Error())
	}
	return c
}

func (m *Metrics) NewGaugeWithTags(opts Options) (*GroupTagGauge, error) {
	gc := &GroupTagGauge{
		options: opts,
	}

	for _, e := range m.engines {
		c, err := e.NewGaugeWithTags(opts)
		if err != nil {
			return nil, err
		}
		gc.AddGauge(c)
	}

	m.gauges = append(m.gauges, gc)
	return gc, nil
}

func (m *Metrics) MustNewGaugeWithTags(opts Options) *GroupTagGauge {
	gc, err := m.NewGaugeWithTags(opts)
	if err != nil {
		panic("Metric Error: " + err.Error())
	}
	return gc
}

func (m *Metrics) NewHistogramWithTags(opts HistogramOptions) (*GroupTagHistogram, error) {
	gc := &GroupTagHistogram{
		options: opts,
	}
	for _, e := range m.engines {
		c, err := e.NewHistogramWithTags(opts)
		if err != nil {
			return nil, err
		}
		gc.AddHistogram(c)
	}
	m.observers = append(m.observers, gc)
	return gc, nil
}

func (m *Metrics) MustNewHistogramWithTags(opts HistogramOptions) *GroupTagHistogram {
	gc, err := m.NewHistogramWithTags(opts)
	if err != nil {
		panic("Metric Error: " + err.Error())
	}
	return gc
}

func (m *Metrics) Handler() http.Handler {
	return m.mh
}

type mhandler struct {
	Handler http.Handler
}

type MetricsHandler struct {
	handler *mhandler
}

func (mh MetricsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if mh.handler != nil {
		mh.handler.Handler.ServeHTTP(w, req)
	}
}
