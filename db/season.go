package db

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db/tablekeys"
	"github.com/mellena1/boston-archery-api/model"
	"github.com/mellena1/boston-archery-api/slices"
)

const seasonEntityType = "Season"

type seasonDynamoItem struct {
	PK         string
	SK         string
	EntityType string
	GSI1PK     string
	GSI1SK     string
	Name       string
	StartDate  string
	EndDate    string
	ByeWeeks   []string
}

func (s seasonDynamoItem) toSeason() model.Season {
	return model.Season{
		ID:        uuid.MustParse(strings.Split(s.SK, "#")[1]),
		Name:      s.Name,
		StartDate: stringToDate(s.StartDate),
		EndDate:   stringToDate(s.EndDate),
		ByeWeeks: slices.Map(s.ByeWeeks, func(s string) time.Time {
			return stringToDate(s)
		}),
	}
}

type SeasonInput struct {
	Name      string
	StartDate time.Time
	EndDate   time.Time
	ByeWeeks  []time.Time
}

func (s SeasonInput) toDynamoItem(uuid uuid.UUID) seasonDynamoItem {
	key := seasonPK(uuid.String())
	gsi1Key := seasonPK(s.Name)
	return seasonDynamoItem{
		PK:         key,
		SK:         key,
		EntityType: seasonEntityType,
		GSI1PK:     gsi1Key,
		GSI1SK:     gsi1Key,
		Name:       s.Name,
		StartDate:  dateToString(s.StartDate),
		EndDate:    dateToString(s.EndDate),
		ByeWeeks: slices.Map(s.ByeWeeks, func(t time.Time) string {
			return dateToString(t)
		}),
	}
}

func (db *DB) AddSeason(ctx context.Context, newSeason SeasonInput) (*model.Season, error) {
	_, err := db.GetSeasonByName(ctx, newSeason.Name)
	switch {
	case errors.Is(err, ErrItemNotFound):
		// Item not found means we have no conflicts
	case err == nil:
		return nil, ErrItemAlreadyExists
	default:
		return nil, err
	}

	dynamoItem := newSeason.toDynamoItem(uuid.New())
	err = db.putItem(ctx, dynamoItem)
	if err != nil {
		return nil, err
	}

	season := dynamoItem.toSeason()
	return &season, nil
}

func (db *DB) UpdateSeason(ctx context.Context, id uuid.UUID, season SeasonInput) (*model.Season, error) {
	_, err := db.GetSeason(ctx, id)
	switch {
	case errors.Is(err, ErrItemNotFound):
		return nil, ErrItemNotFound
	case err == nil:
		// Item found means we do have an item with this ID already
	default:
		return nil, err
	}

	existingSeason, err := db.GetSeasonByName(ctx, season.Name)
	switch {
	case errors.Is(err, ErrItemNotFound):
		// Item not found means we have no conflicts
	case err == nil:
		if existingSeason.ID != id {
			return nil, ErrItemAlreadyExists
		}
	default:
		return nil, err
	}

	dynamoItem := season.toDynamoItem(id)
	err = db.putItem(ctx, dynamoItem)
	if err != nil {
		return nil, err
	}

	updatedSeason := dynamoItem.toSeason()
	return &updatedSeason, nil
}

func (db *DB) GetSeason(ctx context.Context, uuid uuid.UUID) (*model.Season, error) {
	var seasonItem seasonDynamoItem

	key := seasonPK(uuid.String())
	err := db.getItem(ctx, key, key, &seasonItem)
	if err != nil {
		return nil, err
	}

	season := seasonItem.toSeason()
	return &season, err
}

func (db *DB) GetSeasonByName(ctx context.Context, name string) (*model.Season, error) {
	key := seasonPK(name)
	keyCond := expression.Key("GSI1PK").Equal(expression.Value(key))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}
	resp, err := db.dynamoClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &db.tableName,
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Items) == 0 {
		return nil, ErrItemNotFound
	}

	var seasonItem seasonDynamoItem
	err = attributevalue.UnmarshalMap(resp.Items[0], &seasonItem)
	season := seasonItem.toSeason()
	return &season, err
}

func (db *DB) GetAllSeasons(ctx context.Context) ([]model.Season, error) {
	var seasonItems []seasonDynamoItem
	err := db.getAllOfEntity(ctx, seasonEntityType, &seasonItems)
	if err != nil {
		return nil, err
	}

	seasons := make([]model.Season, len(seasonItems))
	for i, item := range seasonItems {
		seasons[i] = item.toSeason()
	}

	return seasons, nil
}

func seasonPK(key string) string {
	return tablekeys.SEASON_KEY_PREFIX + key
}
