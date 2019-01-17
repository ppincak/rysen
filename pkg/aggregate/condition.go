package aggregate

import (
	"github.com/ppincak/rysen/pkg/collections"
)

type ConditionType uint8

// Condition types
const (
	TillSize ConditionType = iota
	TillTime
)

var ConversionTable = map[string]ConditionType{
	"tillSize": TillSize,
	"tillTime": TillTime,
}

// Generic aggregation condition
func AggretateTillSize(size int) AggregationCondition {
	return func(from int64, to int64, entries *collections.SliceList) bool {
		return size == entries.Size()
	}
}

// Generic aggregation condition
func AggretateTillTime(elapsedTime int64) AggregationCondition {
	return func(from int64, to int64, entries *collections.SliceList) bool {
		return from-to >= elapsedTime
	}
}
