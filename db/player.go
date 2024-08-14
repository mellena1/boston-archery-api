package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db/tablekeys"
	"github.com/mellena1/boston-archery-api/model"
	"github.com/mellena1/boston-archery-api/slices"
)

type playerDynamoItem struct {
	PK        string
	SK        string
	GSI1PK    string
	GSI1SK    string
	ID        string
	FirstName string
	LastName  string
}

func (p playerDynamoItem) toPlayer() model.Player {
	return model.Player{
		ID:        uuid.MustParse(p.ID),
		FirstName: p.FirstName,
		LastName:  p.LastName,
	}
}

func playerToDynamoItem(player model.Player) playerDynamoItem {
	key := playerPK(player.ID.String())
	return playerDynamoItem{
		PK:        key,
		SK:        key,
		GSI1PK:    tablekeys.PLAYER_GSI1_PK,
		GSI1SK:    key,
		ID:        player.ID.String(),
		FirstName: player.FirstName,
		LastName:  player.LastName,
	}
}

func playerPK(key string) string {
	return tablekeys.PLAYER_KEY_PREFIX + key
}

func (db *DB) AddPlayer(ctx context.Context, newPlayer model.Player) (*model.Player, error) {
	dynamoItem := playerToDynamoItem(newPlayer)

	putCond := expression.AttributeNotExists(expression.Name(tablekeys.PK))
	putExpr, err := expression.NewBuilder().WithCondition(putCond).Build()
	if err != nil {
		panic("addplayer condition is broken: " + err.Error())
	}

	err = db.putItem(ctx, dynamoItem, withPutItemConditionExpression(putExpr))
	if err != nil {
		var condCheckErr *types.ConditionalCheckFailedException
		switch {
		case errors.As(err, &condCheckErr):
			return nil, ErrItemAlreadyExists
		}

		return nil, err
	}

	return &newPlayer, nil
}

type UpdatePlayerInput struct {
	FirstName *string
	LastName  *string
}

func (db *DB) UpdatePlayer(ctx context.Context, id uuid.UUID, updates UpdatePlayerInput) (*model.Player, error) {
	existsCond := expression.AttributeExists(expression.Name(tablekeys.PK))
	var update expression.UpdateBuilder

	if updates.FirstName != nil {
		update = update.Set(expression.Name("FirstName"), expression.Value(*updates.FirstName))
	}
	if updates.LastName != nil {
		update = update.Set(expression.Name("LastName"), expression.Value(*updates.LastName))
	}

	expr, err := expression.NewBuilder().
		WithCondition(existsCond).
		WithUpdate(update).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression for updatePlayer: %w", err)
	}

	key := playerPK(id.String())
	result, err := db.updateItem(ctx, key, key, withUpdateExpression(expr), withUpdateReturnValues(types.ReturnValueAllNew))
	if err != nil {
		var condCheckErr *types.ConditionalCheckFailedException
		switch {
		case errors.As(err, &condCheckErr):
			return nil, ErrItemNotFound
		}

		return nil, err
	}

	var dynamoItem playerDynamoItem
	err = attributevalue.UnmarshalMap(result.Attributes, &dynamoItem)
	if err != nil {
		return nil, fmt.Errorf("error reading result of player upate: %w", err)
	}

	playerResult := dynamoItem.toPlayer()
	return &playerResult, nil
}

func (db *DB) GetPlayer(ctx context.Context, id uuid.UUID) (*model.Player, error) {
	var playerItem playerDynamoItem

	key := playerPK(id.String())
	err := db.getItem(ctx, key, key, &playerItem)
	if err != nil {
		return nil, fmt.Errorf("failed to get player %q: %w", id, err)
	}

	player := playerItem.toPlayer()
	return &player, nil
}

func (db *DB) GetAllPlayers(ctx context.Context) ([]model.Player, error) {
	keyCond := expression.Key(tablekeys.GSI1PK).Equal(expression.Value(tablekeys.PLAYER_GSI1_PK))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		panic("get all players key condition broken: " + err.Error())
	}

	var playerItems []playerDynamoItem
	err = db.getManyOfEntity(ctx, &playerItems, withQueryKeyConditionExpression(expr), withQueryIndex(tablekeys.GSI1_INDEX))
	if err != nil {
		return nil, fmt.Errorf("failed to get all players: %w", err)
	}

	players := slices.Map(playerItems, func(item playerDynamoItem) model.Player {
		return item.toPlayer()
	})

	return players, nil
}
