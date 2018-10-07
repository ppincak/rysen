package feed

import (
	"sync"

	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/monitor"

	"github.com/ppincak/rysen/pkg/ws"

	"github.com/ppincak/rysen/pkg/bus"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	bus         *bus.Bus
	feeds       map[string]*Feed
	clientFeeds map[string][]*Feed
	handler     *ws.Handler
	lock        *sync.RWMutex
}

var _ monitor.Reporter = (*Service)(nil)

// Create feed service
func NewService(bus *bus.Bus, handler *ws.Handler) *Service {
	return &Service{
		bus:         bus,
		feeds:       make(map[string]*Feed),
		clientFeeds: make(map[string][]*Feed),
		handler:     handler,
		lock:        new(sync.RWMutex),
	}
}

// Initialize the service
func (service *Service) Initialize(feeds []*Metadata) (err error) {
	for _, feed := range feeds {
		_, err = service.Create(feed)
		if err != nil {
			return
		}
	}
	return
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

// Create feed
func (service *Service) Create(metadata *Metadata) (*Feed, error) {
	defer service.lock.Unlock()
	service.lock.Lock()

	if _, ok := service.feeds[metadata.Name]; ok {
		return nil, api.NewError("Feed with name [%s] already exists", metadata.Name)
	}

	feed := NewFeed(metadata, service.bus, service.handler, nil)
	log.Infof("Feed [%s] created", feed.Name)

	feed.Init()
	service.feeds[metadata.Name] = feed

	return feed, nil
}

// Subscribe to feed
func (service *Service) SubscribeTo(name string, client *ws.Client) error {
	defer service.lock.Unlock()
	service.lock.Lock()

	if feed, ok := service.feeds[name]; ok {
		feed.subscribe(client)

		// Add feed to collection of feeds
		var feeds []*Feed
		if f, ok := service.clientFeeds[client.GetSessionId()]; !ok {
			feeds = make([]*Feed, 0)
		} else {
			feeds = f
		}
		feeds = append(feeds, feed)

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

// List all client feeds
func (service *Service) ListClientFeeds(sessionId string) []*Metadata {
	list := make([]*Metadata, 0)
	if feeds, ok := service.clientFeeds[sessionId]; ok {
		for _, feed := range feeds {
			list = append(list, feed.Metadata)
		}
	}
	return list
}
