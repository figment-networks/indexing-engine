package pipeline

import (
	"time"
)

// StateRecorder is responsible for recording statistics during pipeline execution
// TODO: Add stats for every stage and every task
func NewStatsRecorder() *StatsRecorder {
	return &StatsRecorder{
		Stat: Stat{
			StartTime: time.Now(),
		},
	}
}

type StatsRecorder struct {
	Stat
}

func NewStat() *Stat {
	return &Stat{
		StartTime: time.Now(),
	}
}

type Stat struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Success   bool
}

// SetComplete completes stat duration
func (s *Stat) SetCompleted(success bool) {
	s.EndTime = time.Now()
	s.Duration = time.Since(s.StartTime)
	s.Success = success
}
