package core

import (
	"github.com/ppincak/rysen/monitor"
	"github.com/ppincak/rysen/pkg/async"
)

type WsMetrics struct {
	Reads        *async.Counter
	Writes       *async.Counter
	ReadsFailed  *async.Counter
	WritesFailed *async.Counter
}

func NewWsMetrics() *WsMetrics {
	return &WsMetrics{
		Reads:        async.NewCounter(),
		Writes:       async.NewCounter(),
		ReadsFailed:  async.NewCounter(),
		WritesFailed: async.NewCounter(),
	}
}

func (metrics *WsMetrics) ToStatistic(name string) *monitor.Statistic {
	return &monitor.Statistic{
		Name: name,
		Values: map[string]interface{}{
			"reads":        metrics.Reads.Value(),
			"writes":       metrics.Writes.Value(),
			"readsFailed":  metrics.ReadsFailed.Value(),
			"writesFailed": metrics.WritesFailed.Value(),
		},
	}
}
