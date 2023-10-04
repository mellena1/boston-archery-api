package db

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db/tablekeys"
	"github.com/mellena1/boston-archery-api/models"
	"github.com/mellena1/boston-archery-api/slices"
)

const seasonEntityType = "Season"

type SeasonService interface {
	AddSeason(ctx context.Context, newSeason SeasonInput) error
	GetSeason(ctx context.Context, uuid uuid.UUID) (*models.Season, error)
	GetAllSeasons(ctx context.Context) ([]models.Season, error)
}

type seasonDynamoItem struct {
	PK         string
	SK         string
	EntityType string
	Name       string
	StartDate  string
	EndDate    string
	ByeWeeks   []string
}

func (s seasonDynamoItem) toSeason() models.Season {
	return models.Season{
		UUID:      uuid.MustParse(strings.Split(s.SK, "#")[1]),
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
	key := seasonPK(uuid)
	return seasonDynamoItem{
		PK:         key,
		SK:         key,
		EntityType: seasonEntityType,
		Name:       s.Name,
		StartDate:  dateToString(s.StartDate),
		EndDate:    dateToString(s.EndDate),
		ByeWeeks: slices.Map(s.ByeWeeks, func(t time.Time) string {
			return dateToString(t)
		}),
	}
}

func (db *DB) AddSeason(ctx context.Context, newSeason SeasonInput) error {
	item, err := attributevalue.MarshalMap(newSeason.toDynamoItem(uuid.New()))
	if err != nil {
		return err
	}
	_, err = db.dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &db.tableName,
		Item:      item,
	})
	return err
}

func (db *DB) GetSeason(ctx context.Context, uuid uuid.UUID) (*models.Season, error) {
	key := seasonPK(uuid)
	resp, err := db.dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &db.tableName,
		Key: map[string]types.AttributeValue{
			tablekeys.PK: &types.AttributeValueMemberS{Value: tablekeys.SEASON_KEY_PREFIX},
			tablekeys.SK: &types.AttributeValueMemberS{Value: key},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Item) == 0 {
		return nil, ErrItemNotFound
	}

	var seasonItem seasonDynamoItem
	err = attributevalue.UnmarshalMap(resp.Item, &seasonItem)
	season := seasonItem.toSeason()
	return &season, err
}

func (db *DB) GetAllSeasons(ctx context.Context) ([]models.Season, error) {
	keyCond := expression.Key(tablekeys.ENTITY_TYPE).Equal(expression.Value(seasonEntityType))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	resp, err := db.dynamoClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &db.tableName,
		IndexName:                 &db.entityTypeIndexName,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}

	var seasonItems []seasonDynamoItem
	err = attributevalue.UnmarshalListOfMaps(resp.Items, &seasonItems)

	seasons := []models.Season{}
	for _, item := range seasonItems {
		seasons = append(seasons, item.toSeason())
	}

	return seasons, err
}

func seasonPK(uuid uuid.UUID) string {
	return tablekeys.SEASON_KEY_PREFIX + uuid.String()
}
