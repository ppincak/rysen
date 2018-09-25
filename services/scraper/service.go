package scraper

import (
	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/crypto"

	"github.com/ppincak/rysen/pkg/bus"
	"github.com/ppincak/rysen/pkg/collections"
	"github.com/ppincak/rysen/pkg/scrape"
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

// Create scraper from metadata
func (service *Service) Create(metadata *Metadata, caller scrape.Caller, exchange crypto.Exchange) (*scrape.Scraper, error) {
	if collections.ArrayOfStringContains(exchange.Symbols().Symbols, metadata.Symbols) == false {
		return nil, api.NewError("Symbols cannot be used for scraping [%#v]", metadata.Symbols)
	}

	callerFuncType, ok := scrape.StringToFunc[metadata.ScrapeFunc]
	if !ok {
		return nil, api.NewError("Invalid scrapeFunction [%s]", metadata.ScrapeFunc)
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

	scraper := scrape.NewScraper(metadata.Topic, metadata.Symbols, callerFunc, metadata.Interval)
	go scraper.Start()

	return scraper, nil
}
