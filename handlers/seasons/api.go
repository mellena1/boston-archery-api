package seasons

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/model"
)

type SeasonDB interface {
	PutSeason(ctx context.Context, season model.Season) (*model.Season, error)
	GetAllSeasons(ctx context.Context) ([]model.Season, error)
	GetSeason(ctx context.Context, uuid uuid.UUID) (*model.Season, error)
}

type API struct {
	logger *slog.Logger
	db     SeasonDB
}

func NewAPI(logger *slog.Logger, db SeasonDB) *API {
	return &API{
		logger: logger,
		db:     db,
	}
}
