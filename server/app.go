package server

import (
	"github.com/ppincak/rysen/binance"
)

type App interface {
	BinanceClient() *binance.BinanceClient
}
