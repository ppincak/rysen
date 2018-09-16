package main

import (
	"github.com/ppincak/rysen/services/aggregator"
	"github.com/ppincak/rysen/services/feed"
	"github.com/ppincak/rysen/services/schema"
	"github.com/ppincak/rysen/services/scraper"

	b "github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/monitor"
	"github.com/ppincak/rysen/pkg/ws"
	"github.com/ppincak/rysen/server"

	"github.com/ppincak/rysen/binance"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	binanceConfig := binance.NewConfig("https://api.binance.com", nil)

	bus := b.NewBus()
	go bus.Start()

	exchange := binance.NewExchange(binanceConfig, bus)
	err := exchange.Initialize()
	if err != nil {
		log.Error(err)
		return
	}
	wsHandler := ws.NewHandler(nil)

	aggregatorService := aggregator.NewService(bus)
	feedService := feed.NewService(bus, wsHandler)
	scraperService := scraper.NewService(bus)

	monitor := monitor.NewMonitor()
	monitor.Register(exchange)
	monitor.Register(wsHandler)
	monitor.Register(feedService)

	schemas, err := schema.LoadAndCreateSchema("./schema.json")
	schemaService := schema.NewService(aggregatorService, feedService, scraperService)
	schemaService.Register("binance", exchange)
	schemaService.Create(schemas.Component("testSchema"))

	app := &server.App{
		Binance:           exchange,
		Bus:               bus,
		SchemaService:     schemaService,
		AggregatorService: aggregatorService,
		FeedService:       feedService,
		ScraperService:    scraperService,
		Monitor:           monitor,
		WsHandler:         wsHandler,
	}

	s := server.NewServer(app, nil)
	s.Run()
}
