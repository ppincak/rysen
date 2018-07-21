package monitor

type Monitor struct {
	metrics   map[string]Metrics
	reporters map[Reporter]struct{}
}

func NewMonitor() *Monitor {
	return &Monitor{
		metrics: map[string]Metrics{
			"memory": &MemMetrics{},
			"golang": &GoMetrics{},
		},
		reporters: make(map[Reporter]struct{}),
	}
}

func (monitor *Monitor) Statistics() []*Statistic {
	statistics := make([]*Statistic, len(monitor.metrics))
	index := 0
	for key, metrics := range monitor.metrics {
		statistics[index] = &Statistic{
			Name:   key,
			Values: metrics.Map(),
		}
		index++
	}

	for reporter, _ := range monitor.reporters {
		statistics = append(statistics, reporter.Statistics()...)
	}

	return statistics
}

func (monitor *Monitor) Register(reporter Reporter) {
	monitor.reporters[reporter] = struct{}{}
}
