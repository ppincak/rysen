package binance

import (
	"rysen/binance/model"
	"rysen/crypto"
)

func NewSymbols(exchangeInfo *model.ExchangeInfo) *crypto.Symbols {
	return &crypto.Symbols{
		Assets:  ToAssetsMap(exchangeInfo),
		Symbols: ToSymbolSlice(exchangeInfo),
	}
}

func ToAssetsMap(exchangeInfo *model.ExchangeInfo) map[string][]string {
	result := make(map[string][]string)
	for _, symbol := range exchangeInfo.Symbols {
		var symbols []string
		if _, ok := result[symbol.BaseAsset]; ok {
			symbols = result[symbol.BaseAsset]
		} else {
			symbols = []string{}
		}
		symbols = append(symbols, symbol.Symbol)

		result[symbol.BaseAsset] = symbols
	}
	return result
}

func ToSymbolSlice(exchangeInfo *model.ExchangeInfo) []string {
	symbols := make(map[string]bool)
	for i := 0; i < len(exchangeInfo.Symbols); i++ {
		symbols[exchangeInfo.Symbols[i].Symbol] = true
	}
	result := make([]string, len(symbols))
	i := 0
	for symbol, _ := range symbols {
		result[i] = symbol
		i++
	}
	return result
}
