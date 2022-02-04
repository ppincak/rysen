package feed

import (
	"rysen/monitor"
	"rysen/pkg/async"
)

type Metrics struct {
	MessagesSent *async.Counter
}

// Create Metrics
func NewMetrics() *Metrics {
	return &Metrics{
		MessagesSent: async.NewCounter(),
	}
}

// Transform Metrics to Statistics
func (metrics *Metrics) ToStatistic(name string) *monitor.Statistic {
	return &monitor.Statistic{
		Name: name,
		Values: map[string]interface{}{
			"messagesSent": metrics.MessagesSent.Value(),
		},
	}
}
