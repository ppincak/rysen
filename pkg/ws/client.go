package ws

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"rysen/pkg/errors"
)

type Client struct {
	uuid       string
	config     *Config
	handler    *Handler
	conn       *websocket.Conn
	writec     chan []byte
	readc      chan []byte
	stopc      chan struct{}
	pongTicker *time.Ticker
}

func NewClient(conn *websocket.Conn, handler *Handler) *Client {
	return &Client{
		uuid:       uuid.New().String(),
		config:     handler.config,
		conn:       conn,
		handler:    handler,
		writec:     make(chan []byte),
		readc:      make(chan []byte),
		stopc:      make(chan struct{}),
		pongTicker: time.NewTicker(handler.config.PingWait),
	}
}

// Get session id
func (client *Client) GetSessionId() string {
	return client.uuid
}

// Send message to clients write channel
func (client *Client) WriteEvent(event string, message interface{}) error {
	packet := &Packet{
		Id:         uuid.New().String(),
		PacketType: Event,
		Content: EventMessage{
			Event:   event,
			Message: message,
		},
	}
	return client.Write(packet)
}

// Send message to clients write channel
func (client *Client) Write(packet *Packet) error {
	bytes, err := json.Marshal(packet)
	if err != nil {
		log.Error("Failed to marshall message")

		return errors.NewError("Failed to marshall message")
	}
	client.writec <- bytes
	return nil
}

// Send WS ping packet
func (client *Client) ping() error {
	client.conn.SetWriteDeadline(time.Now().Add(client.config.WriteWait))
	if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		log.Error("Failed to send ping request")
		return err
	}
	return nil
}

// TODO refactor, should read Packet
// Send message to clients read channel
func (client *Client) Read(message []byte) {
	client.readc <- message
}

// Stop the client
func (client *Client) stop() {
	client.stopc <- struct{}{}
}

// Sends connect packet
func (client *Client) sendConnectPacket() {
	client.Write(&Packet{
		Id:         uuid.New().String(),
		PacketType: Connect,
		Content: map[string]string{
			"clientId": client.uuid,
		},
	})
}

func (client *Client) readPump() {
	defer func() {
		client.handler.removeClient(client)

		log.Infof("Client: [%s] readpump stopped", client.uuid)
	}()

	log.Infof("Client: [%s] readpump started", client.uuid)

	conn := client.conn
	conn.SetReadLimit(client.config.MaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(client.config.PongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(client.config.PongWait))
		return nil
	})

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Errorf("Failed to read from Client: [%s]", client.uuid)
			log.Error(err)

			client.handler.metrics.ReadsFailed.Inc()

			return
		}
		client.readc <- message
	}
}

func (client *Client) writePump() {
	defer func() {
		client.handler.removeClient(client)
		client.handler.metrics.Disconnects.Inc()

		log.Infof("Client [%s] writepump stopped", client.uuid)
	}()

	log.Infof("Client: [%s] writepump started", client.uuid)

	conn := client.conn

	for {
		select {
		case message, ok := <-client.writec:
			conn.SetWriteDeadline(time.Now().Add(client.config.WriteWait))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})

				log.Errorf("Failed to receive message for Client [%s]", client.uuid)
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Errorf("Failed to obtain writer for Client [%s]", client.uuid)
				log.Error(err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Error("Failed to write message", message)
				log.Error(err)

				client.handler.metrics.WritesFailed.Inc()
			}
		case <-client.pongTicker.C:
			if client.ping() != nil {
				return
			}
		}
	}
}
