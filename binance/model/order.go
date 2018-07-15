package model

type OrderType uint32

const (
	LIMIT OrderType = iota
	MARKET
	STOP_LOSS
	STOP_LOSS_LIMIT
	TAKE_PROFIT
	TAKE_PROFIT_LIMIT
	LIMIT_MAKER
)

type Order interface {
}

type RealOrder struct {
}

type VirtualOrder struct {
}
