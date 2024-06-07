package seasons

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db"
	"github.com/mellena1/boston-archery-api/model"
)

type SeasonDB interface {
	AddSeason(ctx context.Context, newSeason db.SeasonInput) (*model.Season, error)
	GetAllSeasons(ctx context.Context) ([]model.Season, error)
	GetSeasonByName(ctx context.Context, name string) (*model.Season, error)
	UpdateSeason(ctx context.Context, id uuid.UUID, season db.SeasonInput) (*model.Season, error)
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
