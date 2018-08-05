package services

import (
	"sync"

	"github.com/ppincak/rysen/pkg/ws"

	"github.com/ppincak/rysen/bus"
	"github.com/ppincak/rysen/pkg/scrape"

	log "github.com/sirupsen/logrus"
)

type FeedService struct {
	bus     *bus.Bus
	feeds   map[string]*Feed
	handler *ws.Handler
	lock    *sync.RWMutex
}

func NewFeedService(bus *bus.Bus, handler *ws.Handler) *FeedService {
	return &FeedService{
		bus:     bus,
		feeds:   make(map[string]*Feed),
		handler: handler,
		lock:    new(sync.RWMutex),
	}
}

func (service *FeedService) Create(name string, scraper *scrape.Scraper) *Feed {
	defer service.lock.Unlock()
	service.lock.Lock()

	feed := NewFeed(service.bus, service.handler, name, scraper)
	log.Infof("Feed [%s] created", feed.Name)

	feed.Init()
	service.feeds[name] = feed

	return feed
}

func (service *FeedService) SubscribeTo(name string, client *ws.Client) {
	defer service.lock.Unlock()
	service.lock.Lock()

	if feed, ok := service.feeds[name]; ok {
		feed.subscribe(client)
	}
}

func (service *FeedService) GetList() []*Feed {
	list := make([]*Feed, len(service.feeds))
	i := 0
	for _, value := range service.feeds {
		list[i] = value
		i++
	}
	return list
}
