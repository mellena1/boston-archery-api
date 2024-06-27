package players

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db"
	"github.com/mellena1/boston-archery-api/model"
)

type PlayerDB interface {
	AddPlayer(ctx context.Context, newPLayer db.PlayerInput) (*model.Player, error)
	GetPlayer(ctx context.Context, id uuid.UUID) (*model.Player, error)
	GetAllPlayers(ctx context.Context) ([]model.Player, error)
	UpdatePlayer(ctx context.Context, id uuid.UUID, Player db.PlayerInput) (*model.Player, error)
}

type API struct {
	logger *slog.Logger
	db     PlayerDB
}

func NewAPI(logger *slog.Logger, db PlayerDB) *API {
	return &API{
		logger: logger,
		db:     db,
	}
}
