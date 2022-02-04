package ws

import (
	"rysen/monitor"
	"rysen/pkg/async"
)

type WsMetrics struct {
	Clients     *async.Counter
	Connects    *async.Counter
	Disconnects *async.Counter

	Reads        *async.Counter
	Writes       *async.Counter
	ReadsFailed  *async.Counter
	WritesFailed *async.Counter
}

func NewWsMetrics() *WsMetrics {
	return &WsMetrics{
		Clients:      async.NewCounter(),
		Connects:     async.NewCounter(),
		Disconnects:  async.NewCounter(),
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
			"clients":      metrics.Clients.Value(),
			"connects":     metrics.Connects.Value(),
			"disconnects":  metrics.Disconnects.Value(),
			"reads":        metrics.Reads.Value(),
			"writes":       metrics.Writes.Value(),
			"readsFailed":  metrics.ReadsFailed.Value(),
			"writesFailed": metrics.WritesFailed.Value(),
		},
	}
}
