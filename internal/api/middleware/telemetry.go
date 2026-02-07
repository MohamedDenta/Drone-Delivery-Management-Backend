package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// TelemetryMiddleware wraps the otelgin middleware
// This automatically starts spans and propagates context
func TelemetryMiddleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName)
}
