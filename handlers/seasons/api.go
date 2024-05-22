package seasons

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db"
)

type SeasonDB interface {
	AddSeason(ctx context.Context, newSeason db.SeasonInput) (*db.Season, error)
	GetAllSeasons(ctx context.Context) ([]db.Season, error)
	GetSeasonByName(ctx context.Context, name string) (*db.Season, error)
	UpdateSeason(ctx context.Context, id uuid.UUID, season db.SeasonInput) (*db.Season, error)
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
