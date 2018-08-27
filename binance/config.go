package binance

import "github.com/ppincak/rysen/api"

type Config struct {
	// Exchange url
	Url string
	// Secret
	Secret *api.Secret
}
