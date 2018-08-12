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

	"github.com/ppincak/rysen/core"

	"github.com/ppincak/rysen/binance"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	client := binance.NewClient("https://api.binance.com", nil)
	apiCounter := core.NewApiCallCounter(10000, 10000)
	go apiCounter.Start()

	store := binance.NewStore()
	err := store.Initialize(client)
	if err != nil {
		log.Errorf("Failed to initialze BinanceStore")
		return
	}

	bus := b.NewBus()
	go bus.Start()

	caller := binance.NewCaller(client, bus, apiCounter)

	wsHandler := ws.NewHandler(nil)

	feedService := services.NewFeedService(bus, wsHandler)

	monitor := monitor.NewMonitor()
	monitor.Register(client)
	monitor.Register(wsHandler)

	app := &server.App{
		Binance: &server.Binance{
			Client: client,
			Caller: caller,
			Store:  store,
		},
		Bus:         bus,
		FeedService: feedService,
		Monitor:     monitor,
		WsHandler:   wsHandler,
	}

	eos := scrape.NewScraper("eos", []string{"EOSBTC"}, caller.ScrapePrice, 1000)
	go eos.Start()
	ada := scrape.NewScraper("ada", []string{"ADABTC"}, caller.ScrapePrice, 1000)
	go ada.Start()

	eosTrades := scrape.NewScraper("eos/trades", []string{"EOSBTC"}, caller.ScrapeTrades, 1000)
	go eosTrades.Start()

	feedService.Create(services.NewFeedMetadata("eos", "eosPrice", ""))
	feedService.Create(services.NewFeedMetadata("ada", "adaPrice", ""))
	feedService.Create(services.NewFeedMetadata("eos/trades", "eosTrades", ""))

	aggregator := b.NewAggregator("eos/trades", "aggregate", bus, ProcessCallerEvent, data.SumTrades, aggregate.AggretateTillSize(5))
	go aggregator.Start()

	feedService.Create(services.NewFeedMetadata("eos-aggregate", "eosAggregate", ""))

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
