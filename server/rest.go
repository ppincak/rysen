package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.GET(RoutesV1.symbols, getSymbols)
	router.GET(RoutesV1.config, getConfig)
	router.GET(RoutesV1.statistics, getStatistics)
	router.GET(RoutesV1.statistics, getStatistics)
}

func getWs(context *gin.Context) {

}

func getSymbols(context *gin.Context) {
	context.JSON(http.StatusOK, nil)
}

func getConfig(context *gin.Context) {
	context.JSON(http.StatusOK, nil)
}

func getStatistics(context *gin.Context) {
	context.JSON(http.StatusOK, nil)
}
