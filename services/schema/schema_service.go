package schema

import (
	"github.com/ppincak/rysen/pkg/scrape"
	"github.com/ppincak/rysen/services"
)

type SchemaService struct {
	aggregatorService *services.AggregatorService
	feedService       *services.FeedService
	scraperService    *services.ScraperService
}

// Create schema service
func NewSchemaService(
	aggregatorService *services.AggregatorService,
	feedService *services.FeedService,
	scraperService *services.ScraperService) *SchemaService {

	return &SchemaService{
		aggregatorService: aggregatorService,
		feedService:       feedService,
		scraperService:    scraperService,
	}
}

// Create the schema
func (service *SchemaService) Create(schema *ExchangeSchemaMetadata, caller scrape.Caller) {
	for _, metadata := range schema.Scrapers {
		service.scraperService.Create(metadata, caller)
	}
	for _, metadata := range schema.Aggregators {
		service.aggregatorService.Create(metadata)
	}
	for _, metadata := range schema.Feeds {
		service.feedService.Create(metadata)
	}
}
