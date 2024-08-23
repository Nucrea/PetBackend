package middleware

import (
	log "backend/src/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewRequestLogMiddleware(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetHeader("X-Request-Id")
		if requestId == "" {
			requestId = uuid.New().String()
		}

		log.SetCtxRequestId(c, requestId)

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

		ctxLogger := logger.WithContext(c)

		e := ctxLogger.Log()
		e.Str("ip", clientIP)
		e.Msgf("Request %s %s %d %v", method, path, statusCode, latency)
	}
}
