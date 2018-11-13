package server

const (
	GetExchanges        = "/v1/exchanges"
	GetExchangeSymbols  = "/v1/exchanges/:exchange/symbols"
	GetSchemas          = "/v1/schemas"
	GetSchema           = "/v1/schemas/:schemaName"
	CreateSchema        = GetSchemas
	UpdateSchema        = GetSchemas
	DeleteSchema        = GetSchema
	GetStatistics       = "/v1/statistics"
	GetLiveFeed         = "/v1/feeds/live"
	GetFeeds            = "/v1/feeds"
	CreateFeed          = GetFeeds
	DeleteFeed          = GetFeeds
	GetClientFeeds      = "/v1/feeds/client/:sessionId"
	SubscribeToFeed     = "/v1/feeds/subscription"
	UnsubscribeFromFeed = SubscribeToFeed
)
