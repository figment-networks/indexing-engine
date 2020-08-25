package metrics

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	// DefaultMetrics is the generic metrics aggregator
	// needed for global variables so we can attach new metrics
	// before the engine (like prometheus) is linked
	DefaultMetrics *Metrics
	// Debug flag
	Debug bool
)

func init() {
	DefaultMetrics = NewMetrics()
	_, ok := os.LookupEnv("INDEXING_DEBUG_ENABLE")
	if ok {
		Debug = true
	}
}

// Hotload links given engine to previously added metrics
func Hotload(name string) error {
	return DefaultMetrics.Hotload(name)
}

// AddEngine adds new engine to DefaultMetrics
func AddEngine(eng MetricsEngine) error {
	return DefaultMetrics.AddEngine(eng)
}

// Handler returns an http handler from default
func Handler() http.Handler {
	return DefaultMetrics.Handler()
}

// MetricsEngine is an interface for metric engines
type MetricsEngine interface {
	Name() string

	NewCounterWithTags(opts Options) (TagCounter, error)
	NewGaugeWithTags(opts Options) (TagGauge, error)
	NewHistogramWithTags(opts HistogramOptions) (TagObserver, error)
	Handler() http.Handler
}

// Options for metrics
type Options struct {
	Namespace string
	Subsystem string
	Name      string
	Desc      string
	Tags      []string
}

// Metrics a structure that group all the defined metrics
type Metrics struct {
	handler *mhandler
	mh      MetricsHandler

	engines map[string]MetricsEngine

	counters  []*GroupTagCounter
	gauges    []*GroupTagGauge
	observers []*GroupTagHistogram
}

// NewMetrics a metrics constructor
func NewMetrics() *Metrics {
	engines := make(map[string]MetricsEngine)
	handle := &mhandler{}
	return &Metrics{handler: handle, mh: MetricsHandler{handler: handle}, engines: engines}
}

// AddEngine adds new engine (like prometheus) to run
func (m *Metrics) AddEngine(eng MetricsEngine) error {
	_, ok := m.engines[eng.Name()]
	if ok {
		return fmt.Errorf("Engine %s is already registred", eng.Name())
	}
	m.engines[eng.Name()] = eng
	return nil
}

// Hotload loads all the previously defined metrics to given engine
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

// NewCounterWithTags create a group of counters from defined engines
func (m *Metrics) NewCounterWithTags(opts Options) (*GroupTagCounter, error) {
	gc := &GroupTagCounter{
		options:   opts,
		registred: make(map[uint64]*GroupCounter),
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

// MustNewCounterWithTags constructor with panic on error embedded
func (m *Metrics) MustNewCounterWithTags(opts Options) *GroupTagCounter {
	c, err := m.NewCounterWithTags(opts)
	if err != nil {
		panic("Metric Error: " + err.Error())
	}
	return c
}

// NewGaugeWithTags create a group of gauges from defined engines
func (m *Metrics) NewGaugeWithTags(opts Options) (*GroupTagGauge, error) {
	gc := &GroupTagGauge{
		options:   opts,
		registred: make(map[uint64]*GroupGauge),
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

// MustNewGaugeWithTags constructor with panic on error embedded
func (m *Metrics) MustNewGaugeWithTags(opts Options) *GroupTagGauge {
	gc, err := m.NewGaugeWithTags(opts)
	if err != nil {
		panic("Metric Error: " + err.Error())
	}
	return gc
}

// NewHistogramWithTags create a group of histograms from defined engines
func (m *Metrics) NewHistogramWithTags(opts HistogramOptions) (*GroupTagHistogram, error) {
	gc := &GroupTagHistogram{
		options:   opts,
		registred: make(map[uint64]*GroupObserver),
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

// MustNewHistogramWithTags constructor with panic on error embedded
func (m *Metrics) MustNewHistogramWithTags(opts HistogramOptions) *GroupTagHistogram {
	gc, err := m.NewHistogramWithTags(opts)
	if err != nil {
		panic("Metric Error: " + err.Error())
	}
	return gc
}

// Handler Returns metrics custom handler
func (m *Metrics) Handler() http.Handler {
	return m.mh
}

type mhandler struct {
	Handler http.Handler
}

// MetricsHandler a way to merge http handlers
type MetricsHandler struct {
	handler *mhandler
}

// ServeHTTP fulfills http.Handler interface
func (mh MetricsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if mh.handler != nil {
		mh.handler.Handler.ServeHTTP(w, req)
	}
}
