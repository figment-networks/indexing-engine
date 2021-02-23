package pipeline

import "github.com/figment-networks/indexing-engine/metrics"

var (
	taskDurationMetric = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexer",
		Subsystem: "pipeline",
		Name:      "task_duration",
		Desc:      "The total time spent processing an indexing task",
		Tags:      []string{"task"},
	})

	stageDurationMetric = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexer",
		Subsystem: "pipeline",
		Name:      "stage_duration",
		Desc:      "The total time spent processing an indexing stage",
		Tags:      []string{"stage"},
	})

	heightDurationMetric = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexer",
		Subsystem: "pipeline",
		Name:      "height_duration",
		Desc:      "The total time spent indexing a height",
	})

	heightsTotalMetric = metrics.MustNewCounterWithTags(metrics.Options{
		Namespace: "indexer",
		Subsystem: "pipeline",
		Name:      "heights_total",
		Desc:      "The total number of successfully indexed heights",
	})

	errorsTotalMetric = metrics.MustNewCounterWithTags(metrics.Options{
		Namespace: "indexer",
		Subsystem: "pipeline",
		Name:      "errors_total",
		Desc:      "The total number of indexing errors",
	})
)
