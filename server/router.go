package server

import (
	"net/http"

	"github.com/ppincak/rysen/pkg/ws"
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
	engine.GET(GetExchanges, router.getExchanges)
	engine.GET(GetExchangeSymbols, router.getExchangeSymbols)
	engine.GET(GetLiveFeed, router.getLiveFeed)
	engine.GET(GetFeeds, router.getFeeds)
	engine.GET(GetClientFeeds, router.getClientFeeds)
	engine.POST(CreateFeed, router.createFeed)
	engine.DELETE(DeleteFeed, router.deleteFeed)
	engine.POST(SubscribeToFeed, router.subscribeToFeed)
	engine.DELETE(UnsubscribeFromFeed, router.subscribeToFeed)
	engine.GET(GetSchemas, router.getSchemas)
	engine.GET(GetSchema, router.getSchema)
	engine.POST(CreateSchema, router.createSchema)
	engine.PUT(UpdateSchema, router.updateSchema)
	engine.DELETE(DeleteSchema, router.deleteSchema)
}

// Get system statistics/metrics
func (router *Router) getStatistics(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.Monitor.Statistics())
}

// Get exchanges
func (router *Router) getExchanges(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.Exchanges.List())
}

// Get exchange symbols
func (router *Router) getExchangeSymbols(context *gin.Context) {
	exchangeName := context.Param("exchange")
	if exchange, ok := router.app.Exchanges[exchangeName]; !ok {
		errors.BadRequest(context, "Invalid Exchange", "invalid.exchange")
	} else {
		context.JSON(http.StatusOK, exchange.Symbols())
	}
}

// Get live feed served as websocket
func (router *Router) getLiveFeed(context *gin.Context) {
	router.app.WsHandler.ServeWebSocket(context.Writer, context.Request)
}

// Get list of feeds
func (router *Router) getFeeds(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.FeedService.ListFeeds())
}

// Get list of client feeds
func (router *Router) getClientFeeds(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.FeedService.ListClientFeeds(context.Param("sessionId")))
}

// Create a feed
func (router *Router) createFeed(context *gin.Context) {
	var model *feed.Model
	if err := context.ShouldBindJSON(&model); err != nil {
		errors.BadRequest(context, "Deserialization failed", "deserialization.failed")
		return
	}

	_, err := router.app.FeedService.Create(model)
	if err != nil {
		errors.ErrorBadRequest(context, err)
		return
	}
	err = router.app.FeedPersistence.SaveFeed(model)
	if err != nil {
		errors.ErrorBadRequest(context, err)
		return
	}
	context.Status(http.StatusCreated)
}

// Delete a feed
func (router *Router) deleteFeed(context *gin.Context) {
	feedName := context.DefaultQuery("feed", "")
	if feedName == "" {
		errors.BadRequest(context, "Invalid feed name", "invalid.feed")
	}
	if err := router.app.FeedService.Delete(feedName); err != nil {
		errors.ErrorBadRequest(context, err)
		return
	}
	if err := router.app.FeedPersistence.Delete(feedName); err != nil {
		errors.ErrorBadRequest(context, err)
		return
	}
	context.Status(http.StatusOK)
}

// Get list of publishers
func (router *Router) getPublishers(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.PublisherService.ListPublishers())
}

// Subscribe to a feed
func (router *Router) subscribeToFeed(context *gin.Context) {
	router.handleSubscription(context, router.app.FeedService.SubscribeTo)
}

// Unsubscribe from a feed
func (router *Router) unsubscribeFrom(context *gin.Context) {
	router.handleSubscription(context, router.app.FeedService.UnsubscribeFrom)
}

// Get list of schemas
func (router *Router) getSchemas(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.SchemaService.ListSchemas())
}

// Get single schema by name
func (router *Router) getSchema(context *gin.Context) {
	schemaName := context.Param("schemaName")
	if schema, err := router.app.SchemaService.GetSchema(schemaName); err != nil {
		errors.ErrorBadRequest(context, err)
	} else {
		context.JSON(http.StatusOK, schema)
	}
}

// Create a schema
func (router *Router) createSchema(context *gin.Context) {
	router.handleSchema(context, router.app.SchemaService.CreateSchema)
}

// Update a schema
func (router *Router) updateSchema(context *gin.Context) {
	router.handleSchema(context, router.app.SchemaService.UpdateSchema)
}

// Delete a schema
func (router *Router) deleteSchema(context *gin.Context) {
	schemaName := context.Param("schemaName")
	if err := router.app.SchemaService.DeleteSchema(schemaName); err != nil {
		errors.ErrorBadRequest(context, err)
		return
	}
	if err := router.app.SchemaPersistence.Delete(schemaName); err != nil {
		errors.ErrorBadRequest(context, err)
		return
	}
	context.Status(http.StatusOK)
}

func (router *Router) handleSubscription(context *gin.Context, handler func(name string, client *ws.Client) error) {
	sessionId := context.DefaultQuery("sessionId", "")
	if sessionId == "" {
		errors.BadRequest(context, "Missing sessionId param", "missing.sessionId")
		return
	}
	feed := context.DefaultQuery("feed", "")
	if feed == "" {
		errors.BadRequest(context, "Invalid feed", "invalid.feed")
		return
	}
	client := router.app.WsHandler.GetClient(sessionId)
	if client == nil {
		errors.BadRequest(context, "Invalid sessionId", "invalid.sessionId")
		return
	}
	if handler(feed, client) != nil {
		errors.BadRequest(context, "Invalid feed", "invalid.feed")
		return
	}
	context.Status(http.StatusOK)
}

// Schema handler
func (router *Router) handleSchema(context *gin.Context, handlerFunc func(*schema.Model) (*schema.ExchangeSchema, error)) {
	var schema *schema.Model
	if err := context.ShouldBindJSON(&schema); err != nil {
		errors.BadRequest(context, "Deserialization failed", "deserialization.failed")
		return
	}

	_, err := handlerFunc(schema)
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
