package data

import (
	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/binance/model"
	"github.com/ppincak/rysen/math"
)

func AveragePrices(message interface{}) (interface{}, error) {
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

func SumTrades(message interface{}) (interface{}, error) {
	if models, ok := message.([]interface{}); ok {
		trades := make(map[int64]float64)
		for _, m := range models {
			if asserted, ok := m.([]*model.Trade); ok {
				for _, trade := range asserted {
					trades[trade.Id] = trade.PriceFloat64()
				}
			}
		}
	}
	return nil, api.NewError("Invalid assertion")
}

func SumOrders(message interface{}) (interface{}, error) {
	return nil, nil
}
