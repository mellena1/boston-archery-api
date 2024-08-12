package db

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db/tablekeys"
	"github.com/mellena1/boston-archery-api/model"
	"github.com/mellena1/boston-archery-api/slices"
)

type seasonDynamoItem struct {
	PK        string
	SK        string
	GSI1PK    string
	GSI1SK    string
	ID        string
	Name      string
	StartDate string
	EndDate   string
	ByeWeeks  []string
}

func (s seasonDynamoItem) toSeason() model.Season {
	return model.Season{
		ID:        uuid.MustParse(s.ID),
		Name:      s.Name,
		StartDate: stringToDate(s.StartDate),
		EndDate:   stringToDate(s.EndDate),
		ByeWeeks:  strSliceToDates(s.ByeWeeks),
	}
}

func seasonToDynamoItem(season model.Season) seasonDynamoItem {
	pk := seasonPK(season.ID.String())
	return seasonDynamoItem{
		PK:        pk,
		SK:        pk,
		GSI1PK:    tablekeys.SEASON_GSI1_PK,
		GSI1SK:    seasonGSI1SK(season.ID, season.StartDate),
		ID:        season.ID.String(),
		Name:      season.Name,
		StartDate: dateToString(season.StartDate),
		EndDate:   dateToString(season.EndDate),
		ByeWeeks:  dateSliceToStrs(season.ByeWeeks),
	}
}

func (db *DB) AddSeason(ctx context.Context, season model.Season) (*model.Season, error) {
	putCond := expression.AttributeNotExists(expression.Name(tablekeys.PK))
	putExpr, err := expression.NewBuilder().WithCondition(putCond).Build()
	if err != nil {
		panic("addseason condition is invalid")
	}

	return db.putSeason(ctx, season, putExpr)
}

type UpdateSeasonInput struct {
	Name      *string
	StartDate *time.Time
	EndDate   *time.Time
	ByeWeeks  *[]time.Time
}

func (db *DB) UpdateSeason(ctx context.Context, id uuid.UUID, updates UpdateSeasonInput) (*model.Season, error) {
	cond := expression.AttributeExists(expression.Name(tablekeys.PK))
	var update expression.UpdateBuilder

	if updates.Name != nil {
		update = update.Set(
			expression.Name("Name"),
			expression.Value(*updates.Name),
		)
	}
	if updates.StartDate != nil {
		update = update.Set(
			expression.Name("StartDate"),
			expression.Value(dateToString(*updates.StartDate)),
		).Set(
			expression.Name(tablekeys.GSI1SK),
			expression.Value(seasonGSI1SK(id, *updates.StartDate)),
		)
	}
	if updates.EndDate != nil {
		update = update.Set(
			expression.Name("EndDate"),
			expression.Value(dateToString(*updates.EndDate)),
		)
	}
	if updates.ByeWeeks != nil {
		update = update.Set(
			expression.Name("ByeWeeks"),
			expression.Value(dateSliceToStrs(*updates.ByeWeeks)),
		)
	}

	expr, err := expression.NewBuilder().
		WithCondition(cond).
		WithUpdate(update).
		Build()
	if err != nil {
		panic("updateseason condition is invalid: " + err.Error())
	}

	result, err := db.updateItem(ctx, seasonPK(id.String()), seasonPK(id.String()), withUpdateExpression(expr), withUpdateReturnValues(types.ReturnValueAllNew))
	if err != nil {
		return nil, fmt.Errorf("error updating season: %w", err)
	}

	var dynamoItem seasonDynamoItem
	err = attributevalue.UnmarshalMap(result.Attributes, &dynamoItem)
	if err != nil {
		return nil, fmt.Errorf("error reading result of season upate: %w", err)
	}

	seasonResult := dynamoItem.toSeason()
	return &seasonResult, nil
}

func (db *DB) putSeason(ctx context.Context, season model.Season, cond expression.Expression) (*model.Season, error) {
	err := db.putItem(ctx, seasonToDynamoItem(season), withPutItemConditionExpression(cond))
	if err != nil {
		return nil, err
	}

	return &season, nil
}

func (db *DB) GetSeason(ctx context.Context, uuid uuid.UUID) (*model.Season, error) {
	var seasonItem seasonDynamoItem

	key := seasonPK(uuid.String())
	err := db.getItem(ctx, key, key, &seasonItem)
	if err != nil {
		return nil, err
	}

	// TODO: error handling

	season := seasonItem.toSeason()
	return &season, err
}

func (db *DB) GetAllSeasons(ctx context.Context) ([]model.Season, error) {
	keyCond := expression.Key(tablekeys.GSI1PK).Equal(expression.Value(tablekeys.SEASON_GSI1_PK))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		panic("get all seasons key condition is broken: " + err.Error())
	}

	var seasonItems []seasonDynamoItem
	err = db.getManyOfEntity(ctx, &seasonItems, withQueryKeyConditionExpression(expr), withQueryIndex(tablekeys.GSI1_INDEX))
	if err != nil {
		return nil, err
	}

	seasons := slices.Map(seasonItems, func(item seasonDynamoItem) model.Season {
		return item.toSeason()
	})
	return seasons, nil
}

func seasonPK(key string) string {
	return tablekeys.SEASON_KEY_PREFIX + key
}

func seasonGSI1SK(id uuid.UUID, startDate time.Time) string {
	return seasonPK(fmt.Sprintf("%s#%s", dateToString(startDate), id.String()))
}
