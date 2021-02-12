package pipeline

import "github.com/figment-networks/indexing-engine/metrics"

var (
	taskDurationMetric = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexers",
		Subsystem: "pipeline_task",
		Name:      "task_duration",
		Desc:      "The total time spent processing an indexing task",
		Tags:      []string{"task"},
	})
)
