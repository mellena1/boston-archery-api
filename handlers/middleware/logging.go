package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// formatDuration formats a duration to one decimal point.
func formatDuration(d time.Duration) string {
	div := time.Duration(10)
	switch {
	case d > time.Second:
		d = d.Round(time.Second / div)
	case d > time.Millisecond:
		d = d.Round(time.Millisecond / div)
	case d > time.Microsecond:
		d = d.Round(time.Microsecond / div)
	case d > time.Nanosecond:
		d = d.Round(time.Nanosecond / div)
	}
	return d.String()
}

func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		clientIP := c.ClientIP()
		if clientIP == "" {
			clientIP = c.Request.RemoteAddr
		}

		logger.InfoContext(c, "",
			slog.Time("timestamp", time.Now()),
			slog.String("latency", formatDuration(time.Since(start))),
			slog.String("client_ip", clientIP),
			slog.String("method", c.Request.Method),
			slog.Int("status_code", c.Writer.Status()),
			slog.Int("body_size", c.Writer.Size()),
			slog.String("path", c.Request.URL.Path),
			slog.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}
