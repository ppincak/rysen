package scrape

type CallerFuncType uint32

const (
	OrdersEvent CallerFuncType = iota
	TradesEvent
	PriceEvent
)

type CallerFunc func(topic string, symbols []string)

type CallerEvent struct {
	Symbol    string
	EventType CallerFuncType
	Data      interface{}
	Timestamp int64
}
