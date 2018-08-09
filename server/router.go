package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	app *App
}

func NewRouter(app *App) *Router {
	return &Router{
		app: app,
	}
}

func (router *Router) Init(engine *gin.Engine) {
	engine.GET(RoutesV1.live, router.getLive)
	engine.GET(RoutesV1.symbols, router.getSymbols)
	engine.GET(RoutesV1.config, router.getConfig)
	engine.GET(RoutesV1.statistics, router.getStatistics)
	engine.GET(RoutesV1.feeds, router.getFeeds)
	engine.POST(RoutesV1.subscribeToFeed, router.postSubscribeToFeed)
}

func (router *Router) getLive(context *gin.Context) {
	router.app.WsHandler.ServeWebSocket(context.Writer, context.Request)
}

func (router *Router) postSubscribeToFeed(context *gin.Context) {
	clientId := context.DefaultQuery("clientId", "")
	feed := context.DefaultQuery("feed", "")
	if clientId == "" {
		context.JSON(http.StatusBadRequest, NewApiError("Missing clientId param", ""))
		return
	}
	if feed == "" {
		context.JSON(http.StatusBadRequest, NewApiError("Missing feedId param", ""))
		return
	}
	client := router.app.WsHandler.GetClient(clientId)
	if client == nil {
		context.Status(http.StatusBadRequest)
		return
	}
	router.app.FeedService.SubscribeTo(feed, client)

	context.Status(http.StatusOK)
}

func (router *Router) getFeeds(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.FeedService.GetList())
}

func (router *Router) getSymbols(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.Binance.Store.Symbols)
}

func (router *Router) getConfig(context *gin.Context) {
	context.JSON(http.StatusOK, nil)
}

func (router *Router) getStatistics(context *gin.Context) {
	context.JSON(http.StatusOK, router.app.Monitor.Statistics())
}
