package bus

type TransformFunc func(interface{}) (interface{}, error)
type AggregationFunc func(entries interface{}, lastEntry interface{}, from int64) (interface{}, error)
