package db

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db/tablekeys"
	"github.com/mellena1/boston-archery-api/model"
)

const playerEntityType = "Player"

type playerDynamoItem struct {
	PK         string
	SK         string
	EntityType string
	FirstName  string
	LastName   string
}

func (p playerDynamoItem) toPlayer() model.Player {
	return model.Player{
		ID:        uuid.MustParse(strings.Split(p.PK, "#")[1]),
		FirstName: p.FirstName,
		LastName:  p.LastName,
	}
}

type PlayerInput struct {
	FirstName string
	LastName  string
}

func (p PlayerInput) toDynamoItem(id uuid.UUID) playerDynamoItem {
	key := playerPK(id.String())
	return playerDynamoItem{
		PK:         key,
		SK:         key,
		EntityType: playerEntityType,
		FirstName:  p.FirstName,
		LastName:   p.LastName,
	}
}

func playerPK(key string) string {
	return tablekeys.PLAYER_KEY_PREFIX + key
}

func (db *DB) AddPlayer(ctx context.Context, newPlayer PlayerInput) (*model.Player, error) {
	dynamoItem := newPlayer.toDynamoItem(uuid.New())
	err := db.putItem(ctx, dynamoItem)
	if err != nil {
		return nil, err
	}

	player := dynamoItem.toPlayer()
	return &player, nil
}

func (db *DB) UpdatePlayer(ctx context.Context, id uuid.UUID, player PlayerInput) (*model.Player, error) {
	dynamoItem := player.toDynamoItem(uuid.New())
	err := db.putItem(ctx, dynamoItem)
	if err != nil {
		return nil, err
	}

	updatedPlayer := dynamoItem.toPlayer()
	return &updatedPlayer, nil
}

func (db *DB) GetPlayer(ctx context.Context, id uuid.UUID) (*model.Player, error) {
	var playerItem playerDynamoItem

	key := playerPK(id.String())
	err := db.getItem(ctx, key, key, &playerItem)
	if err != nil {
		return nil, err
	}

	player := playerItem.toPlayer()
	return &player, nil
}

func (db *DB) GetAllPlayers(ctx context.Context) ([]model.Player, error) {
	var playerItems []playerDynamoItem
	err := db.getAllOfEntity(ctx, playerEntityType, &playerItems)
	if err != nil {
		return nil, err
	}

	players := make([]model.Player, len(playerItems))
	for i, item := range playerItems {
		players[i] = item.toPlayer()
	}

	return players, nil
}
