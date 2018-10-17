package feed

import (
	"sync"

	"github.com/ppincak/rysen/monitor"

	"github.com/ppincak/rysen/pkg/errors"
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
func (service *Service) Initialize(feeds []*Model) (err error) {
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
func (service *Service) Create(model *Model) (*Feed, error) {
	defer service.lock.Unlock()
	service.lock.Lock()

	if _, ok := service.feeds[model.Name]; ok {
		return nil, errors.NewError("Feed with name [%s] already exists", model.Name)
	}

	feed := NewFeed(model, service.bus, service.handler, nil)
	log.Infof("Feed [%s] created", feed.Name)

	feed.Init()
	service.feeds[model.Name] = feed

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
		return errors.NewError("Feed not found")
	}
}

// Unsubscribe from feed
func (service *Service) UnsubscribeFrom(name string, client *ws.Client) error {
	defer service.lock.Unlock()
	service.lock.Lock()

	if feed, ok := service.feeds[name]; ok {
		delete(service.clientFeeds, client.GetSessionId())

		feed.unsubscribe(client)

		return nil
	} else {
		return errors.NewError("Feed not found")
	}
}

// List all feeds
func (service *Service) ListFeeds() []*Model {
	list := make([]*Model, len(service.feeds))
	i := 0
	for _, value := range service.feeds {
		list[i] = value.Model
		i++
	}
	return list
}

// List all client feeds
func (service *Service) ListClientFeeds(sessionId string) []*Model {
	list := make([]*Model, 0)
	if feeds, ok := service.clientFeeds[sessionId]; ok {
		for _, feed := range feeds {
			list = append(list, feed.Model)
		}
	}
	return list
}
