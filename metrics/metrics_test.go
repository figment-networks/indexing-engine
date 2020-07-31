package metrics_test

import (
	"expvar"
	"log"
	"net/http"
	"strings"

	"github.com/figment-networks/indexing-engine/metrics"
)

// Dummy test engine implementation
type MetricsEngineTest struct {
}

func (m *MetricsEngineTest) Name() string {
	return "MetricsEngineTest"
}

func (m *MetricsEngineTest) NewCounterWithTags(opts metrics.Options) (metrics.TagCounter, error) {
	return NewMetricElementTestCouter(opts), nil
}

func (m *MetricsEngineTest) NewGaugeWithTags(opts metrics.Options) (metrics.TagGauge, error) {
	return nil, nil // NOOP
}

func (m *MetricsEngineTest) NewHistogramWithTags(opts metrics.HistogramOptions) (metrics.TagObserver, error) {
	return nil, nil // NOOP
}

func (m *MetricsEngineTest) Handler() http.Handler {
	return nil
}

type MetricElementTestCouter struct {
	metric *MyCounter
}

func NewMetricElementTestCouter(opts metrics.Options) *MetricElementTestCouter {
	return &MetricElementTestCouter{
		metric: &MyCounter{
			expvar.NewInt(strings.Join([]string{opts.Namespace, opts.Subsystem, opts.Name}, "_")),
		},
	}
}

func (me *MetricElementTestCouter) WithTags(map[string]string) (metrics.Counter, error) {
	// (lukanus): discard tags for the sake of simple example
	return me.metric, nil
}

func (me *MetricElementTestCouter) WithLabels(...string) metrics.Counter {
	// (lukanus): discard tags for the sake of simple example
	return me.metric
}

type MyCounter struct {
	metric *expvar.Int
}

func (mc MyCounter) Inc() {
	mc.metric.Add(1)
}
func (mc MyCounter) Add(a float64) {
	mc.metric.Add(int64(a))
}

func ExampleMetrics_NewCounterWithTags() {
	// Add metrics engine of yout choice metrics global
	me := &MetricsEngineTest{}
	err := metrics.AddEngine(me)
	if err != nil {
		log.Fatal(err)
	}

	// If needed link predeclared values
	err = metrics.Hotload(me.Name())
	if err != nil {
		log.Fatal(err)
	}

	// Create counter in all registred engines
	counter, err := metrics.NewCounterWithTags(metrics.Options{Namespace: "my", Subsystem: "test", Name: "metric"})

	// Use counter
	counter.WithLabels().Add(1)
}
