package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func NewTracingMiddleware(tracer trace.Tracer) gin.HandlerFunc {
	prop := otel.GetTextMapPropagator()

	return func(c *gin.Context) {
		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()

		ctx := prop.Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))

		ctx, span := tracer.Start(ctx, c.Request.URL.Path)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
