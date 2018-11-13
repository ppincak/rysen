package crypto

import (
	"github.com/ppincak/rysen/pkg/aggregate"
	"github.com/ppincak/rysen/pkg/scrape"
)

// Represent crypto exchange eg: Binance, Idex
// each exchange should have its own API Client and Caller
// and this Exchange is a form of packaging for this
type Exchange interface {
	// Get name of the exchange
	Name() string
	// Get all aggregations available for the exchange
	Aggregations() aggregate.AggregationsMap
	// Get all scraping options available for the exchange
	Scrapers() []string
	// Get caller
	Caller() scrape.Caller
	// Get all symbols available for exchange
	Symbols() *Symbols
}

// Base model for exchange
type ExchangeModel struct {
	Name         string              `json:"name"`
	Assets       map[string][]string `json:"assets"`
	Symbols      []string            `json:"symbols"`
	Scrapers     []string            `json:"scrapers"`
	Aggregations []string            `json:"aggregations"`
}

// Collections of multiple exchanges where the key is the name of exchange
type Exchanges map[string]Exchange

// Create new Exchange collection
func NewExchanges() Exchanges {
	return make(Exchanges)
}

// Register new Exchange
func (exchanges Exchanges) Register(exchange Exchange) {
	exchanges[exchange.Name()] = exchange
}

// Get list of all exchanges
func (exchanges Exchanges) List() []*ExchangeModel {
	result := make([]*ExchangeModel, len(exchanges))
	i := 0
	for _, exchange := range exchanges {
		exchangeSymbols := exchange.Symbols()

		j := 0
		aggregations := make([]string, len(exchange.Aggregations()))
		for aggregation, _ := range exchange.Aggregations() {
			aggregations[j] = string(aggregation)
			j++
		}

		result[i] = &ExchangeModel{
			Name:         exchange.Name(),
			Assets:       exchangeSymbols.Assets,
			Symbols:      exchangeSymbols.Symbols,
			Scrapers:     exchange.Scrapers(),
			Aggregations: aggregations,
		}
		i++
	}
	return result
}
