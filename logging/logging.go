package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/handlers/middleware"
)

func NewLogger(level slog.Level) *slog.Logger {
	return slog.New(&ginHandler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})})
}

type ginHandler struct {
	slog.Handler
}

func (h *ginHandler) Handle(ctx context.Context, record slog.Record) error {
	ginCtx, ok := ctx.(*gin.Context)
	if ok {
		requestId := requestid.Get(ginCtx)
		record.AddAttrs(slog.String("request_id", requestId))

		claims := middleware.GetJWTClaimsCtx(ginCtx)
		if claims != nil {
			record.AddAttrs(
				slog.Group("auth",
					slog.String("user_id", claims.UserID),
					slog.String("username", claims.Username),
					slog.Bool("is_admin", claims.IsAdmin),
					slog.String("nickname", claims.Nickname),
				),
			)
		}
	}

	return h.Handler.Handle(ctx, record)
}

func (h *ginHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ginHandler{
		Handler: h.Handler.WithAttrs(attrs),
	}
}

func (h *ginHandler) WithGroup(name string) slog.Handler {
	return &ginHandler{
		Handler: h.Handler.WithGroup(name),
	}
}
