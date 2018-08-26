package aggregate

import (
	"github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/pkg/collections"
)

type ConditionType uint8

const (
	TillSize ConditionType = iota
	TillTime
)

// Generic aggregation condition
func AggretateTillSize(size int) bus.AggregationCondition {
	return func(from int64, to int64, entries *collections.SliceList) bool {
		return size == entries.Size()
	}
}

// Generic aggregation condition
func AggretateTillTime(elapsedTime int64) bus.AggregationCondition {
	return func(from int64, to int64, entries *collections.SliceList) bool {
		return from-to >= elapsedTime
	}
}
