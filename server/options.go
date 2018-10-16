package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create new handler for HTTP OPTIONS requests
func NewOptionsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
			c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.AbortWithStatus(http.StatusOK)
		}
	}
}
