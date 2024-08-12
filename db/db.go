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
	TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

type DB struct {
	tableName    string
	dynamoClient DynamoDBClient
}

func NewDB(tableName string, dynamoClient DynamoDBClient) *DB {
	return &DB{
		tableName:    tableName,
		dynamoClient: dynamoClient,
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

type updateInputOption func(input *dynamodb.UpdateItemInput)

func withUpdateExpression(expr expression.Expression) updateInputOption {
	return func(input *dynamodb.UpdateItemInput) {
		input.ConditionExpression = expr.Condition()
		input.UpdateExpression = expr.Update()
		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
	}
}

func withUpdateReturnValues(ret types.ReturnValue) updateInputOption {
	return func(input *dynamodb.UpdateItemInput) {
		input.ReturnValues = ret
	}
}

func (db *DB) putItem(ctx context.Context, item any, opts ...putItemOption) error {
	marshaledItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: &db.tableName,
		Item:      marshaledItem,
	}

	for _, opt := range opts {
		opt(input)
	}

	_, err = db.dynamoClient.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("error with PutItem call: %w", err)
	}

	return nil
}

func (db *DB) updateItem(ctx context.Context, PK string, SK string, opts ...updateInputOption) (*dynamodb.UpdateItemOutput, error) {
	input := &dynamodb.UpdateItemInput{
		TableName: &db.tableName,
		Key: map[string]types.AttributeValue{
			tablekeys.PK: &types.AttributeValueMemberS{Value: PK},
			tablekeys.SK: &types.AttributeValueMemberS{Value: SK},
		},
	}

	for _, opt := range opts {
		opt(input)
	}

	result, err := db.dynamoClient.UpdateItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error with UpdateItem call: %w", err)
	}

	return result, nil
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

func (db *DB) getManyOfEntity(ctx context.Context, v any, opts ...queryInputOption) error {
	query := &dynamodb.QueryInput{
		TableName: &db.tableName,
	}
	for _, opt := range opts {
		opt(query)
	}

	resp, err := db.dynamoClient.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("getManyOfEntity query failed: %w", err)
	}

	return attributevalue.UnmarshalListOfMaps(resp.Items, v)
}
