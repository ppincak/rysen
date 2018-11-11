package aggregator

import (
	"github.com/ppincak/rysen/pkg/aggregate"
	"github.com/ppincak/rysen/pkg/collections"
)

// Model structure describing aggregator
type Model struct {
	ReadTopic  string `json:"readTopic"`
	WriteTopic string `json:"writeTopic"`
	// TODO clean local storage and rename
	AggregationFunc string               `json:"aggregationFunc"`
	Condition       AggregationCondition `json:"condition"`
}

// Note: maybe create some dsl in the future
type AggregationCondition struct {
	Func  aggregate.ConditionType `json:"func"`
	Value int64                   `json:"value"`
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
