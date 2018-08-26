package services

import (
	"github.com/ppincak/rysen/pkg/aggregate"
)

// Metadata structure describing aggregator
type AggregatorMetadata struct {
	ReadTopic       string               `json:"readTopic"`
	WriteTopic      string               `json:"writeTopic"`
	AggregationFunc string               `json:"aggregationFunc"`
	ProcessFunc     string               `json:"processFunc"`
	Condition       AggregationCondition `json:"condition"`
}

// Note: maybe create some dsl in the future
type AggregationCondition struct {
	Func  aggregate.ConditionType
	Value int64
}
