package teams

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db"
	"github.com/mellena1/boston-archery-api/model"
)

type TeamDB interface {
	AddTeam(ctx context.Context, team model.Team) (*model.Team, error)
	UpdateTeam(ctx context.Context, id uuid.UUID, updates db.UpdateTeamInput) (*model.Team, error)
	GetAllTeams(ctx context.Context) ([]model.Team, error)
	GetTeam(ctx context.Context, id uuid.UUID) (*model.Team, error)
}

type API struct {
	logger *slog.Logger
	db     TeamDB
}

func NewAPI(logger *slog.Logger, db TeamDB) *API {
	return &API{
		logger: logger,
		db:     db,
	}
}
