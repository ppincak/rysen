package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.GET(RoutesV1.symbols, GetSymbols)
	router.GET(RoutesV1.config, GetConfig)
}

func GetSymbols(context *gin.Context) {
	context.JSON(http.StatusOK, nil)
}

func GetConfig(context *gin.Context) {
	context.JSON(http.StatusOK, nil)
}

func Mertrics(context *gin.Context) {
	context.JSON(http.StatusOK, nil)
}
