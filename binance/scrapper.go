package binance

import (
	"time"

	"github.com/ppincak/rysen/core"
	log "github.com/sirupsen/logrus"
)

//
type Scraper struct {
	topic   string
	symbols []string
	fun     CallerFunc
	ticker  *time.Ticker
	outc    chan *CallerEvent
	stopc   chan struct{}
}

var _ core.Worker = (*Scraper)(nil)

// Returns new Scraper where the interval value is in miliseconds
func NewScraper(topic string,
	symbols []string,
	fun CallerFunc,
	interval int64) *Scraper {

	return &Scraper{
		topic:   topic,
		symbols: symbols,
		fun:     fun,
		ticker:  time.NewTicker(core.ToDurationMillis(interval)),
		stopc:   make(chan struct{}),
	}
}

// Start the scrapper
func (scraper *Scraper) Start() {
	for {
		select {
		case <-scraper.ticker.C:
			scraper.fun(scraper.topic, scraper.symbols)
		case <-scraper.stopc:
			log.Infof("Worker for topic: %s stopped", scraper.topic)
			return
		}
	}
}

// Stop the scrapper
func (worker *Scraper) Stop() {
	worker.stopc <- struct{}{}
}
