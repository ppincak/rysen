package ws

import (
	"time"
)

const (
	ReadBufferSize  = 1024
	WriteBufferSize = 1024
	MaxMessageSize  = 1024 * 20
	ReadWait        = time.Second
	WriteWait       = time.Second / 2
	PongWait        = time.Second * 60
	PingWait        = time.Second - (PongWait / 10)
)

type WsConfig struct {
	Port            int           `json:"port"`
	ReadBufferSize  int           `json:"readBufferSize"`
	WriteBufferSize int           `json:"writeBufferSize"`
	MaxMessageSize  int64         `json:"maxMessageSize"`
	ReadWait        time.Duration `json:"readWait"`
	WriteWait       time.Duration `json:"writeWait"`
	PingWait        time.Duration `json:"pingWait"`
	PongWait        time.Duration `json:"pongWait"`
}

var DefaultConfig = &WsConfig{
	ReadBufferSize:  ReadBufferSize,
	WriteBufferSize: WriteBufferSize,
	MaxMessageSize:  MaxMessageSize,
	ReadWait:        ReadWait,
	WriteWait:       WriteWait,
	PingWait:        PingWait,
	PongWait:        PongWait,
}
