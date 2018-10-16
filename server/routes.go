package server

const (
	GetExchangeSymbols = "/v1/symbols/:exchange"
	GetSchemas         = "/v1/schemas"
	GetSchema          = "/v1/schemas/:schemaName"
	CreateSchema       = GetSchemas
	UpdateSchema       = GetSchemas
	DeleteSchema       = GetSchema
	GetStatistics      = "/v1/statistics"
	GetLiveFeed        = "/v1/feeds/live"
	GetFeeds           = "/v1/feeds"
	CreateFeed         = GetFeeds
	GetClientFeeds     = "/v1/feeds/client/:sessionId"
	SubscribeToFeed    = "/v1/feeds/subscribe"
)
