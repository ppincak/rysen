package ws

import "github.com/google/uuid"

type PacketType int

const (
	Connect PacketType = iota
	Disconnect
	Event
	Error
)

type Packet struct {
	Id         string      `json:"id"`
	PacketType PacketType  `json:"type"`
	Content    interface{} `json:"content,omitempty"`
}

func NewPacket(packetType PacketType, content interface{}) *Packet {
	return &Packet{
		Id:         uuid.New().String(),
		PacketType: packetType,
		Content:    content,
	}
}

type EventMessage struct {
	Event   string      `json:"event"`
	Message interface{} `json:"message,omitempty"`
}
