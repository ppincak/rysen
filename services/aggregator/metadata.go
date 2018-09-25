package aggregator

import (
	"github.com/ppincak/rysen/pkg/aggregate"
)

// Metadata structure describing aggregator
type Metadata struct {
	ReadTopic  string `json:"readTopic"`
	WriteTopic string `json:"writeTopic"`
	// Note: probably may be removed
	AggregationFunc string               `json:"aggregationFunc"`
	ProcessFunc     string               `json:"processFunc"`
	Condition       AggregationCondition `json:"condition"`
}

// Note: maybe create some dsl in the future
type AggregationCondition struct {
	Func  aggregate.ConditionType `json:"func"`
	Value int64                   `json:"value"`
}
