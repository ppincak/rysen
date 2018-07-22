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
}

func (router *Router) getLive(context *gin.Context) {
	router.app.WsHandler.ServeWebSocket(context.Writer, context.Request)
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
