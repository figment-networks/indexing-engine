package prometheusmetrics

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	reg       *prometheus.Registry
	gatherers prometheus.Gatherers
}

func New() (m *Metrics) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	return &Metrics{reg: reg, gatherers: prometheus.Gatherers{
		reg,
	}}
}

func (m *Metrics) Name() string {
	return "prometheus"
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(
		m.gatherers[0],
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	)
}

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
		//Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		options.Buckets = bList
	}
	histogram := prometheus.NewHistogramVec(options, opts.Tags)

	err := m.reg.Register(histogram)
	return &TagHistogram{h: histogram}, err
}

type TagCounter struct {
	c *prometheus.CounterVec
}

func (tc *TagCounter) WithTags(tags map[string]string) (metrics.Counter, error) {
	return tc.c.GetMetricWith(tags)
}

func (tc *TagCounter) WithLabels(lv ...string) metrics.Counter {
	return tc.c.WithLabelValues(lv...)
}

type TagGauge struct {
	g *prometheus.GaugeVec
}

func (tg *TagGauge) WithTags(tags map[string]string) (metrics.Gauge, error) {
	return tg.g.GetMetricWith(tags)
}

func (tg *TagGauge) WithLabels(lv ...string) metrics.Gauge {
	return tg.g.WithLabelValues(lv...)
}

type TagHistogram struct {
	h *prometheus.HistogramVec
}

func (th *TagHistogram) WithTags(tags map[string]string) (metrics.Observer, error) {
	return th.h.GetMetricWith(tags)
}

func (th *TagHistogram) WithLabels(lv ...string) metrics.Observer {
	return th.h.WithLabelValues(lv...)
}
