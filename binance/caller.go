package binance

import (
	"time"

	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/core"
	"github.com/ppincak/rysen/pkg/scrape"
)

type Caller struct {
	client  *Client
	bus     *bus.Bus
	limiter *core.CallLimiter
}

func NewCaller(client *Client, bus *bus.Bus, limiter *core.CallLimiter) *Caller {
	return &Caller{
		client:  client,
		bus:     bus,
		limiter: limiter,
	}
}

type Callable func(*Client, string) (api.ApiResponse, error)

func (caller *Caller) call(topic string, symbols []string, eventType scrape.CallerFuncType, callable Callable) {
	for _, symbol := range symbols {
		go func(client *Client, symbol string) {
			t := time.Now()

			resp, err := callable(client, symbol)

			go caller.limiter.IncrRequests(nil)

			if err == nil {
				caller.bus.Publish(topic, &scrape.CallerEvent{
					Symbol:    symbol,
					EventType: eventType,
					Data:      resp,
					Timestamp: t.UnixNano(),
				})
			}
		}(caller.client, symbol)
	}
}

func (caller *Caller) ScrapeOrders(topic string, symbols []string) {
	caller.call(topic, symbols, scrape.OrdersEvent, func(client *Client, symbol string) (api.ApiResponse, error) {
		return client.OrderBook(symbol, 0)
	})
}

func (caller *Caller) ScrapeTrades(topic string, symbols []string) {
	caller.call(topic, symbols, scrape.TradesEvent, func(client *Client, symbol string) (api.ApiResponse, error) {
		return client.Trades(symbol, 0)
	})
}

func (caller *Caller) ScrapePrice(topic string, symbols []string) {
	caller.call(topic, symbols, scrape.PriceEvent, func(client *Client, symbol string) (api.ApiResponse, error) {
		return client.TickerPrice(symbol)
	})
}
