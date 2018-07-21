package monitor

import "runtime"

type GoMetrics struct{}

func (metrics *GoMetrics) Map() map[string]interface{} {
	return map[string]interface{}{
		"GOARCH":     runtime.GOARCH,
		"GOOS":       runtime.GOOS,
		"GOMAXPROCS": runtime.GOMAXPROCS(-1),
		"version":    runtime.Version(),
	}
}
