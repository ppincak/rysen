package ws

type Packet struct {
	Id      string      `json:"id"`
	Topic   string      `json:"topic"`
	Message interface{} `json:"message"`
}
