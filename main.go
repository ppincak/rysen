package main

import (
	"fmt"

	"github.com/ppincak/rysen/monitor"

	"github.com/ppincak/rysen/core"

	"github.com/ppincak/rysen/binance"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	monitor := monitor.NewMonitor()
	client := binance.NewClient("https://api.binance.com", nil)

	apiCounter := core.NewApiCallCounter(10000, 10000)
	go apiCounter.Start()
	bus := binance.NewBinanceBus()
	go bus.Start()
	caller := binance.NewCaller(client, bus, apiCounter)
	eos := binance.NewScraper("eos", []string{"EOSBTC"}, caller.ScrapePrice, 1000)
	go eos.Start()
	ada := binance.NewScraper("ada", []string{"ADABTC"}, caller.ScrapePrice, 1000)
	go ada.Start()

	outc := make(chan *core.BusEvent)
	bus.Subscribe("eos", outc)
	bus.Subscribe("ada", outc)

	go func(outc chan *core.BusEvent) {
		fmt.Println("running")
		for {
			select {
			case msg := <-outc:
				fmt.Println(msg.Message)
			}
		}
	}(outc)

	monitor.Register(client)

	select {}
}
