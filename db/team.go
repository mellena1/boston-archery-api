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

type teamDynamoItem struct {
	PK         string
	SK         string
	GSI1PK     string
	GSI1SK     string
	ID         string
	Name       string
	TeamColors []string
}

func (t teamDynamoItem) toTeam() model.Team {
	return model.Team{
		ID:         uuid.MustParse(t.ID),
		Name:       t.Name,
		TeamColors: t.TeamColors,
	}
}

func teamToDynamoItem(team model.Team) teamDynamoItem {
	key := teamPK(team.ID.String())
	return teamDynamoItem{
		PK:         key,
		SK:         key,
		GSI1PK:     tablekeys.TEAM_GSI1_PK,
		GSI1SK:     key,
		ID:         team.ID.String(),
		Name:       team.Name,
		TeamColors: team.TeamColors,
	}
}

func teamPK(key string) string {
	return tablekeys.TEAM_KEY_PREFIX + key
}

func (db *DB) AddTeam(ctx context.Context, newTeam model.Team) (*model.Team, error) {
	dynamoItem := teamToDynamoItem(newTeam)

	putCond := expression.AttributeNotExists(expression.Name(tablekeys.PK))
	putExpr, err := expression.NewBuilder().WithCondition(putCond).Build()
	if err != nil {
		panic("addteam condition is broken: " + err.Error())
	}

	err = db.putItem(ctx, dynamoItem, withPutItemConditionExpression(putExpr))
	if err != nil {
		var condCheckErr *types.ConditionalCheckFailedException
		switch {
		case errors.As(err, &condCheckErr):
			return nil, ErrItemAlreadyExists
		}

		return nil, fmt.Errorf("failed to add team: %w", err)
	}

	return &newTeam, nil
}

type UpdateTeamInput struct {
	Name       *string
	TeamColors *[]string
}

func (db *DB) UpdateTeam(ctx context.Context, id uuid.UUID, updates UpdateTeamInput) (*model.Team, error) {
	existsCond := expression.AttributeExists(expression.Name(tablekeys.PK))
	var update expression.UpdateBuilder

	if updates.Name != nil {
		update = update.Set(expression.Name("Name"), expression.Value(*updates.Name))
	}
	if updates.TeamColors != nil {
		update = update.Set(expression.Name("TeamColors"), expression.Value(*updates.TeamColors))
	}

	expr, err := expression.NewBuilder().
		WithCondition(existsCond).
		WithUpdate(update).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression for updateTeam: %w", err)
	}

	key := teamPK(id.String())
	result, err := db.updateItem(ctx, key, key, withUpdateExpression(expr), withUpdateReturnValues(types.ReturnValueAllNew))
	if err != nil {
		var condCheckErr *types.ConditionalCheckFailedException
		switch {
		case errors.As(err, &condCheckErr):
			return nil, ErrItemNotFound
		}

		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	var dynamoItem teamDynamoItem
	err = attributevalue.UnmarshalMap(result.Attributes, &dynamoItem)
	if err != nil {
		return nil, fmt.Errorf("error reading result of team upate: %w", err)
	}

	teamResult := dynamoItem.toTeam()
	return &teamResult, nil
}

func (db *DB) GetTeam(ctx context.Context, id uuid.UUID) (*model.Team, error) {
	var teamItem teamDynamoItem

	key := teamPK(id.String())
	err := db.getItem(ctx, key, key, &teamItem)
	if err != nil {
		return nil, fmt.Errorf("failed to get team %q: %w", id, err)
	}

	team := teamItem.toTeam()
	return &team, nil
}

func (db *DB) GetAllTeams(ctx context.Context) ([]model.Team, error) {
	keyCond := expression.Key(tablekeys.GSI1PK).Equal(expression.Value(tablekeys.TEAM_GSI1_PK))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		panic("get all teams key condition broken: " + err.Error())
	}

	var teamItems []teamDynamoItem
	err = db.getManyOfEntity(ctx, &teamItems, withQueryKeyConditionExpression(expr), withQueryIndex(tablekeys.GSI1_INDEX))
	if err != nil {
		return nil, fmt.Errorf("failed to get all teams: %w", err)
	}

	teams := slices.Map(teamItems, func(item teamDynamoItem) model.Team {
		return item.toTeam()
	})

	return teams, nil
}
