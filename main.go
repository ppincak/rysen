package main

import (
	"github.com/ppincak/rysen/api"

	"github.com/ppincak/rysen/binance/data"
	b "github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/monitor"
	"github.com/ppincak/rysen/pkg/aggregate"
	"github.com/ppincak/rysen/pkg/scrape"
	"github.com/ppincak/rysen/pkg/ws"
	"github.com/ppincak/rysen/server"
	"github.com/ppincak/rysen/services"

	"github.com/ppincak/rysen/binance"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	binanceConfig := &binance.Config{
		Url:    "https://api.binance.com",
		Secret: nil,
	}
	bus := b.NewBus()
	go bus.Start()

	exchange := binance.NewExchange(binanceConfig, bus)
	err := exchange.Initialize()
	if err != nil {
		log.Error(err)
		return
	}
	wsHandler := ws.NewHandler(nil)

	feedService := services.NewFeedService(bus, wsHandler)

	monitor := monitor.NewMonitor()
	monitor.Register(exchange)
	monitor.Register(wsHandler)

	app := &server.App{
		Binance:     exchange,
		Bus:         bus,
		FeedService: feedService,
		Monitor:     monitor,
		WsHandler:   wsHandler,
	}

	eos := scrape.NewScraper("eos", []string{"EOSBTC"}, exchange.Caller.ScrapePrice, 1000)
	go eos.Start()
	ada := scrape.NewScraper("ada", []string{"ADABTC"}, exchange.Caller.ScrapePrice, 1000)
	go ada.Start()

	eosTrades := scrape.NewScraper("eos/trades", []string{"EOSBTC"}, exchange.Caller.ScrapeTrades, 1000)
	go eosTrades.Start()

	feedService.Create(services.NewFeedMetadata("eos", "eosPrice", ""))
	feedService.Create(services.NewFeedMetadata("ada", "adaPrice", ""))
	feedService.Create(services.NewFeedMetadata("eos/trades/aggregate", "eosTrades", ""))

	aggregator := aggregate.NewAggregator("eos/trades", "eos/trades/aggregate", bus, ProcessCallerEvent, data.SumTrades, aggregate.AggretateTillSize(5))
	go aggregator.Start()

	s := server.NewServer(app, nil)
	s.Run()
}

func ProcessCallerEvent(event *b.BusEvent) (interface{}, error) {
	if assertion, ok := event.Message.(*scrape.CallerEvent); ok {
		return assertion.Data, nil
	} else {
		return nil, api.NewError("Failed to assert")
	}
}
