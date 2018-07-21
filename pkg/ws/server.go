package ws

import (
	"net/http"
	"sync"

	"github.com/ppincak/rysen/core"
	"github.com/ppincak/rysen/monitor"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WsServer struct {
	clients  map[string]*Client
	config   *WsConfig
	lock     *sync.RWMutex
	upgrader *websocket.Upgrader
	metrics  *core.WsMetrics
}

func NewWsServer(config *WsConfig) *WsServer {
	if config == nil {
		config = DefaultConfig
	}

	return &WsServer{
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

func (client *WsServer) Statistics() []*monitor.Statistic {
	return []*monitor.Statistic{
		client.metrics.ToStatistic("wsServerStatistics"),
	}
}

func (server *WsServer) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}

	client := NewClient(ws, server)
	go client.readPump()
	go client.writePump()

	server.addClient(client)

	log.Infof("Client [%s] connection established", client.uuid)
}

func (server *WsServer) addClient(client *Client) {
	defer server.lock.Unlock()
	server.lock.Lock()

	server.clients[client.uuid] = client

	log.Infof("Client [%s] added to server collection")
}

func (server *WsServer) removeClient(client *Client) {
	defer server.lock.Unlock()
	server.lock.Lock()

	if _, ok := server.clients[client.uuid]; ok {
		delete(server.clients, client.uuid)

		log.Infof("Client [%s] removed from server collection")
	}
}
