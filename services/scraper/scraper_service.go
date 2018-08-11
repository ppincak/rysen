package scraper

import (
	"sync"

	"github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/pkg/scrape"
)

type ScraperService struct {
	bus      *bus.Bus
	caller   scrape.Caller
	scrapers map[string]*scrape.Scraper
	lock     sync.RWMutex
}

func NewScraperService(bus *bus.Bus, caller scrape.Caller) *ScraperService {
	return &ScraperService{}
}

func (service *ScraperService) Create(topic string, symbols []string, callerFuncType scrape.CallerFuncType, interval int64) *scrape.Scraper {
	var callerFunc scrape.CallerFunc

	switch callerFuncType {
	case scrape.OrdersEvent:
		callerFunc = service.caller.ScrapeOrders
	case scrape.PriceEvent:
		callerFunc = service.caller.ScrapePrices
	case scrape.TradesEvent:
		callerFunc = service.caller.ScrapeTrades
	}

	scraper := scrape.NewScraper(topic, symbols, callerFunc, interval)
	service.scrapers[topic] = scraper

	return scraper
}
