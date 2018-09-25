package server

import (
	"github.com/ppincak/rysen/crypto"
	"github.com/ppincak/rysen/monitor"
	"github.com/ppincak/rysen/pkg/bus"
	"github.com/ppincak/rysen/pkg/ws"
	"github.com/ppincak/rysen/services/aggregator"
	"github.com/ppincak/rysen/services/feed"
	"github.com/ppincak/rysen/services/schema"
	"github.com/ppincak/rysen/services/scraper"
)

// Application component container
type App struct {
	Exchanges         crypto.Exchanges
	Bus               *bus.Bus
	SchemaService     *schema.Service
	AggregatorService *aggregator.Service
	FeedService       *feed.Service
	ScraperService    *scraper.Service
	Monitor           *monitor.Monitor
	WsHandler         *ws.Handler
}
