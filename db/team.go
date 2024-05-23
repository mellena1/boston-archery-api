package db

import "github.com/google/uuid"

type Team struct {
	ID         uuid.UUID
	Name       string
	TeamColors []string
	Captain    uuid.UUID
}
