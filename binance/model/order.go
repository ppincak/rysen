package model

type OrderType uint8

const (
	LIMIT OrderType = iota
	MARKET
	STOP_LOSS
	STOP_LOSS_LIMIT
	TAKE_PROFIT
	TAKE_PROFIT_LIMIT
	LIMIT_MAKER
)

var OrderTypes = map[OrderType]string{}

type OrderStatus uint8

const (
	NEW OrderStatus = iota
	PARTIALLY_FILLED
	FILLED
	CANCELED
	PENDING_CANCEL
	REJECTED
	EXPIRED
)

var OrderStatuses = map[OrderStatus]string{}

type Order struct {
	Symbol   string  `json:"symbol"`
	Side     string  `json:"side"`
	Type     string  `json:"type"`
	Quantity int64   `json:"quantity"`
	Price    float64 `json:"price"`
}
