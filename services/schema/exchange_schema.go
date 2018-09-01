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
	Feeds       []*feed.Metadata       `json:"feeds"`
	Scrapers    []*scraper.Metadata    `json:"scrapers"`
	Aggregators []*aggregator.Metadata `json:"aggregators"`
}

// Todo: rethink
type ExchangeSchema struct {
	metadata *ExchangeSchemaMetadata

	scrapers    []*scrape.Scraper
	aggregators []*aggregate.Aggregator
	feeds       []*feed.Feed
}
