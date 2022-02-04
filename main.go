package main

import (
	"runtime"

	"github.com/onrik/logrus/filename"

	"github.com/ppincak/rysen/crypto"
	"github.com/ppincak/rysen/services/aggregator"
	"github.com/ppincak/rysen/services/feed"
	"github.com/ppincak/rysen/services/publisher"
	"github.com/ppincak/rysen/services/schema"
	"github.com/ppincak/rysen/services/scraper"
	"github.com/ppincak/rysen/services/security"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/ppincak/rysen/monitor"
	b "github.com/ppincak/rysen/pkg/bus"
	"github.com/ppincak/rysen/pkg/ws"
	"github.com/ppincak/rysen/server"

	"github.com/ppincak/rysen/binance"

	log "github.com/sirupsen/logrus"
)

func main() {
	runtime.GOMAXPROCS(4)

	// Setup logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.AddHook(filename.NewHook())

	bus := b.NewBus()
	go bus.Start()

	binanceConfig := binance.NewConfig("https://api.binance.com")
	binanceExchange := binance.NewExchange(binanceConfig, bus)
	err := binanceExchange.Initialize()
	if err != nil {
		log.Error(err)
		return
	}

	exchanges := crypto.NewExchanges()
	exchanges.Register(binanceExchange)

	wsHandler := ws.NewHandler(nil)

	aggregatorService := aggregator.NewService(bus)
	scraperService := scraper.NewService(bus)

	db, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		log.Error(err)
		return
	}

	securityPersistence := security.NewPersistence(db, nil)
	accounts, err := securityPersistence.GetAccounts()
	if err != nil {
		log.Error(err)
		return
	}

	publisherService := publisher.NewService(bus, nil)

	securityService := security.NewService()
	securityService.Initialize(accounts)

	feedPersistence := feed.NewPersistence(db, nil)
	feeds, err := feedPersistence.GetFeeds()
	if err != nil {
		return
	}

	feedService := feed.NewService(bus, wsHandler)
	feedService.Initialize(feeds)

	// Register all monitors
	monitor := monitor.NewMonitor()
	monitor.Register(binanceExchange)
	monitor.Register(wsHandler)
	monitor.Register(feedService)

	schemaPersistence := schema.NewPersistence(db, nil)
	schemas, err := schemaPersistence.GetSchemas()
	if err != nil {
		return
	}

	schemaService := schema.NewService(
		aggregatorService,
		feedService,
		scraperService,
		exchanges)

	schemaService.Initialize(schemas)

	app := &server.App{
		Bus:       bus,
		Exchanges: exchanges,

		FeedService:     feedService,
		FeedPersistence: feedPersistence,

		PublisherService: publisherService,

		SecurityService:     securityService,
		SecurityPersistence: securityPersistence,

		SchemaService:     schemaService,
		SchemaPersistence: schemaPersistence,

		AggregatorService: aggregatorService,
		ScraperService:    scraperService,

		Monitor:   monitor,
		WsHandler: wsHandler,
	}

	s := server.NewServer(app, nil)
	s.Run()
}
