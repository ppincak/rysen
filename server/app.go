package server

import (
	"github.com/ppincak/rysen/binance"
	"github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/monitor"
	"github.com/ppincak/rysen/pkg/ws"
)

type Binance struct {
	Client  *binance.Client
	Caller  *binance.Caller
	Scraper *binance.Scraper
	Store   *binance.Store
}

type App struct {
	Binance   *Binance
	Bus       *bus.Bus
	Monitor   *monitor.Monitor
	WsHandler *ws.Handler
}
