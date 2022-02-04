package scrape

import (
	"time"

	log "github.com/sirupsen/logrus"
	"rysen/core"
)

// Scrapes from binance
type Scraper struct {
	topic    string
	symbols  []string
	fun      CallerFunc
	ticker   *time.Ticker
	interval int64
	outc     chan *CallerEvent
	stopc    chan struct{}
}

var _ core.Worker = (*Scraper)(nil)

// Returns new Scraper where the interval value is in miliseconds
func NewScraper(topic string,
	symbols []string,
	fun CallerFunc,
	interval int64) *Scraper {

	return &Scraper{
		topic:    topic,
		symbols:  symbols,
		fun:      fun,
		ticker:   time.NewTicker(core.ToDurationMillis(interval)),
		interval: interval,
		stopc:    make(chan struct{}),
	}
}

// Return scraper topic
func (scraper *Scraper) Topic() string {
	return scraper.topic
}

// Return scraper symbols
func (scraper *Scraper) Symbols() []string {
	return scraper.symbols
}

// Return ticker interval
func (scraper *Scraper) Interval() int64 {
	return scraper.interval
}

// Start the scrapper
func (scraper *Scraper) Start() {
	defer func() {
		log.Infof("Worker for topic: [%s] stopped", scraper.topic)
	}()

	for {
		select {
		case <-scraper.ticker.C:
			scraper.fun(scraper.topic, scraper.symbols)
		case <-scraper.stopc:
			return
		}
	}
}

// Stop the scrapper
func (worker *Scraper) Stop() {
	worker.stopc <- struct{}{}
}
