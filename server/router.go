package server

import (
	"net/http"

	"github.com/ppincak/rysen/services/feed"

	"github.com/gin-gonic/gin"
	"github.com/ppincak/rysen/server/errors"
	"github.com/ppincak/rysen/services/schema"
)

type Router struct {
	app *App
}

// Create Router
func NewRouter(app *App) *Router {
	return &Router{
		app: app,
	}
}

// initialize http Router
func (router *Router) Init(engine *gin.Engine) {
	engine.GET(GetStatistics, router.getStatistics)
	engine.GET(GetExchangeSymbols, router.getExchangeSymbols)
	engine.GET(GetLiveFeed, router.getLiveFeed)
	engine.GET(GetFeeds, router.getFeeds)
	engine.GET(GetClientFeeds, router.getClientFeeds)
	engine.POST(CreateFeed, router.createFeed)
	engine.POST(SubscribeToFeed, router.subscribeToFeed)
	engine.GET(GetSchemas, router.getSchemas)
	engine.POST(CreateSchema, router.createSchema)
}

// get system statistics/metrics
func (router *Router) getStatistics(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.Monitor.Statistics())
}

//get exchange symbols
func (router *Router) getExchangeSymbols(context *gin.Context) {
	exchangeName := context.Param("exchange")
	if exchange, ok := router.app.Exchanges[exchangeName]; !ok {
		errors.BadRequest(context, "Invalid Exchange", "invalid.exchange")
	} else {
		context.JSON(http.StatusOK, exchange.Symbols())
	}
}

// get live feed served as websocket
func (router *Router) getLiveFeed(context *gin.Context) {
	router.app.WsHandler.ServeWebSocket(context.Writer, context.Request)
}

// get list of feeds
func (router *Router) getFeeds(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.FeedService.ListFeeds())
}

// get list of client feeds
func (router *Router) getClientFeeds(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.FeedService.ListClientFeeds(context.Param("sessionId")))
}

// create a feed
func (router *Router) createFeed(context *gin.Context) {
	var metadata *feed.Metadata
	if err := context.ShouldBindJSON(&metadata); err != nil {
		errors.BadRequest(context, "Deserialization failed", "deserialization.failed")
		return
	}

	_, err := router.app.FeedService.Create(metadata)
	if err != nil {
		errors.ErrorBadRequest(context, err)
		return
	}
	err = router.app.FeedPersistence.SaveFeed(metadata)
	if err != nil {
		errors.ErrorBadRequest(context, err)
		return
	}
	context.Status(http.StatusOK)
}

// get list of publishers
func (router *Router) getPublishers(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.PublisherService.ListPublishers())
}

// subscribe to a feed
func (router *Router) subscribeToFeed(context *gin.Context) {
	clientId := context.DefaultQuery("clientId", "")
	if clientId == "" {
		errors.BadRequest(context, "Missing clientId param", "missing.clienId")
		return
	}
	client := router.app.WsHandler.GetClient(clientId)
	if client == nil {
		errors.BadRequest(context, "Invalid clientId", "invalid.clientId")
		return
	}
	feed := context.Param("feed")
	if router.app.FeedService.SubscribeTo(feed, client) != nil {
		errors.BadRequest(context, "Invalid feed", "invalid.feed")
		return
	}
	context.Status(http.StatusOK)
}

// get list of schemas
func (router *Router) getSchemas(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.SchemaService.ListSchemas())
}

// create a schema
func (router *Router) createSchema(context *gin.Context) {
	var schema *schema.ExchangeSchemaMetadata
	if err := context.ShouldBindJSON(&schema); err != nil {
		errors.BadRequest(context, "Deserialization failed", "deserialization.failed")
		return
	}

	_, err := router.app.SchemaService.Create(schema)
	if err != nil {
		errors.ErrorBadRequest(context, err)
		return
	}
	err = router.app.SchemaPersistence.SaveSchema(schema)
	if err != nil {
		errors.ErrorBadRequest(context, err)
		return
	}
	context.Status(http.StatusOK)
}
