package scrape

// Caller function type
type CallerFuncType uint32

// All possible caller events
const (
	OrdersEvent CallerFuncType = iota
	TradesEvent
	PriceEvent
)

// Mapping from caller func type to its string representation
var FuncToString = map[CallerFuncType]string{
	0: "scrapeOrders",
	1: "scrapePrices",
	2: "scrapeTrades",
}

// Mapping from string representation to caller func type
var StringToFunc = map[string]CallerFuncType{
	"scrapeOrders": 0,
	"scrapePrices": 1,
	"scrapeTrades": 2,
}

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
