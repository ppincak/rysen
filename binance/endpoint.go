package binance

type endpoints struct {
	AggregateTrades  string
	Candlesticks     string
	ExchangeInfo     string
	HistoricalTrades string
	OrderBook        string
	OrderBookTicker  string
	Ticker24         string
	TickerPrice      string
	Trades           string
}

func newV1Endpoints() endpoints {
	return endpoints{
		AggregateTrades:  "/api/v1/aggTrades",
		Candlesticks:     "/api/v1/klines",
		ExchangeInfo:     "/api/v1/exchangeInfo",
		HistoricalTrades: "/api/v1/historicalTrades",
		OrderBook:        "/api/v1/depth",
		OrderBookTicker:  "/api/v3/ticker/bookTicker",
		Ticker24:         "/api/v1/ticker/24hr",
		TickerPrice:      "/api/v3/ticker/price",
		Trades:           "/api/v1/trades",
	}
}

var Endpoints = newV1Endpoints()
