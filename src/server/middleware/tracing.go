package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NewTracingMiddleware(tracer trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, span := tracer.Start(c.Request.Context(), c.Request.URL.Path)
		defer span.End()

		span.SetAttributes(attribute.String("requestId", c.ClientIP()))

		c.Next()
	}
}
