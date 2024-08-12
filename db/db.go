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

type putItemOption func(input *dynamodb.PutItemInput)

func withPutItemConditionExpression(expr expression.Expression) putItemOption {
	return func(input *dynamodb.PutItemInput) {
		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
		input.ConditionExpression = expr.Condition()
	}
}

type queryInputOption func(query *dynamodb.QueryInput)

func withQueryIndex(indexName string) queryInputOption {
	return func(query *dynamodb.QueryInput) {
		query.IndexName = &indexName
	}
}

func withQueryKeyConditionExpression(expr expression.Expression) queryInputOption {
	return func(query *dynamodb.QueryInput) {
		query.KeyConditionExpression = expr.KeyCondition()
		query.ExpressionAttributeNames = expr.Names()
		query.ExpressionAttributeValues = expr.Values()
	}
}

func (db *DB) putItem(ctx context.Context, item any, options ...putItemOption) error {
	marshaledItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: &db.tableName,
		Item:      marshaledItem,
	}

	for _, opt := range options {
		opt(input)
	}

	_, err = db.dynamoClient.PutItem(ctx, input)
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

func (db *DB) getManyOfEntity(ctx context.Context, v any, options ...queryInputOption) error {
	query := &dynamodb.QueryInput{
		TableName: &db.tableName,
	}
	for _, opt := range options {
		opt(query)
	}

	resp, err := db.dynamoClient.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("getManyOfEntity query failed: %w", err)
	}

	return attributevalue.UnmarshalListOfMaps(resp.Items, v)
}
