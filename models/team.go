package models

import "github.com/google/uuid"

type Team struct {
	ID         uuid.UUID
	SeasonID   uuid.UUID
	Name       string
	TeamColors []string
	Captain    string
}
