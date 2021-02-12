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
)
