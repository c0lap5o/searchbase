package middlewares

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type sloggerMiddleware struct {
	logger *slog.Logger
}

// GinSlogger returns privacy-preserving structured request logging middleware.
func GinSlogger(logger *slog.Logger) gin.HandlerFunc {
	m := &sloggerMiddleware{
		logger: logger.With(slog.String("component", "gin")),
	}
	return m.handle
}

func (m *sloggerMiddleware) handle(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path

	defer func() {
		latency := time.Since(start)
		status := c.Writer.Status()

		// Pre-allocate slice with capacity 5 to prevent resizing
		attrs := make([]slog.Attr, 0, 7)
		attrs = append(attrs,
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("latency", latency),
		)

		spanContext := trace.SpanFromContext(c.Request.Context()).SpanContext()
		if spanContext.HasTraceID() {
			attrs = append(attrs, slog.String("trace_id", spanContext.TraceID().String()))
		}

		if spanContext.HasSpanID() {
			attrs = append(attrs, slog.String("span_id", spanContext.SpanID().String()))
		}

		// Trust the edge proxy (Cloudflare/Nginx) for region data, never log the raw IP
		if country := c.GetHeader("CF-IPCountry"); country != "" {
			attrs = append(attrs, slog.String("country", country))
		}

		if len(c.Errors) > 0 {
			if errs := c.Errors.ByType(gin.ErrorTypePrivate).String(); errs != "" {
				attrs = append(attrs, slog.String("errors", errs))
			}
		}

		level := slog.LevelInfo
		msg := "request handled"

		if status >= 500 {
			level = slog.LevelError
			msg = "server error"
		}

		m.logger.LogAttrs(c.Request.Context(), level, msg, attrs...)
	}()

	c.Next()
}
