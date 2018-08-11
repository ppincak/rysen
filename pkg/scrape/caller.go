package scrape

// Caller function type
type CallerFuncType uint32

// All possible caller events
const (
	OrdersEvent CallerFuncType = iota
	TradesEvent
	PriceEvent
)

// Caller
type Caller interface {
	ScrapeOrders(topic string, symbols []string)
	ScrapePrices(topic string, symbols []string)
	ScrapeTrades(topic string, symbols []string)
}

// Single caller function
type CallerFunc func(topic string, symbols []string)

// Event emitted by caller
type CallerEvent struct {
	Symbol    string
	EventType CallerFuncType
	Data      interface{}
	Timestamp int64
}
