package server

import (
	"github.com/ppincak/rysen/crypto"
	"github.com/ppincak/rysen/monitor"
	"github.com/ppincak/rysen/pkg/bus"
	"github.com/ppincak/rysen/pkg/ws"
	"github.com/ppincak/rysen/services/aggregator"
	"github.com/ppincak/rysen/services/feed"
	"github.com/ppincak/rysen/services/publisher"
	"github.com/ppincak/rysen/services/schema"
	"github.com/ppincak/rysen/services/scraper"
	"github.com/ppincak/rysen/services/security"
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
