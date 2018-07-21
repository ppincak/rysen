package ws

import (
	"encoding/json"
	"time"

	"github.com/ppincak/rysen/api"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	uuid       string
	config     *WsConfig
	server     *WsServer
	conn       *websocket.Conn
	writec     chan []byte
	readc      chan []byte
	stopc      chan struct{}
	pongTicker *time.Ticker
}

func NewClient(writec *websocket.Conn, server *WsServer) *Client {
	return &Client{
		uuid:       uuid.New().String(),
		config:     server.config,
		server:     server,
		writec:     make(chan []byte),
		readc:      make(chan []byte),
		stopc:      make(chan struct{}),
		pongTicker: time.NewTicker(server.config.PingWait),
	}
}

func (client *Client) GetSessionId() string {
	return client.uuid
}

func (client *Client) Write(topic string, message interface{}) error {
	packet := &Packet{
		Id:      uuid.New().String(),
		Topic:   topic,
		Message: message,
	}
	bytes, err := json.Marshal(packet)
	if err != nil {
		return api.NewError("Failed to marshall message")
	}
	client.writec <- bytes
	return nil
}

func (client *Client) Read(message []byte) {
	client.readc <- message
}

func (client *Client) stop() {
	client.stopc <- struct{}{}
}

func (client *Client) readPump() {
	defer func() {
		log.Infof("Client: %s readpump stopped", client.uuid)
	}()

	log.Infof("Client: %s readpump started", client.uuid)

	conn := client.conn
	conn.SetReadLimit(client.config.MaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(client.config.ReadWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(client.config.PongWait))
		return nil
	})

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Errorf("Failed to read from Client: %s", client.uuid)
			return
		}
		client.readc <- message
	}
}

func (client *Client) writePump() {
	defer func() {
		log.Infof("Client: %s writepump stopped", client.uuid)
	}()

	conn := client.conn

	for {
		select {
		case message, ok := <-client.writec:
			conn.SetWriteDeadline(time.Now().Add(client.config.WriteWait))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})

				log.Errorf("Failed to received message for Client %s", client.uuid)
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Errorf("Failed to obtain writer for Client %s", client.uuid)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Error("Failed to write message", message)
			}
		case <-client.pongTicker.C:
			conn.SetWriteDeadline(time.Now().Add(client.config.WriteWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
