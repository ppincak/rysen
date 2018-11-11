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
		_, err = service.CreateSchema(schema)
		if err != nil {
			return
		}
	}
	return
}

// CreateSchema the schema
func (service *Service) CreateSchema(schema *Model) (instance *ExchangeSchema, err error) {
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
	if _, ok := service.schemaInstances[schema.Name]; ok {
		return nil, errors.NewError("Schema with name [%s] already exists", schema.Name)
	}

	instance = NewExchangeSchema(schema)

	err = service.createSchema(schema, instance, exchange)
	if err != nil {
		return nil, err
	} else {
		return instance, err
	}
}

func (service *Service) createSchema(schema *Model, instance *ExchangeSchema, exchange crypto.Exchange) (err error) {
	// Create Scrapers
	for i, model := range schema.Scrapers {
		if scraper, err := service.scraperService.Create(model, exchange.Caller(), exchange); err != nil {
			log.Error("Failed to create scraper from model [%#v]", model)
			return err
		} else {
			instance.scrapers[i] = scraper

			log.Debugf("Created scraper from model [%#v]", model)
		}
	}

	// Create Aggregators
	for i, model := range schema.Aggregators {
		if aggregator, err := service.aggregatorService.Create(model, exchange.Aggregations()); err != nil {
			log.Error("Failed to create aggregator from model [%#v]", model)
			return err
		} else {
			instance.aggregators[i] = aggregator

			log.Debugf("Created aggregator [%#v]", model)
		}
	}

	service.schemaInstances[schema.Name] = instance
	return nil
}

// Update the schema
// Note: may be not the best solution performance wise, but seems like the most reasonable to implement
func (service *Service) UpdateSchema(schema *Model) (instance *ExchangeSchema, err error) {
	defer func() {
		service.lock.Unlock()

		// Teardown in case of an error
		if err != nil && instance != nil {
			instance.Destroy()
		}
	}()
	service.lock.Lock()

	exchange, ok := service.exchanges[schema.Exchange]
	if !ok {
		return nil, errors.NewError("Exchange [%s] not found", schema.Exchange)
	}

	if i, ok := service.schemaInstances[schema.Name]; !ok {
		return nil, errors.NewError("Schema with name [%s] not found", schema.Name)
	} else {
		instance = i
	}
	instance.Destroy()

	instance = NewExchangeSchema(schema)

	err = service.createSchema(schema, instance, exchange)
	if err != nil {
		log.Error("Failed to create schema [%#v]", schema)
		return nil, err
	}
	return instance, nil
}

// Delete schema
func (service *Service) DeleteSchema(schemaName string) error {
	defer service.lock.Unlock()
	service.lock.Lock()

	instance, ok := service.schemaInstances[schemaName]
	if !ok {
		return errors.NewError("Schema not found")
	}
	instance.Destroy()

	return nil
}

// Return schema by name
func (service *Service) GetSchema(schemaName string) (*Model, error) {
	defer service.lock.RUnlock()
	service.lock.RLock()

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
