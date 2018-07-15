package model

type ExchangeInfo struct {
	RateLimits []RateLimit `json:"rateLimits"`
	Symbols    []Symbol    `json:"symbols"`
}

type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	Limit         uint32 `json:"limit"`
}

type Symbol struct {
	Symbol    string `json:"symbol"`
	BaseAsset string `json:"baseAsset"`
}
