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

// Creates empty ExchangeSchemaMetadata
func NewExchangeSchemaMetadata(name string) *ExchangeSchemaMetadata {
	return &ExchangeSchemaMetadata{
		Name:        name,
		Scrapers:    make([]*scraper.Metadata, 0),
		Aggregators: make([]*aggregator.Metadata, 0),
		Feeds:       make([]*feed.Metadata, 0),
	}
}

// Note: maybe rename to ExchangeSchema
// Represents single instance of the exchange schema
type ExchangeSchemaInstance struct {
	metadata *ExchangeSchemaMetadata

	scrapers    []*scrape.Scraper
	aggregators []*aggregate.Aggregator
	feeds       []*feed.Feed
}

// Create new schema instance
func NewExchangeSchemaInstance(metadata *ExchangeSchemaMetadata) *ExchangeSchemaInstance {
	return &ExchangeSchemaInstance{
		metadata:    metadata,
		scrapers:    make([]*scrape.Scraper, len(metadata.Scrapers)),
		aggregators: make([]*aggregate.Aggregator, len(metadata.Aggregators)),
		feeds:       make([]*feed.Feed, len(metadata.Feeds)),
	}
}
