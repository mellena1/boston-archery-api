package db

import (
	"strings"

	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db/tablekeys"
	"github.com/mellena1/boston-archery-api/model"
)

const playerEntityType = "Player"

type playerDynamoItem struct {
	PK              string
	SK              string
	EntityType      string
	FirstName       string
	LastName        string
	DiscordUserID   string
	DiscordUserName string
}

func (p playerDynamoItem) toPlayer() model.Player {
	return model.Player{
		ID:              uuid.MustParse(strings.Split(p.PK, "#")[1]),
		FirstName:       p.FirstName,
		LastName:        p.LastName,
		DiscordUserID:   p.DiscordUserID,
		DiscordUserName: p.DiscordUserName,
	}
}

type PlayerInput struct {
	FirstName       string
	LastName        string
	DiscordUserID   string
	DiscordUserName string
}

func (p PlayerInput) toDynamoItem(id uuid.UUID) playerDynamoItem {
	key := playerPK(id.String())
	return playerDynamoItem{
		PK:              key,
		SK:              key,
		EntityType:      playerEntityType,
		FirstName:       p.FirstName,
		LastName:        p.LastName,
		DiscordUserID:   p.DiscordUserID,
		DiscordUserName: p.DiscordUserName,
	}
}

func playerPK(key string) string {
	return tablekeys.PLAYER_KEY_PREFIX + key
}
