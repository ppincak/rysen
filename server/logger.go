package server

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Create new logger handler function
func NewLogger(level log.Level) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		path := c.Request.URL.Path
		method := c.Request.Method
		statusCode := c.Writer.Status()
		latency := end.Sub(start)
		clientIP := c.ClientIP()

		entry := log.WithFields(log.Fields{
			"path":       path,
			"method":     method,
			"statusCode": statusCode,
			"latency":    latency,
			"clientIP":   clientIP,
		})

		levelSwitch(level, entry)("HTTP Request")
	}
}

// Generic log function
type LogFunc func(args ...interface{})

// Invoke correct log level function
func levelSwitch(level log.Level, entry *log.Entry) LogFunc {
	var logFunc LogFunc

	switch level {
	case log.DebugLevel:
		logFunc = entry.Debug
	case log.InfoLevel:
		logFunc = entry.Info
	}

	return logFunc
}
