package services

import (
	"sync"

	"github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/pkg/scrape"
)

type ScraperService struct {
	bus      *bus.Bus
	lock     *sync.RWMutex
	scrapers map[string]*scrape.Scraper
}

func NewScraperService(bus *bus.Bus) *ScraperService {
	return &ScraperService{
		bus:      bus,
		lock:     new(sync.RWMutex),
		scrapers: make(map[string]*scrape.Scraper),
	}
}

func (service *ScraperService) Create(metadata *ScraperMetadata, caller scrape.Caller) *scrape.Scraper {
	defer service.lock.Unlock()
	service.lock.Lock()

	var callerFunc scrape.CallerFunc

	switch metadata.ScrapeFunc {
	case scrape.OrdersEvent:
		callerFunc = caller.ScrapeOrders
	case scrape.PriceEvent:
		callerFunc = caller.ScrapePrices
	case scrape.TradesEvent:
		callerFunc = caller.ScrapeTrades
	}

	scraper := scrape.NewScraper(metadata.Topic, metadata.Symbols, callerFunc, metadata.Interval)
	go scraper.Start()

	service.scrapers[metadata.Topic] = scraper
	return scraper
}
