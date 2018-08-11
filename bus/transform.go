package bus

type TransformFunc func(interface{}) (interface{}, error)
type AggregationFunc func(interface{}) (interface{}, error)
