package schema

import (
	"github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/pkg/scrape"
	"github.com/ppincak/rysen/services"
)

// Metadata for a single exchange
type ExchangeSchemaMetadata struct {
	Name        string                         `json:"name"`
	Exchange    string                         `json:"exchange"`
	Feeds       []*services.FeedMetadata       `json:"feeds"`
	Scrapers    []*services.ScraperMetadata    `json:"scrapers"`
	Aggregators []*services.AggregatorMetadata `json:"aggregators"`
}

// Todo: rethink
type ExchangeSchema struct {
	metadata *ExchangeSchemaMetadata

	scrapers    []*scrape.Scraper
	aggregators []*bus.Aggregator
	feeds       []*services.Feed
}
