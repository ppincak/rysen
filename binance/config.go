package binance

type Config struct {
	// Exchange url
	url string
}

// Create new config
func NewConfig(url string) *Config {
	return &Config{
		url: url,
	}
}
