package server

import (
	"github.com/ppincak/rysen/binance"
	"github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/monitor"
	"github.com/ppincak/rysen/pkg/ws"
	"github.com/ppincak/rysen/services"
)

type Binance struct {
	Client *binance.Client
	Caller *binance.Caller
	Store  *binance.Store
}

type App struct {
	Binance     *Binance
	Bus         *bus.Bus
	FeedService *services.FeedService
	Monitor     *monitor.Monitor
	WsHandler   *ws.Handler
}
