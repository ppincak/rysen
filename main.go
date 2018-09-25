package main

import (
	"flag"

	"github.com/ppincak/rysen/crypto"
	"github.com/ppincak/rysen/services/aggregator"
	"github.com/ppincak/rysen/services/feed"
	"github.com/ppincak/rysen/services/schema"
	"github.com/ppincak/rysen/services/scraper"

	"github.com/ppincak/rysen/monitor"
	b "github.com/ppincak/rysen/pkg/bus"
	"github.com/ppincak/rysen/pkg/ws"
	"github.com/ppincak/rysen/server"

	"github.com/ppincak/rysen/binance"

	log "github.com/sirupsen/logrus"
)

func iniFlags() (string, string, string) {
	schemaFile := flag.String("schema", "", "Schema json file")
	accountsFile := flag.String("accounts", "", "Accounts json file")
	key := flag.String("decryptionKey", "", "Decryption key")

	flag.Parse()

	return *schemaFile, *accountsFile, *key
}

func iniBinance(bus *b.Bus) *binance.Exchange {
	binanceConfig := binance.NewConfig("https://api.binance.com")
	exchange := binance.NewExchange(binanceConfig, bus)
	err := exchange.Initialize()
	if err != nil {
		panic(err)
	}
	return exchange
}

func main() {
	//schemaFile, accountsFile, decryptionKey := iniFlags()

	// Setup logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	bus := b.NewBus()
	go bus.Start()

	binanceExchange := iniBinance(bus)

	exchanges := crypto.NewExchanges()
	exchanges.Register(binanceExchange)

	wsHandler := ws.NewHandler(nil)

	aggregatorService := aggregator.NewService(bus)
	feedService := feed.NewService(bus, wsHandler)
	scraperService := scraper.NewService(bus)
	schemaService := schema.NewService(aggregatorService, feedService, scraperService, exchanges)

	// Register all monitors
	monitor := monitor.NewMonitor()
	monitor.Register(binanceExchange)
	monitor.Register(wsHandler)
	monitor.Register(feedService)

	schemas, _ := schema.LoadAndCreateSchema("./schema.json")
	schemaService.Create(schemas.Component("testSchema"))

	schemaBackup := schema.NewSchemaBackup("./db", nil)
	schemaBackup.Open()
	schemaBackup.SaveSchema(schemas.Component("testSchema"))

	app := &server.App{
		Bus:               bus,
		Exchanges:         exchanges,
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
