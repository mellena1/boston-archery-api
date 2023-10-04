package models

import (
	"time"

	"github.com/google/uuid"
)

type Season struct {
	UUID      uuid.UUID
	Name      string
	StartDate time.Time
	EndDate   time.Time
	ByeWeeks  []time.Time
}
