package binance

import (
	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/store"
)

type BinanceStore struct {
	*store.Store
}

func NewBinanceStore() *BinanceStore {
	return &BinanceStore{
		Store: store.NewStore(),
	}
}

func (store *BinanceStore) Initialize(client *BinanceClient) error {
	info, err := client.ExchangeInfo()
	if err != nil {
		return api.NewError("Failed to initialize BinanceStore")
	}
	store.Symbols = NewSymbols(info)
	return nil
}
