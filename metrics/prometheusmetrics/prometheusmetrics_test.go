package prometheusmetrics

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/stretchr/testify/require"
)

func TestMetrics_NewCounterWithTags(t *testing.T) {

	t.Run("New counter with tags", func(t *testing.T) {

		m := New()
		err := metrics.DetaultMetrics.AddEngine(m)
		require.NoError(t, err)

		mux := http.NewServeMux()
		mux.Handle("/metrics", m.Handler())
		srv := httptest.NewServer(mux)
		defer srv.Close()

		got, err := metrics.NewCounterWithTags(metrics.Options{
			Namespace: "a",
			Subsystem: "b",
			Name:      "c",
			Desc:      "d",
			Tags:      []string{"e", "f", "g"},
		})
		require.NoError(t, err)

		counter := got.WithLabels([]string{"e1", "f1", "g1"})
		counter.Inc()
		counter.Inc()
		counter.Inc()

		res, err := http.Get(srv.URL + "/metrics")
		data, err := ioutil.ReadAll(res.Body)
		res.Body.Close()

		require.Contains(t, string(data), `a_b_c{e="e1",f="f1",g="g1"} 3`)

	})
}
