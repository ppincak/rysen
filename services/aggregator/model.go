package aggregator

import (
	"github.com/ppincak/rysen/pkg/aggregate"
)

// Model structure describing aggregator
type Model struct {
	ReadTopic       string               `json:"readTopic"`
	WriteTopic      string               `json:"writeTopic"`
	AggregationFunc string               `json:"aggregationFunc"`
	Condition       AggregationCondition `json:"condition"`
}

// Note: maybe create some dsl in the future
type AggregationCondition struct {
	Func  aggregate.ConditionType `json:"func"`
	Value int64                   `json:"value"`
}
