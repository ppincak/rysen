package feed

import (
	"sync"

	"rysen/pkg/bus"
	"rysen/pkg/ws"

	log "github.com/sirupsen/logrus"
)

type Feed struct {
	*Model

	// map of clients connected to the feed
	clients map[string]*ws.Client
	// bus event used for subscription to bus
	eventsc chan *bus.BusEvent
	// bus
	bus *bus.Bus
	// ws handler
	handler *ws.Handler
	// uuid of the on remove handler
	onRemoveUuid string
	// mutex
	lock *sync.RWMutex
	// transformer
	transformer *bus.Transformer
	// transformation function
	transformFunc bus.TransformFunc
	// subscription to bus
	sub *bus.BusSubscription

	// feeds metrics
	metrics *Metrics
}

// Create Feed
func NewFeed(model *Model, b *bus.Bus, handler *ws.Handler, transformFunc bus.TransformFunc) *Feed {
	if transformFunc == nil {
		transformFunc = TransformForWsClient
	}

	return &Feed{
		Model:         model,
		bus:           b,
		clients:       make(map[string]*ws.Client),
		eventsc:       make(chan *bus.BusEvent),
		handler:       handler,
		lock:          new(sync.RWMutex),
		transformFunc: transformFunc,
		metrics:       NewMetrics(),
	}
}

// Transform bus message to message for the feed
func (feed *Feed) transform(event *bus.BusEvent) {
	defer feed.lock.RUnlock()
	feed.lock.RLock()

	for _, client := range feed.clients {
		// Note: rethink and maybe refactor in the future
		// call has to be async, so we don't block the route for a log period of time
		go func(client *ws.Client, name string, message interface{}, transformFunc bus.TransformFunc) {
			if transformed, err := transformFunc(message); err == nil {
				client.WriteEvent(name, transformed)
			}
		}(client, feed.Name, event.Message, feed.transformFunc)
	}

	feed.metrics.MessagesSent.Add(int64(len(feed.clients)))
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
	feed.transformer = bus.NewTransformer(feed.eventsc, feed.transform)
	go feed.transformer.Start()

	feed.sub = feed.bus.Subscribe(feed.Topic, feed.eventsc)
	feed.onRemoveUuid = feed.handler.OnRemove(feed.unsubscribe)

	log.Infof("Feed [%s] initialized", feed.Name)
}

// Destroy the feed
func (feed *Feed) Destroy() {
	defer feed.lock.Unlock()
	feed.lock.Lock()

	feed.sub.Cancel()
	feed.transformer.Stop()

	log.Infof("Feed [%s] destroyed", feed.Name)
}
