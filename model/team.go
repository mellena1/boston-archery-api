package model

import "github.com/google/uuid"

type Team struct {
	ID         uuid.UUID
	Name       string
	TeamColors []string
}
