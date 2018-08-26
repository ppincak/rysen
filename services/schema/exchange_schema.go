package schema

import (
	"github.com/ppincak/rysen/pkg/aggregate"
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
	aggregators []*aggregate.Aggregator
	feeds       []*services.Feed
}
