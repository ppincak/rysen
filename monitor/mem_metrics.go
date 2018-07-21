package monitor

import "runtime"

type MemMetrics struct{}

func (metrics *MemMetrics) Map() map[string]interface{} {
	var rt runtime.MemStats
	runtime.ReadMemStats(&rt)
	return map[string]interface{}{
		"Alloc":        int64(rt.Alloc),
		"Frees":        int64(rt.Frees),
		"Lookups":      int64(rt.Lookups),
		"Mallocs":      int64(rt.Mallocs),
		"NumGC":        int64(rt.NumGC),
		"NumGoroutine": int64(runtime.NumGoroutine()),
		"HeapAlloc":    int64(rt.HeapAlloc),
		"HeapSys":      int64(rt.HeapSys),
		"HeapIdle":     int64(rt.HeapIdle),
		"HeapInUse":    int64(rt.HeapInuse),
		"HeapReleased": int64(rt.HeapReleased),
		"HeapObjects":  int64(rt.HeapObjects),
		"PauseTotalNs": int64(rt.PauseTotalNs),
		"Sys":          int64(rt.Sys),
		"TotalAlloc":   int64(rt.TotalAlloc),
	}
}
