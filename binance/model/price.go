package model

import "github.com/ppincak/rysen/pkg/converters"

type Price struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (price *Price) PriceFloat64() float64 {
	p, _ := converters.ToFloat64(price.Price)
	return p
}
