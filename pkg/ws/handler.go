package ws

import (
	"net/http"
	"sync"

	"github.com/google/uuid"

	"github.com/ppincak/rysen/monitor"
	"github.com/ppincak/rysen/pkg/errors"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type OnRemoveCallback func(*Client)

// Connection handler
type Handler struct {
	clients           map[string]*Client
	config            *Config
	lock              *sync.RWMutex
	upgrader          *websocket.Upgrader
	metrics           *WsMetrics
	onRemoveCallbacks map[string]OnRemoveCallback
}

// Create new handler
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
			CheckOrigin:     CheckOriginAny,
		},
		metrics:           NewWsMetrics(),
		onRemoveCallbacks: make(map[string]OnRemoveCallback),
	}
}

func (client *Handler) Statistics() []*monitor.Statistic {
	return []*monitor.Statistic{
		client.metrics.ToStatistic("wsHandlerStatistics"),
	}
}

func (handler *Handler) ServeWebSocket(w http.ResponseWriter, r *http.Request) (*Client, error) {
	ws, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return nil, errors.NewError("Failed to upgrade connection to Websocket")
	}

	client := NewClient(ws, handler)
	go client.readPump()
	go client.writePump()
	// blocks till writepump isn't started
	client.sendConnectPacket()
	client.ping()

	handler.addClient(client)
	handler.metrics.Connects.Inc()

	log.Infof("Client [%s] connection established", client.uuid)

	return client, nil
}

func (handler *Handler) GetClient(sessionId string) *Client {
	defer handler.lock.RUnlock()
	handler.lock.RLock()

	if client, ok := handler.clients[sessionId]; ok {
		return client
	}
	return nil
}

func (handler *Handler) OnRemove(callback OnRemoveCallback) string {
	defer handler.lock.Unlock()
	handler.lock.Lock()

	uuid := uuid.New().String()
	handler.onRemoveCallbacks[uuid] = callback
	return uuid
}

func (handler *Handler) runOnRemove(client *Client) {
	// Note: maybe needed to add a lock here in the future
	for _, callback := range handler.onRemoveCallbacks {
		callback(client)
	}
}

func (handler *Handler) addClient(client *Client) {
	defer handler.lock.Unlock()
	handler.lock.Lock()

	handler.clients[client.uuid] = client
	handler.metrics.Clients.Inc()

	log.Infof("Client [%s] added to handler collection", client.uuid)
}

func (handler *Handler) removeClient(client *Client) {
	defer handler.lock.Unlock()
	handler.lock.Lock()

	if _, ok := handler.clients[client.uuid]; ok {
		delete(handler.clients, client.uuid)
		handler.metrics.Clients.Dec()

		go handler.runOnRemove(client)

		log.Infof("Client [%s] removed from handler collection", client.uuid)
	}
}
