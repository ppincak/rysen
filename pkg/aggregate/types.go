package aggregate

import (
	"rysen/pkg/bus"
	"rysen/pkg/collections"
)

// Result of the aggregation
type AggregationResult struct {
	Result interface{} `json:"result"`
	From   int64       `json:"from"`
	To     int64       `json:"to"`
}

// Processing function
type ProcessFunc func(event *bus.BusEvent) (interface{}, error)

type AggregationFunc func(entries interface{}, lastEntry interface{}, from int64) (interface{}, error)
type AggregationCondition func(from int64, to int64, entries *collections.SliceList) bool
type AggregationType string
type AggregationsMap map[AggregationType]AggregationFunc
