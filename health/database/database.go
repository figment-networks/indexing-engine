package database

import "github.com/figment-networks/indexing-engine/metrics"

var (
	PingMetric = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexerhealth",
		Subsystem: "database",
		Name:      "ping",
		Desc:      "Duration how long it takes to ping",
		Tags:      []string{"type"},
	})

	SizeMetric = metrics.MustNewGaugeWithTags(metrics.Options{
		Namespace: "indexerhealth",
		Subsystem: "database",
		Name:      "size",
		Desc:      "Current size of database",
		Tags:      []string{"type"},
	})
)
