package main

import (
	"github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/monitor"
	"github.com/ppincak/rysen/pkg/ws"
	"github.com/ppincak/rysen/server"

	"github.com/ppincak/rysen/core"

	"github.com/ppincak/rysen/binance"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	client := binance.NewClient("https://api.binance.com", nil)

	apiCounter := core.NewApiCallCounter(10000, 10000)
	go apiCounter.Start()
	// eos := binance.NewScraper("eos", []string{"EOSBTC"}, caller.ScrapePrice, 1000)
	// go eos.Start()
	// ada := binance.NewScraper("ada", []string{"ADABTC"}, caller.ScrapePrice, 1000)
	// go ada.Start()

	// outc := make(chan *core.BusEvent)
	// bus.Subscribe("eos", outc)
	// bus.Subscribe("ada", outc)

	// go func(outc chan *core.BusEvent) {
	// 	fmt.Println("running")
	// 	for {
	// 		select {
	// 		case msg := <-outc:
	// 			fmt.Println(msg.Message)
	// 		}
	// 	}
	// }(outc)

	store := binance.NewStore()
	err := store.Initialize(client)
	if err != nil {
		log.Errorf("Failed to initialze BinanceStore")
		return
	}

	bus := bus.NewBus()
	go bus.Start()

	caller := binance.NewCaller(client, bus, apiCounter)

	wsHandler := ws.NewHandler(nil)

	monitor := monitor.NewMonitor()
	monitor.Register(client)
	monitor.Register(wsHandler)

	app := &server.App{
		Binance: &server.Binance{
			Client: client,
			Caller: caller,
			Store:  store,
		},
		Bus:       bus,
		Monitor:   monitor,
		WsHandler: wsHandler,
	}

	s := server.NewServer(app, nil)
	s.Run()
}
