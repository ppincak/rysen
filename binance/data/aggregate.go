package data

import (
	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/math"
)

func AggregatePrice(message interface{}) (interface{}, error) {
	if result, ok := message.([]interface{}); ok {
		prices, err := ExtractPrices(result)
		if err == nil {
			return math.NewAverage(prices), nil
		} else {
			return nil, err
		}
	}
	return nil, api.NewError("Invalid assertion")
}

func AggretateTill(size int) bus.AggregationCondition {
	return func(from int64, to int64, s int) bool {
		return size == s
	}
}
