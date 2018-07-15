package main

import (
	"fmt"

	"github.com/ppincak/rysen/core"

	"github.com/ppincak/rysen/binance"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	client := binance.NewClient("https://api.binance.com", nil)
	// resp, _ := client.OrderBook("LTCBTC", 100)
	// fmt.Println(resp)
	// resp, _ = client.Trades("LTCBTC", 100)
	// fmt.Println(resp)
	//fmt.Println(client.ExchangeInfo())

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

	select {}
}
