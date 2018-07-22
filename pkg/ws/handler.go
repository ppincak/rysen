package ws

import (
	"net/http"
	"sync"

	"github.com/ppincak/rysen/api"

	"github.com/ppincak/rysen/core"
	"github.com/ppincak/rysen/monitor"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	clients  map[string]*Client
	config   *Config
	lock     *sync.RWMutex
	upgrader *websocket.Upgrader
	metrics  *core.WsMetrics
}

func NewHandler(config *Config) *Handler {
	if config == nil {
		config = DefaultConfig
	}

	return &Handler{
		clients: make(map[string]*Client),
		config:  config,
		lock:    new(sync.RWMutex),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  config.ReadBufferSize,
			WriteBufferSize: config.WriteBufferSize,
		},
		metrics: core.NewWsMetrics(),
	}
}

func (client *Handler) Statistics() []*monitor.Statistic {
	return []*monitor.Statistic{
		client.metrics.ToStatistic("wshandlerStatistics"),
	}
}

func (handler *Handler) ServeWebSocket(w http.ResponseWriter, r *http.Request) (*Client, error) {
	ws, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return nil, api.NewError("Failed to upgrade connection to Websocket")
	}

	client := NewClient(ws, handler)
	go client.readPump()
	go client.writePump()

	handler.addClient(client)

	log.Infof("Client [%s] connection established", client.uuid)

	return client, nil
}

func (handler *Handler) addClient(client *Client) {
	defer handler.lock.Unlock()
	handler.lock.Lock()

	handler.clients[client.uuid] = client

	log.Infof("Client [%s] added to handler collection")
}

func (handler *Handler) removeClient(client *Client) {
	defer handler.lock.Unlock()
	handler.lock.Lock()

	if _, ok := handler.clients[client.uuid]; ok {
		delete(handler.clients, client.uuid)

		log.Infof("Client [%s] removed from handler collection")
	}
}
