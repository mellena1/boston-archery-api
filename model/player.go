package model

import "github.com/google/uuid"

type Player struct {
	ID              uuid.UUID
	FirstName       string
	LastName        string
	DiscordUserID   string
	DiscordUserName string
}
