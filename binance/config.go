package binance

import "github.com/ppincak/rysen/api"

type Config struct {
	// Exchange url
	url string
	// Secret
	secret *api.Secret
}

// Create new config
func NewConfig(url string, secret *api.Secret) *Config {
	return &Config{
		url:    url,
		secret: secret,
	}
}
