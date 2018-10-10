package schema

import (
	"sync"

	"github.com/ppincak/rysen/crypto"
	"github.com/ppincak/rysen/pkg/errors"
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
	schemas           map[string]*Model
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
func (service *Service) Initialize(schemas []*Model) (err error) {
	for _, schema := range schemas {
		_, err = service.Create(schema)
		if err != nil {
			return
		}
	}
	return
}

// Create the schema
func (service *Service) Create(schema *Model) (instance *ExchangeSchema, err error) {
	service.lock.Lock()

	defer func() {
		service.lock.Unlock()

		// Teardown in case of an error
		if err != nil && instance != nil {
			instance.Destroy()
		}
	}()

	exchange, ok := service.exchanges[schema.Exchange]
	if !ok {
		return nil, errors.NewError("Exchange [%s] not found", schema.Exchange)
	}
	instance = NewExchangeSchema(schema)

	// Create Scrapers
	for i, model := range schema.Scrapers {
		if scraper, err := service.scraperService.Create(model, exchange.Caller(), exchange); err != nil {
			log.Error("Failed to create scraper from model [%#v]", model)
			return nil, err
		} else {
			instance.scrapers[i] = scraper

			log.Debugf("Created scraper from model [%#v]", model)
		}
	}

	// Create Aggregators
	for i, model := range schema.Aggregators {
		if aggregator, err := service.aggregatorService.Create(model, exchange.Aggregations()); err != nil {
			log.Error("Failed to create aggregator from model [%#v]", model)
			return nil, err
		} else {
			instance.aggregators[i] = aggregator

			log.Debugf("Created aggregator [%#v]", model)
		}
	}

	service.schemaInstances[schema.Name] = instance
	return instance, nil
}

// Delete schema
func (service *Service) DeleteSchema(schemaName string) error {
	schema, ok := service.schemaInstances[schemaName]
	if !ok {
		return errors.NewError("Schema not found")
	}
	schema.Destroy()

	return nil
}

// Return schema by name
func (service *Service) GetSchema(schemaName string) (*Model, error) {
	schema, ok := service.schemaInstances[schemaName]
	if !ok {
		return nil, errors.NewError("Schema not found")
	}
	return schema.model, nil
}

// Return list of all registered schemas
func (service *Service) ListSchemas() []*Model {
	defer service.lock.RUnlock()
	service.lock.RLock()

	result := make([]*Model, len(service.schemaInstances))
	i := 0
	for _, instance := range service.schemaInstances {
		result[i] = instance.model
		i++
	}
	return result
}
