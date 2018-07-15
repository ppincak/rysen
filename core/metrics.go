package core

import (
	"github.com/ppincak/rysen/common"
)

const (
	SuccessfullCalls = iota
	FailedCalls
)

type Metrics struct {
	successfullCalls *common.Counter
	failedCalls      *common.Counter
}

func NewMetrics() *Metrics {
	metrics := &Metrics{
		successfullCalls: common.NewCounter(),
		failedCalls:      common.NewCounter(),
	}
	return metrics
}

func (metrics *Metrics) SuccessfullCalls() int64 {
	return metrics.successfullCalls.Value()
}

func (metrics *Metrics) FailedCalls() int64 {
	return metrics.failedCalls.Value()
}

func (metrics *Metrics) Inc(field int) {
	switch field {
	case SuccessfullCalls:
		metrics.successfullCalls.Inc()
	case FailedCalls:
		metrics.failedCalls.Inc()
	}
}

func (metrics *Metrics) ToStatistic(name string) *common.Statistic {
	return &common.Statistic{
		Name: name,
		Values: map[string]interface{}{
			"succesfullCalls": metrics.successfullCalls.Value(),
			"failedCalls":     metrics.failedCalls.Value(),
		},
	}
}
