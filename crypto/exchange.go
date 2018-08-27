package crypto

import (
	"github.com/ppincak/rysen/pkg/aggregate"
	"github.com/ppincak/rysen/pkg/scrape"
)

type Exchange interface {
	// Get all aggregations available for the exchange
	Aggregations() aggregate.AggregationsMap
	// Get caller
	Caller() scrape.Caller
	// Get all symbols available for exchange
	Symbols() *Symbols
}

type ExchangeContainer struct {
	exchanges map[string]Exchange
}

func NewExchangeContainer() *ExchangeContainer {
	return &ExchangeContainer{
		exchanges: make(map[string]Exchange),
	}
}
