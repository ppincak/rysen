package binance

import (
	"rysen/binance/data"
	"rysen/crypto"
	"rysen/monitor"
	"rysen/pkg/aggregate"
	"rysen/pkg/bus"
	"rysen/pkg/errors"
	"rysen/pkg/scrape"
)

type Exchange struct {
	client *Client
	caller *Caller
	config *Config
	bus    *bus.Bus

	symbols *crypto.Symbols
}

var _ crypto.Exchange = (*Exchange)(nil)
var _ monitor.Reporter = (*Exchange)(nil)

func NewExchange(config *Config, bus *bus.Bus) *Exchange {
	client := NewClient(config.url)

	return &Exchange{
		client: client,
		config: config,
		bus:    bus,
	}
}

// Initialize the exchange
func (exchange *Exchange) Initialize() error {
	info, err := exchange.client.ExchangeInfo()
	if err != nil {
		return errors.NewError("Failed to initialize Binance Exchange")
	}
	caller := NewCaller(exchange.client, exchange.bus)
	exchange.caller = caller
	exchange.symbols = NewSymbols(info)

	return nil
}

// Get monitoring statistics
func (exchange *Exchange) Statistics() []*monitor.Statistic {
	return exchange.client.Statistics()
}

// Get aggregations
func (exchange *Exchange) Name() string {
	return "binance"
}

// Get aggregations
func (exchange *Exchange) Aggregations() aggregate.AggregationsMap {
	return data.Aggregations
}

// Get scrapers
func (exchange *Exchange) Scrapers() []string {
	return []string{
		"scrapeOrders",
		"scrapePrice",
		"scrapeTrades",
	}
}

// Get caller
func (exchange *Exchange) Caller() scrape.Caller {
	return exchange.caller
}

// Get symbols
func (exchange *Exchange) Symbols() *crypto.Symbols {
	return exchange.symbols
}
