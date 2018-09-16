package scraper

import (
	"sync"

	"github.com/ppincak/rysen/api"

	"github.com/ppincak/rysen/pkg/bus"
	"github.com/ppincak/rysen/pkg/scrape"
)

// Scraper Service
type Service struct {
	bus      *bus.Bus
	lock     *sync.RWMutex
	scrapers map[string]*scrape.Scraper
}

// Create new scraper service
func NewService(bus *bus.Bus) *Service {
	return &Service{
		bus:      bus,
		lock:     new(sync.RWMutex),
		scrapers: make(map[string]*scrape.Scraper),
	}
}

// Create scraper from metadata
func (service *Service) Create(metadata *Metadata, caller scrape.Caller) (*scrape.Scraper, error) {
	defer service.lock.Unlock()
	service.lock.Lock()

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

	service.scrapers[metadata.Topic] = scraper
	return scraper, nil
}
