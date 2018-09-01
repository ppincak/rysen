package schema

import (
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
	schemas           map[string]*ExchangeSchemaMetadata
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
		schemas:           make(map[string]*ExchangeSchemaMetadata),
	}
}

// Register an crypto exchange
func (service *Service) Register(name string, exchange crypto.Exchange) {
	service.exchanges[name] = exchange
}

// Create the schema
func (service *Service) Create(schema *ExchangeSchemaMetadata) error {
	exchange, ok := service.exchanges[schema.Exchange]
	if !ok {
		return api.NewError("Exchange [%s] not found", schema.Exchange)
	}

	for _, metadata := range schema.Scrapers {
		if _, err := service.scraperService.Create(metadata, exchange.Caller()); err != nil {
			log.Error("Failed to create scraper from metadata [%#v]", metadata)
			log.Error(err)
			continue
		}

		log.Debugf("Created scraper from metadata [%#v]", metadata)
	}

	for _, metadata := range schema.Aggregators {
		if _, err := service.aggregatorService.Create(metadata, exchange.Aggregations()); err != nil {
			log.Error("Failed to create aggregator from metadata [%#v]", metadata)
			log.Error(err)
			continue
		}

		log.Debugf("Created aggregator [%#v]", metadata)
	}

	for _, metadata := range schema.Feeds {
		log.Debugf("Creating feed [%#v]", metadata)

		service.feedService.Create(metadata)
	}

	return nil
}
