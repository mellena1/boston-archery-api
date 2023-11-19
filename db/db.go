package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
