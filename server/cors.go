package server

import "github.com/gin-gonic/gin"

// Create new CORS handler function
func NewCorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}
}
