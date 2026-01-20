package api

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// MiddlewareTrace 接口层链路追踪，设置span
func MiddlewareTrace(trace Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {

		newCtx, span := trace.AddServerTrace(c)
		if span != nil {
			defer span.End()
			// defer func() { trace.TelemetrySpanEnd(span, c.Err()) }()
			req := c.Request.WithContext(newCtx)
			c.Request = req

			c.Next()

			status := c.Writer.Status()
			if status/100 >= 4 {
				span.SetStatus(codes.Error, "REQUEST FAILED")
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
			if status > 0 {
				span.SetAttributes(semconv.HTTPStatusCode(status))
			}
		}
	}
}
