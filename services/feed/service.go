package feed

import (
	"sync"

	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/monitor"

	"github.com/ppincak/rysen/pkg/ws"

	"github.com/ppincak/rysen/pkg/bus"

	log "github.com/sirupsen/logrus"
)

var _ monitor.Reporter = (*Service)(nil)

type Service struct {
	bus     *bus.Bus
	feeds   map[string]*Feed
	handler *ws.Handler
	lock    *sync.RWMutex
}

// Create feed service
func NewService(bus *bus.Bus, handler *ws.Handler) *Service {
	return &Service{
		bus:     bus,
		feeds:   make(map[string]*Feed),
		handler: handler,
		lock:    new(sync.RWMutex),
	}
}

// Get feed statistics
func (service *Service) Statistics() []*monitor.Statistic {
	defer service.lock.Unlock()
	service.lock.Lock()

	i := 0
	statistics := make([]*monitor.Statistic, len(service.feeds))
	for _, feed := range service.feeds {
		statistics[i] = feed.metrics.ToStatistic(feed.Name)
		i++
	}
	return statistics
}

// Note: maybe check for unique feed name
// Create feed
func (service *Service) Create(metadata *Metadata) *Feed {
	defer service.lock.Unlock()
	service.lock.Lock()

	feed := NewFeed(metadata, service.bus, service.handler, nil)
	log.Infof("Feed [%s] created", feed.Name)

	feed.Init()
	service.feeds[metadata.Name] = feed

	return feed
}

// Subscribe to feed
func (service *Service) SubscribeTo(name string, client *ws.Client) error {
	defer service.lock.Unlock()
	service.lock.Lock()

	if feed, ok := service.feeds[name]; ok {
		feed.subscribe(client)
		return nil
	} else {
		return api.NewError("Feed not found")
	}
}

// List all feeds
func (service *Service) ListFeeds() []*Metadata {
	list := make([]*Metadata, len(service.feeds))
	i := 0
	for _, value := range service.feeds {
		list[i] = value.Metadata
		i++
	}
	return list
}
