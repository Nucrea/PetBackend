package middleware

import (
	"backend/src/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewRequestLogMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetHeader("X-Request-Id")
		if requestId == "" {
			requestId = uuid.New().String()
		}

		path := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			path = path + "?" + c.Request.URL.RawQuery
		}

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		method := c.Request.Method
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		e := logger.Log()
		e.Str("id", requestId)
		e.Str("ip", clientIP)
		e.Msgf("[REQUEST] %s %s %d %v", method, path, statusCode, latency)
	}
}
