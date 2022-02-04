package server

import (
	"rysen/crypto"
	"rysen/monitor"
	"rysen/pkg/bus"
	"rysen/pkg/ws"
	"rysen/services/aggregator"
	"rysen/services/feed"
	"rysen/services/publisher"
	"rysen/services/schema"
	"rysen/services/scraper"
	"rysen/services/security"
)

// Application component container
type App struct {
	Bus       *bus.Bus
	Exchanges crypto.Exchanges

	// Feed
	FeedService     *feed.Service
	FeedPersistence *feed.Persistence

	// Publisher
	PublisherService *publisher.Service

	// Schema
	SchemaService     *schema.Service
	SchemaPersistence *schema.Persistence

	// Security
	SecurityService     *security.Service
	SecurityPersistence *security.Persistence

	AggregatorService *aggregator.Service
	ScraperService    *scraper.Service

	Monitor   *monitor.Monitor
	WsHandler *ws.Handler
}
