package aggregate

import "github.com/ppincak/rysen/bus"

// Generic aggregation condition
func AggretateTillSize(size int) bus.AggregationCondition {
	return func(from int64, to int64, s int) bool {
		return size == s
	}
}

// Generic aggregation condition
func AggretateTillTime(elapsedTime int64) bus.AggregationCondition {
	return func(from int64, to int64, s int) bool {
		return from-to >= elapsedTime
	}
}
