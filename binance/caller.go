package binance

import (
	"time"

	"github.com/ppincak/rysen/binance/model"

	"github.com/ppincak/rysen/pkg/bus"
	"github.com/ppincak/rysen/pkg/scrape"

	log "github.com/sirupsen/logrus"
)

type Caller struct {
	client *Client
	bus    *bus.Bus
}

var _ scrape.Caller = (*Caller)(nil)

func NewCaller(client *Client, bus *bus.Bus) *Caller {
	return &Caller{
		client: client,
		bus:    bus,
	}
}

type Callable func(*Client, string) (model.Model, error)

func (caller *Caller) call(topic string, symbols []string, eventType scrape.CallerFuncType, callable Callable) {
	// Note: check if this affects performance, if yes refactor this
	workFunc := func(client *Client, symbol string) {
		t := time.Now()

		resp, err := callable(client, symbol)

		if err == nil {
			caller.bus.Publish(topic, &scrape.CallerEvent{
				Symbol:    symbol,
				EventType: eventType,
				Data:      resp,
				Timestamp: t.UnixNano(),
			})
		} else {
			log.Errorf("Error during scraping call")
			log.Error(err)
		}
	}

	if len(symbols) == 0 {
		go workFunc(caller.client, "")
	} else {
		for _, symbol := range symbols {
			go workFunc(caller.client, symbol)
		}
	}
}

func (caller *Caller) ScrapeOrders(topic string, symbols []string) {
	caller.call(topic, symbols, scrape.OrdersEvent, func(client *Client, symbol string) (model.Model, error) {
		response, err := client.OrderBook(symbol, 0)
		var book *model.OrderBookProcessed
		if err == nil {
			if asserted, ok := response.(*model.OrderBook); ok {
				book = asserted.ProcessAll()
			}
		}
		return book, err
	})
}

func (caller *Caller) ScrapeTrades(topic string, symbols []string) {
	caller.call(topic, symbols, scrape.TradesEvent, func(client *Client, symbol string) (model.Model, error) {
		return client.Trades(symbol, 0)
	})
}

func (caller *Caller) ScrapePrice(topic string, symbols []string) {
	caller.call(topic, symbols, scrape.PriceEvent, func(client *Client, symbol string) (model.Model, error) {
		return client.TickerPrice(symbol)
	})
}
