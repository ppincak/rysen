package binance

import (
	"time"

	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/core"
)

type BinanceCaller struct {
	client  *BinanceClient
	bus     *BinanceBus
	limiter *core.CallLimiter
}

func NewCaller(client *BinanceClient, bus *BinanceBus, limiter *core.CallLimiter) *BinanceCaller {
	return &BinanceCaller{
		client:  client,
		bus:     bus,
		limiter: limiter,
	}
}

type CallerFuncType uint32

const (
	OrdersEvent CallerFuncType = iota
	TradesEvent
	PriceEvent
)

type CallerFunc func(topic string, symbols []string)

type CallerEvent struct {
	symbol    string
	eventType CallerFuncType
	data      interface{}
	timestamp int64
}

type Callable func(*BinanceClient, string) (api.ApiResponse, error)

func (caller *BinanceCaller) call(topic string, symbols []string, eventType CallerFuncType, callable Callable) {
	for _, symbol := range symbols {
		go func(client *BinanceClient, symbol string) {
			t := time.Now()

			resp, err := callable(client, symbol)

			go caller.limiter.IncrRequests(nil)

			if err == nil {
				caller.bus.Publish(topic, &CallerEvent{
					symbol:    symbol,
					eventType: eventType,
					data:      resp,
					timestamp: t.UnixNano(),
				})
			}
		}(caller.client, symbol)
	}
}

func (caller *BinanceCaller) ScrapeOrders(topic string, symbols []string) {
	caller.call(topic, symbols, OrdersEvent, func(client *BinanceClient, symbol string) (api.ApiResponse, error) {
		return client.OrderBook(symbol, 0)
	})
}

func (caller *BinanceCaller) ScrapeTrades(topic string, symbols []string) {
	caller.call(topic, symbols, TradesEvent, func(client *BinanceClient, symbol string) (api.ApiResponse, error) {
		return client.Trades(symbol, 0)
	})
}

func (caller *BinanceCaller) ScrapePrice(topic string, symbols []string) {
	caller.call(topic, symbols, PriceEvent, func(client *BinanceClient, symbol string) (api.ApiResponse, error) {
		return client.TickerPrice(symbol)
	})
}
