package models

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID            uuid.UUID
	DateTime      time.Time
	SeasonID      uuid.UUID
	HomeTeamID    uuid.UUID
	AwayTeamID    uuid.UUID
	HomeTeamScore int
	AwayTeamScore int
	YoutubeURL    string
}
