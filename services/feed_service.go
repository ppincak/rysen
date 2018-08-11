package services

import (
	"sync"

	"github.com/ppincak/rysen/pkg/ws"

	"github.com/ppincak/rysen/bus"

	log "github.com/sirupsen/logrus"
)

type FeedService struct {
	bus     *bus.Bus
	feeds   map[string]*Feed
	handler *ws.Handler
	lock    *sync.RWMutex
}

// Create feed service
func NewFeedService(bus *bus.Bus, handler *ws.Handler) *FeedService {
	return &FeedService{
		bus:     bus,
		feeds:   make(map[string]*Feed),
		handler: handler,
		lock:    new(sync.RWMutex),
	}
}

func (service *FeedService) Create(metadata *FeedMetadata) *Feed {
	defer service.lock.Unlock()
	service.lock.Lock()

	feed := NewFeed(metadata, service.bus, service.handler, nil)
	log.Infof("Feed [%s] created", feed.Name)

	feed.Init()
	service.feeds[metadata.Name] = feed

	return feed
}

func (service *FeedService) SubscribeTo(name string, client *ws.Client) error {
	defer service.lock.Unlock()
	service.lock.Lock()

	if feed, ok := service.feeds[name]; ok {
		feed.subscribe(client)
	}
	return nil
}

func (service *FeedService) GetList() []*FeedMetadata {
	list := make([]*FeedMetadata, len(service.feeds))
	i := 0
	for _, value := range service.feeds {
		list[i] = value.FeedMetadata
		i++
	}
	return list
}
