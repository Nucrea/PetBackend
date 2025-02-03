package middleware

import (
	"backend/internal/integrations"
	log "backend/pkg/logger"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

func NewRequestLogMiddleware(logger log.Logger, tracer trace.Tracer, prometheus *integrations.Prometheus) gin.HandlerFunc {
	return func(c *gin.Context) {
		prometheus.RequestInc()
		defer prometheus.RequestDec()

		requestId := c.GetHeader("X-Request-Id")
		if requestId == "" {
			requestId = uuid.New().String()
		}
		c.Header("X-Request-Id", requestId)
		c.Header("Access-Control-Allow-Origin", "*")

		log.SetCtxRequestId(c, requestId)

		path := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			path = path + "?" + c.Request.URL.RawQuery
		}

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		prometheus.AddRequestTime(float64(latency.Microseconds()))

		method := c.Request.Method
		statusCode := c.Writer.Status()

		ctxLogger := logger.WithContext(c)

		msg := fmt.Sprintf("Request %s %s %d %v", method, path, statusCode, latency)

		if statusCode >= 200 && statusCode < 400 {
			// ctxLogger.Log().Msg(msg)
			return
		}

		if statusCode >= 400 && statusCode < 500 {
			prometheus.Add4xxError()
			ctxLogger.Warning().Msg(msg)
			return
		}

		prometheus.Add5xxError()
		ctxLogger.Error().Msg(msg)
	}
}
