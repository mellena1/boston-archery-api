package db

import "github.com/google/uuid"

type Player struct {
	ID   uuid.UUID
	Name string
}
