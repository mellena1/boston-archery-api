package model

import (
	"time"

	"github.com/google/uuid"
)

type Season struct {
	ID        uuid.UUID
	Name      string
	StartDate time.Time
	EndDate   time.Time
	ByeWeeks  []time.Time
}
