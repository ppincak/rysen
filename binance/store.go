package binance

import (
	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/store"
)

type Store struct {
	*store.Store
}

func NewStore() *Store {
	return &Store{
		Store: store.NewStore(),
	}
}

func (store *Store) Initialize(client *Client) error {
	info, err := client.ExchangeInfo()
	if err != nil {
		return api.NewError("Failed to initialize BinanceStore")
	}
	store.Symbols = NewSymbols(info)
	return nil
}
