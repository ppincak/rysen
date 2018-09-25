package schema

import (
	"sync"

	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/crypto"
	"github.com/ppincak/rysen/services/aggregator"
	"github.com/ppincak/rysen/services/feed"
	"github.com/ppincak/rysen/services/scraper"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	aggregatorService *aggregator.Service
	feedService       *feed.Service
	scraperService    *scraper.Service
	exchanges         crypto.Exchanges
	schemaBackup      *SchemaBackup
	schemas           map[string]*ExchangeSchemaMetadata
	schemaInstances   map[string]*ExchangeSchema
	lock              *sync.RWMutex
}

// Create schema service
func NewService(
	aggregatorService *aggregator.Service,
	feedService *feed.Service,
	scraperService *scraper.Service,
	exchanges crypto.Exchanges) *Service {

	return &Service{
		aggregatorService: aggregatorService,
		feedService:       feedService,
		scraperService:    scraperService,
		exchanges:         exchanges,
		schemaInstances:   make(map[string]*ExchangeSchema),
		lock:              new(sync.RWMutex),
	}
}

// Intialize the whole schema from backup
func (service *Service) InitializeFromBackup() (err error) {
	schemas, err := service.schemaBackup.GetSchemas()
	if err != nil {
		return
	}
	for _, schema := range schemas {
		service.Create(schema)
	}
	return
}

// Create the schema
func (service *Service) Create(schema *ExchangeSchemaMetadata) (*ExchangeSchema, error) {
	defer service.lock.Unlock()
	service.lock.Lock()

	exchange, ok := service.exchanges[schema.Exchange]
	if !ok {
		return nil, api.NewError("Exchange [%s] not found", schema.Exchange)
	}

	instance := NewExchangeSchema(schema)

	// Create Scrapers
	for i, metadata := range schema.Scrapers {
		if scraper, err := service.scraperService.Create(metadata, exchange.Caller(), exchange); err != nil {
			log.Error("Failed to create scraper from metadata [%#v]", metadata)
			log.Error(err)
		} else {
			instance.scrapers[i] = scraper

			log.Debugf("Created scraper from metadata [%#v]", metadata)
		}
	}

	// Create Aggregators
	for i, metadata := range schema.Aggregators {
		if aggregator, err := service.aggregatorService.Create(metadata, exchange.Aggregations()); err != nil {
			log.Error("Failed to create aggregator from metadata [%#v]", metadata)
			log.Error(err)
		} else {
			instance.aggregators[i] = aggregator

			log.Debugf("Created aggregator [%#v]", metadata)
		}
	}

	// Create Feeds
	for i, metadata := range schema.Feeds {
		log.Debugf("Creating feed [%#v]", metadata)

		if feed, err := service.feedService.Create(metadata); err != nil {
			log.Error("Failed to create feed from metadata [%#v]", metadata)
			log.Error(err)
		} else {
			instance.feeds[i] = feed

			log.Debugf("Created feed from metadata [%#v]", metadata)
		}

	}

	service.schemaInstances[schema.Name] = instance
	return instance, nil
}

// Delete schema
func (service *Service) Delete(schemaName string) error {
	schema, ok := service.schemaInstances[schemaName]
	if !ok {
		return api.NewError("Schema not found")
	}
	schema.Destroy()

	return nil
}

// Return list of all registered schemas
func (service *Service) ListSchemas() []*ExchangeSchemaMetadata {
	defer service.lock.RUnlock()
	service.lock.RLock()

	result := make([]*ExchangeSchemaMetadata, len(service.schemaInstances))
	i := 0
	for _, instance := range service.schemaInstances {
		result[i] = instance.metadata
		i++
	}
	return result
}
