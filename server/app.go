package server

import (
	"github.com/ppincak/rysen/binance"
	"github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/monitor"
	"github.com/ppincak/rysen/pkg/ws"
	"github.com/ppincak/rysen/services"
)

type App struct {
	Binance           *binance.Exchange
	Bus               *bus.Bus
	AggregatorService *services.AggregatorService
	FeedService       *services.FeedService
	ScraperService    *services.ScraperService
	Monitor           *monitor.Monitor
	WsHandler         *ws.Handler
}
