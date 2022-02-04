package scraper

import (
	"rysen/crypto"

	"rysen/pkg/bus"
	"rysen/pkg/collections"
	"rysen/pkg/errors"
	"rysen/pkg/scrape"
)

// Scraper Service
type Service struct {
	bus *bus.Bus
}

// Create new scraper service
func NewService(bus *bus.Bus) *Service {
	return &Service{
		bus: bus,
	}
}

// Create scraper from model
func (service *Service) Create(model *Model, caller scrape.Caller, exchange crypto.Exchange) (*scrape.Scraper, error) {
	if collections.ArrayOfStringContains(exchange.Symbols().Symbols, model.Symbols) == false {
		return nil, errors.NewError("Symbols cannot be used for scraping [%#v]", model.Symbols)
	}

	callerFuncType, ok := scrape.StringToFunc[model.ScrapeFunc]
	if !ok {
		return nil, errors.NewError("Invalid scrapeFunction [%s]", model.ScrapeFunc)
	}

	var callerFunc scrape.CallerFunc

	switch callerFuncType {
	case scrape.OrdersEvent:
		callerFunc = caller.ScrapeOrders
	case scrape.PriceEvent:
		callerFunc = caller.ScrapePrice
	case scrape.TradesEvent:
		callerFunc = caller.ScrapeTrades
	}

	scraper := scrape.NewScraper(model.Topic, model.Symbols, callerFunc, model.Interval)
	go scraper.Start()

	return scraper, nil
}
