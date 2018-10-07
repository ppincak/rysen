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

// Initialize Router
func (router *Router) Init(engine *gin.Engine) {
	engine.GET(RoutesV1.statistics, router.getStatistics)
	engine.GET(RoutesV1.symbols, router.getSymbols)
	engine.GET(RoutesV1.live, router.getLive)
	engine.GET(RoutesV1.feeds, router.getFeeds)
	engine.GET(RoutesV1.clientFeeds, router.getClientFeeds)
	engine.POST(RoutesV1.feeds, router.createFeed)
	engine.POST(RoutesV1.subscribeToFeed, router.subscribeToFeed)
	engine.GET(RoutesV1.schema, router.getSchemas)
	engine.POST(RoutesV1.schema, router.createSchema)
}

func (router *Router) getStatistics(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.Monitor.Statistics())
}

func (router *Router) getSymbols(context *gin.Context) {
	exchangeName := context.Param("exchange")
	if exchange, ok := router.app.Exchanges[exchangeName]; !ok {
		errors.BadRequest(context, "Invalid Exchange", "invalid.exchange")
	} else {
		context.JSON(http.StatusOK, exchange.Symbols())
	}
}

func (router *Router) getLive(context *gin.Context) {
	router.app.WsHandler.ServeWebSocket(context.Writer, context.Request)
}

func (router *Router) getFeeds(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.FeedService.ListFeeds())
}

func (router *Router) getClientFeeds(context *gin.Context) {
	sessionId := context.Param("sessionId")
	if sessionId == "" {
		errors.BadRequest(context, "Invalid session id", "invalid.sessionId")
	} else {
		context.JSON(http.StatusOK, router.app.FeedService.ListClientFeeds(sessionId))
	}
}

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

func (router *Router) getPublishers(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.PublisherService.ListPublishers())
}

func (router *Router) createPublisher(context *gin.Context) {

}

func (router *Router) subscribeToFeed(context *gin.Context) {
	clientId := context.DefaultQuery("clientId", "")
	feed := context.Param("feed")
	if clientId == "" {
		errors.BadRequest(context, "Missing clientId param", "missing.clienId")
		return
	}
	if feed == "" {
		errors.BadRequest(context, "Missing feed param", "missing.feed")
		return
	}
	client := router.app.WsHandler.GetClient(clientId)
	if client == nil {
		errors.BadRequest(context, "Invalid clientId", "invalid.clientId")
		return
	}
	if router.app.FeedService.SubscribeTo(feed, client) != nil {
		errors.BadRequest(context, "Invalid feed", "invalid.feed")
		return
	}
	context.Status(http.StatusOK)
}

func (router *Router) getSchemas(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.SchemaService.ListSchemas())
}

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
