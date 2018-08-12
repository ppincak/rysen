package model

import "github.com/ppincak/rysen/pkg/converters"

type Trade struct {
	Id           int64  `json:"id"`
	Price        string `json:"price"`
	Quantity     string `json:"qty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	IsBestMatch  bool   `json:"isBestMatch"`
}

func (trade *Trade) PriceFloat64() float64 {
	p, _ := converters.ToFloat64(trade.Price)
	return p
}

func (trade *Trade) QuantityFloat64() float64 {
	qty, _ := converters.ToFloat64(trade.Quantity)
	return qty
}
