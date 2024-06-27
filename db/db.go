package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mellena1/boston-archery-api/db/tablekeys"
)

type DynamoDBClient interface {
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(context.Context, *dynamodb.QueryInput, ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	UpdateItem(context.Context, *dynamodb.UpdateItemInput, ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

type DB struct {
	tableName           string
	entityTypeIndexName string
	dynamoClient        DynamoDBClient
}

func NewDB(tableName string, entityTypeIndexName string, dynamoClient DynamoDBClient) *DB {
	return &DB{
		tableName:           tableName,
		entityTypeIndexName: entityTypeIndexName,
		dynamoClient:        dynamoClient,
	}
}

func (db *DB) putItem(ctx context.Context, item any) error {
	marshaledItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	_, err = db.dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &db.tableName,
		Item:      marshaledItem,
	})
	if err != nil {
		return fmt.Errorf("error with PutItem call: %w", err)
	}

	return nil
}

func (db *DB) getItem(ctx context.Context, PK string, SK string, v any) error {
	resp, err := db.dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &db.tableName,
		Key: map[string]types.AttributeValue{
			tablekeys.PK: &types.AttributeValueMemberS{Value: PK},
			tablekeys.SK: &types.AttributeValueMemberS{Value: SK},
		},
	})
	if err != nil {
		return fmt.Errorf("error with GetItem call: %w", err)
	}

	if len(resp.Item) == 0 {
		return ErrItemNotFound
	}

	return attributevalue.UnmarshalMap(resp.Item, v)
}

func (db *DB) getAllOfEntity(ctx context.Context, entityType string, vs any) error {
	keyCond := expression.Key(tablekeys.ENTITY_TYPE).Equal(expression.Value(entityType))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return fmt.Errorf("making entity expression failed")
	}

	resp, err := db.dynamoClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &db.tableName,
		IndexName:                 &db.entityTypeIndexName,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return fmt.Errorf("query on entity GSI failed: %w", err)
	}

	return attributevalue.UnmarshalListOfMaps(resp.Items, vs)
}
