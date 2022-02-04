package core

import (
	"rysen/monitor"
	"rysen/pkg/async"
)

const (
	SuccessfullCalls = iota
	FailedCalls
)

type ApiMetrics struct {
	successfullCalls *async.Counter
	failedCalls      *async.Counter
}

func NewApiMetrics() *ApiMetrics {
	metrics := &ApiMetrics{
		successfullCalls: async.NewCounter(),
		failedCalls:      async.NewCounter(),
	}
	return metrics
}

func (metrics *ApiMetrics) SuccessfullCalls() int64 {
	return metrics.successfullCalls.Value()
}

func (metrics *ApiMetrics) FailedCalls() int64 {
	return metrics.failedCalls.Value()
}

func (metrics *ApiMetrics) Inc(field int) {
	switch field {
	case SuccessfullCalls:
		metrics.successfullCalls.Inc()
	case FailedCalls:
		metrics.failedCalls.Inc()
	}
}

func (metrics *ApiMetrics) ToStatistic(name string) *monitor.Statistic {
	return &monitor.Statistic{
		Name: name,
		Values: map[string]interface{}{
			"succesfullCalls": metrics.successfullCalls.Value(),
			"failedCalls":     metrics.failedCalls.Value(),
		},
	}
}
