package prometheusmetrics

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics Prometheus metrics engine
type Metrics struct {
	reg       *prometheus.Registry
	gatherers prometheus.Gatherers
}

// New Prometheus Engine metrics constructor
func New() (m *Metrics) {
	reg := prometheus.NewRegistry()

	reg.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	reg.Register(prometheus.NewGoCollector())
	return &Metrics{reg: reg, gatherers: prometheus.Gatherers{reg}}
}

// Name returns unique engine name
func (m *Metrics) Name() string {
	return "prometheus"
}

// Handler http Hangler for prometheus
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(
		m.gatherers[0],
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	)
}

// NewCounterWithTags creates new prometheus counter, registering it
func (m *Metrics) NewCounterWithTags(opts metrics.Options) (metrics.TagCounter, error) {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Desc,
		},
		opts.Tags,
	)

	err := m.reg.Register(counter)
	return &TagCounter{c: counter}, err
}

// NewGaugeWithTags creates new prometheus gauge, registering it
func (m *Metrics) NewGaugeWithTags(opts metrics.Options) (metrics.TagGauge, error) {
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Desc,
		},
		opts.Tags,
	)

	err := m.reg.Register(gauge)
	return &TagGauge{g: gauge}, err
}

// NewHistogramWithTags creates new histogram, registering it
func (m *Metrics) NewHistogramWithTags(opts metrics.HistogramOptions) (metrics.TagObserver, error) {

	options := prometheus.HistogramOpts{
		Namespace: opts.Namespace,
		Subsystem: opts.Subsystem,
		Name:      opts.Name,
		Help:      opts.Desc,
	}

	if opts.Buckets.Type != "" || len(opts.Buckets.Buckets) > 0 {
		var bList []float64
		for _, b := range opts.Buckets.Buckets {
			f, err := strconv.ParseFloat(b, 64)
			if err != nil {
				return nil, fmt.Errorf("Error parsing histogram buckets floats : %w", err)
			}
			bList = append(bList, f)
		}
		options.Buckets = bList
	}
	histogram := prometheus.NewHistogramVec(options, opts.Tags)

	err := m.reg.Register(histogram)
	return &TagHistogram{h: histogram}, err
}

// TagCounter Prometheus counter wrapper
type TagCounter struct {
	c *prometheus.CounterVec
}

// WithTags creates metric with given tags
func (tc *TagCounter) WithTags(tags map[string]string) (metrics.Counter, error) {
	return tc.c.GetMetricWith(tags)
}

// WithLabels creates metric with given tag values in defined order
func (tc *TagCounter) WithLabels(lv ...string) metrics.Counter {
	return tc.c.WithLabelValues(lv...)
}

// TagGauge Prometheus gauge wrapper
type TagGauge struct {
	g *prometheus.GaugeVec
}

// WithTags creates metric with given tags
func (tg *TagGauge) WithTags(tags map[string]string) (metrics.Gauge, error) {
	return tg.g.GetMetricWith(tags)
}

// WithLabels creates metric with given tag values in defined order
func (tg *TagGauge) WithLabels(lv ...string) metrics.Gauge {
	return tg.g.WithLabelValues(lv...)
}

// TagHistogram Prometheus histogram wrapper
type TagHistogram struct {
	h *prometheus.HistogramVec
}

// WithTags creates metric with given tags
func (th *TagHistogram) WithTags(tags map[string]string) (metrics.Observer, error) {
	return th.h.GetMetricWith(tags)
}

// WithLabels creates metric with given tag values in defined order
func (th *TagHistogram) WithLabels(lv ...string) metrics.Observer {
	return th.h.WithLabelValues(lv...)
}
