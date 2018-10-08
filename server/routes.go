package server

const (
	GetExchangeSymbols = "/v1/symbols/:exchange"
	GetSchemas         = "/v1/schemas"
	CreateSchema       = GetSchemas
	GetStatistics      = "/v1/statistics"
	GetLiveFeed        = "/v1/feeds/live"
	GetFeeds           = "/v1/feeds"
	CreateFeed         = GetFeeds
	GetClientFeeds     = "/v1/feeds/client/:sessionId"
	SubscribeToFeed    = "/v1/feeds/:feedName/subscribe"
)
