package services

import (
	"sync"

	"github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/pkg/scrape"
	"github.com/ppincak/rysen/pkg/ws"

	log "github.com/sirupsen/logrus"
)

type Feed struct {
	Name        string   `json:"name"`
	Symbols     []string `json:"symbols"`
	bus         *bus.Bus
	clients     map[string]*ws.Client
	eventsc     chan *bus.BusEvent
	handler     *ws.Handler
	handlerUuid string
	interval    int64
	lock        *sync.RWMutex
	transformer *bus.Transformer
	sub         *bus.BusSubscription
	scraper     *scrape.Scraper
}

func NewFeed(b *bus.Bus, handler *ws.Handler, name string, scraper *scrape.Scraper) *Feed {
	return &Feed{
		Name:     name,
		Symbols:  scraper.Symbols(),
		bus:      b,
		clients:  make(map[string]*ws.Client),
		eventsc:  make(chan *bus.BusEvent),
		handler:  handler,
		interval: scraper.Interval(),
		lock:     new(sync.RWMutex),
		scraper:  scraper,
	}
}

func (feed *Feed) transform(event *bus.BusEvent) {
	defer feed.lock.RUnlock()
	feed.lock.RLock()

	for _, client := range feed.clients {
		go func(client *ws.Client, name string, message interface{}) {
			switch message.(type) {
			case *scrape.CallerEvent:
				client.WriteEvent(name, message.(scrape.CallerEvent).Data)
			}

		}(client, feed.Name, event.Message)
	}
}

// Subscribe client to feed
func (feed *Feed) subscribe(client *ws.Client) {
	defer feed.lock.Unlock()
	feed.lock.Lock()

	if _, ok := feed.clients[client.GetSessionId()]; !ok {
		feed.clients[client.GetSessionId()] = client
	}

	log.Infof("Client [%s] subscribed to feed [%s]", client.GetSessionId(), feed.Name)
}

// Unsubscribe client from feed
func (feed *Feed) unsubscribe(client *ws.Client) {
	defer feed.lock.Unlock()
	feed.lock.Lock()

	if _, ok := feed.clients[client.GetSessionId()]; ok {
		delete(feed.clients, client.GetSessionId())

		log.Infof("Client [%s] unsubscribed from feed [%s]", client.GetSessionId(), feed.Name)
	}
}

// Initialize the feed
func (feed *Feed) Init() {
	// Note: add locks ?

	feed.transformer = bus.NewTransformer(feed.eventsc, feed.transform)
	go feed.transformer.Start()

	feed.sub = feed.bus.Subscribe(feed.scraper.Topic(), feed.eventsc)
	feed.handlerUuid = feed.handler.OnRemove(feed.unsubscribe)

	log.Infof("Feed [%s] initialized", feed.Name)
}

// Destroy the feed
func (feed *Feed) Destroy() {
	// Note: add locks ?

	feed.sub.Cancel()
	feed.transformer.Stop()

	log.Infof("Feed [%s] destroyed", feed.Name)
}
