package schema

import (
	"github.com/ppincak/rysen/pkg/aggregate"
	"github.com/ppincak/rysen/pkg/scrape"
	"github.com/ppincak/rysen/services/aggregator"
	"github.com/ppincak/rysen/services/feed"
	"github.com/ppincak/rysen/services/scraper"
)

// Metadata for a single exchange
type ExchangeSchemaMetadata struct {
	Name        string                 `json:"name"`
	Exchange    string                 `json:"exchange"`
	Scrapers    []*scraper.Metadata    `json:"scrapers"`
	Aggregators []*aggregator.Metadata `json:"aggregators"`
	Feeds       []*feed.Metadata       `json:"feeds"`
}

// Represents instance of the exchange schema,
// basically it is just a container for all the components
type ExchangeSchema struct {
	metadata *ExchangeSchemaMetadata

	scrapers    []*scrape.Scraper
	aggregators []*aggregate.Aggregator
	feeds       []*feed.Feed
}

// Create new schema instance
func NewExchangeSchema(metadata *ExchangeSchemaMetadata) *ExchangeSchema {
	return &ExchangeSchema{
		metadata:    metadata,
		scrapers:    make([]*scrape.Scraper, len(metadata.Scrapers)),
		aggregators: make([]*aggregate.Aggregator, len(metadata.Aggregators)),
		feeds:       make([]*feed.Feed, len(metadata.Feeds)),
	}
}

// Destroy the schema
func (schema *ExchangeSchema) Destroy() {
	for _, scraper := range schema.scrapers {
		scraper.Stop()
	}
	for _, aggregator := range schema.aggregators {
		aggregator.Stop()
	}
}
