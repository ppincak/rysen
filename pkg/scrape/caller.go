package scrape

// Caller function type
type CallerFuncType uint32

// All possible caller events
const (
	OrdersEvent CallerFuncType = iota
	PriceEvent
	TradesEvent
)

// Mapping from caller func type to its string representation
var FuncToString = map[CallerFuncType]string{
	OrdersEvent: "scrapeOrders",
	PriceEvent:  "scrapePrices",
	TradesEvent: "scrapeTrades",
}

// Mapping from string representation to caller func type
var StringToFunc = map[string]CallerFuncType{
	"scrapeOrders": OrdersEvent,
	"scrapePrice":  PriceEvent,
	"scrapeTrades": TradesEvent,
}

// Caller
type Caller interface {
	ScrapeOrders(topic string, symbols []string)
	ScrapePrice(topic string, symbols []string)
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
