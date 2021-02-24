package database

import "github.com/figment-networks/indexing-engine/metrics"

var (
	PingMetric = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "health",
		Subsystem: "database",
		Name:      "ping",
		Desc:      "Duration how long it takes to ping",
		Tags:      []string{"database_type"},
	})

	SizeMetric = metrics.MustNewGaugeWithTags(metrics.Options{
		Namespace: "health",
		Subsystem: "database",
		Name:      "size",
		Desc:      "Current size of database",
		Tags:      []string{"database_type"},
	})
)
