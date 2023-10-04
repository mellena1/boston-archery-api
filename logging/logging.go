package logging

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func NewLogger(level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}

func GinMiddlewareLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		logger.Info("",
			"timestamp", time.Now(),
			"latency", time.Since(start),
			"client_ip", c.ClientIP(),
			"method", c.Request.Method,
			"status_code", c.Writer.Status(),
			"body_size", c.Writer.Size(),
			"path", c.Request.URL.RequestURI(),
			"error", c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
	}
}
