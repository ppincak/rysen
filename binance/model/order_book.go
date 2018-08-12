package model

import "github.com/ppincak/rysen/pkg/converters"

type OrderBook struct {
	Asks []interface{} `json:"asks"`
	Bids []interface{} `json:"bids"`
}

func (book *OrderBook) ProcessAll() *OrderBookProcessed {
	processed := NewOrderBookProcessed(len(book.Asks), len(book.Bids))

	for i, order := range book.Asks {
		processed.Asks[i] = book.ProcessOne(order)
	}
	for i, order := range book.Bids {
		processed.Bids[i] = book.ProcessOne(order)
	}
	return processed
}

func (book *OrderBook) ProcessOne(order interface{}) *BookOrder {
	if asserted, ok := order.([]interface{}); ok {
		price, _ := converters.ToFloat64(asserted[0])
		volume, _ := converters.ToFloat64(asserted[1])

		return &BookOrder{
			Price:  price,
			Volume: volume,
		}
	}
	return nil
}

type OrderBookProcessed struct {
	Asks []*BookOrder `json:"asks"`
	Bids []*BookOrder `json:"bids"`
}

func NewOrderBookProcessed(askLen int, bidsLen int) *OrderBookProcessed {
	return &OrderBookProcessed{
		Asks: make([]*BookOrder, askLen),
		Bids: make([]*BookOrder, bidsLen),
	}
}

type BookOrder struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}
