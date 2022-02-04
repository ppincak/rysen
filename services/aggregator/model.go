package aggregator

import (
	"rysen/pkg/collections"
)

// Model structure describing aggregator
type Model struct {
	ReadTopic       string               `json:"readTopic"`
	WriteTopic      string               `json:"writeTopic"`
	AggregationFunc string               `json:"aggregationFunction"`
	Condition       AggregationCondition `json:"condition"`
}

// Note: maybe create some dsl in the future
type AggregationCondition struct {
	Func  string `json:"function"`
	Value int64  `json:"value"`
}

var _ collections.Comparable = (*Model)(nil)

// Equals func
func (model *Model) Equals(value interface{}) bool {
	assertion, ok := value.(*Model)
	if !ok {
		return false
	}
	if assertion.ReadTopic != model.ReadTopic {
		return false
	}
	if assertion.WriteTopic != model.WriteTopic {
		return false
	}
	return true
}
