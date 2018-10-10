package schema

import (
	"github.com/ppincak/rysen/pkg/aggregate"
	"github.com/ppincak/rysen/pkg/scrape"
	"github.com/ppincak/rysen/services/aggregator"
	"github.com/ppincak/rysen/services/scraper"
)

// Model for a single exchange
type Model struct {
	Name        string              `json:"name"`
	Exchange    string              `json:"exchange"`
	Scrapers    []*scraper.Model    `json:"scrapers"`
	Aggregators []*aggregator.Model `json:"aggregators"`
}

// Represents instance of the exchange schema,
// basically it is just a container for all the components
type ExchangeSchema struct {
	model *Model

	scrapers    []*scrape.Scraper
	aggregators []*aggregate.Aggregator
}

// Create new schema instance
func NewExchangeSchema(model *Model) *ExchangeSchema {
	return &ExchangeSchema{
		model:       model,
		scrapers:    make([]*scrape.Scraper, len(model.Scrapers)),
		aggregators: make([]*aggregate.Aggregator, len(model.Aggregators)),
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
