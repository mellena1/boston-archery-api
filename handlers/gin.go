package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/logging"
)

func NewGin(logger *slog.Logger) *gin.Engine {
	engine := gin.New()

	engine.Use(logging.GinMiddlewareLogger(logger), gin.Recovery())

	return engine
}
