package ws

import (
	"time"
)

const (
	ReadBufferSize  = 1024
	WriteBufferSize = 1024
	MaxMessageSize  = 1024 * 50
	ReadWait        = time.Hour
	WriteWait       = time.Second / 2
	PongWait        = time.Second * 60
	PingWait        = (time.Second * 60) - (PongWait / 10)
)

type Config struct {
	Port            int           `json:"port"`
	ReadBufferSize  int           `json:"readBufferSize"`
	WriteBufferSize int           `json:"writeBufferSize"`
	MaxMessageSize  int64         `json:"maxMessageSize"`
	ReadWait        time.Duration `json:"readWait"`
	WriteWait       time.Duration `json:"writeWait"`
	PingWait        time.Duration `json:"pingWait"`
	PongWait        time.Duration `json:"pongWait"`
}

var DefaultConfig = &Config{
	ReadBufferSize:  ReadBufferSize,
	WriteBufferSize: WriteBufferSize,
	MaxMessageSize:  MaxMessageSize,
	ReadWait:        ReadWait,
	WriteWait:       WriteWait,
	PingWait:        PingWait,
	PongWait:        PongWait,
}
