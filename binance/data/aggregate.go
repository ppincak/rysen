package data

import (
	"reflect"

	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/binance/model"
	"github.com/ppincak/rysen/math"
	"github.com/ppincak/rysen/pkg/aggregate"
	"gonum.org/v1/gonum/floats"
)

const (
	AveragePricesFunc  aggregate.AggregationType = "averagePrices"
	SumTradesFunc      aggregate.AggregationType = "sumTrades"
	SumGroupTradesFunc aggregate.AggregationType = "sumGroupTrades"
)

var Aggregations = map[aggregate.AggregationType]aggregate.AggregationFunc{
	AveragePricesFunc:  AveragePrices,
	SumTradesFunc:      SumTrades,
	SumGroupTradesFunc: SumGroupTrades,
}

// Average prices
func AveragePrices(message interface{}, lastEntry interface{}, from int64) (interface{}, error) {
	if models, ok := message.([]interface{}); ok {
		prices := make([]float64, len(models))
		for i, m := range models {
			if asserted, ok := m.(*model.Price); ok {
				prices[i] = asserted.PriceFloat64()
			}
		}
		return math.NewAverage(prices), nil
	}
	return nil, api.NewError("Invalid assertion")
}

// Create map from array of trades
func TradesMap(message interface{}, lastEntry interface{}, from int64) (map[int64]*model.Trade, error) {
	if models, ok := message.([]interface{}); ok {
		trades := make(map[int64]*model.Trade)
		for _, m := range models {
			if asserted, ok := m.([]*model.Trade); ok {
				var lastTrade *model.Trade

				if lastEntry != nil {
					lastTrades := lastEntry.([]*model.Trade)
					lastTrade = lastTrades[len(lastTrades)-1]
				}

				for _, trade := range asserted {
					if lastTrade != nil && trade.Id < lastTrade.Id {
						continue
					}
					if lastTrade == nil && trade.Time >= from {
						continue
					}
					trades[trade.Id] = trade
				}
			}
		}
		return trades, nil
	}
	return nil, api.NewError("Invalid assertion of type [%s]", reflect.TypeOf(message).Name)
}

// Sum trades
func SumTrades(message interface{}, lastEntry interface{}, from int64) (interface{}, error) {
	trades, err := TradesMap(message, lastEntry, from)
	if err != nil {
		return nil, err
	}

	prices := make([]float64, 0)
	quantities := make([]float64, 0)

	for _, trade := range trades {
		prices = append(prices, trade.PriceFloat64())
		quantities = append(quantities, trade.QuantityFloat64())
	}
	return &math.Sum{
		Metric: floats.Sum(prices),
		Volume: floats.Sum(quantities),
	}, nil
}

// Sum trades by price
func SumGroupTrades(message interface{}, lastEntry interface{}, from int64) (interface{}, error) {
	trades, err := TradesMap(message, lastEntry, from)
	if err != nil {
		return nil, err
	}

	priceQuantity := make(map[float64][]float64)
	for _, trade := range trades {
		val, ok := priceQuantity[trade.PriceFloat64()]
		if !ok {
			val = make([]float64, 0)
		}
		priceQuantity[trade.PriceFloat64()] = append(val, trade.QuantityFloat64())
	}
	return math.SumGroupedByFloat64(priceQuantity), nil
}
