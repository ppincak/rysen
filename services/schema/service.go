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
	exchanges         map[string]crypto.Exchange
	schemaInstances   map[string]*ExchangeSchemaInstance
	lock              *sync.RWMutex
}

// Create schema service
func NewService(
	aggregatorService *aggregator.Service,
	feedService *feed.Service,
	scraperService *scraper.Service) *Service {

	return &Service{
		aggregatorService: aggregatorService,
		feedService:       feedService,
		scraperService:    scraperService,
		exchanges:         make(map[string]crypto.Exchange),
		schemaInstances:   make(map[string]*ExchangeSchemaInstance),
		lock:              new(sync.RWMutex),
	}
}

// Register an crypto exchange
func (service *Service) RegisterExchange(name string, exchange crypto.Exchange) {
	service.exchanges[name] = exchange
}

// Create the schema
func (service *Service) Create(schema *ExchangeSchemaMetadata) (*ExchangeSchemaInstance, error) {
	defer service.lock.Unlock()
	service.lock.Lock()

	exchange, ok := service.exchanges[schema.Exchange]
	if !ok {
		return nil, api.NewError("Exchange [%s] not found", schema.Exchange)
	}

	instance := NewExchangeSchemaInstance(schema)

	// Create Scrapers
	for i, metadata := range schema.Scrapers {
		if scraper, err := service.scraperService.Create(metadata, exchange.Caller()); err != nil {
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
