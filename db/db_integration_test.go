//go:build integration

package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	dynamodblocal "github.com/abhirockzz/dynamodb-local-testcontainers-go"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mellena1/boston-archery-api/db/tablekeys"
)

var dynamodbTestContainer *dynamodblocal.DynamodbLocalContainer
var dynamoClient *dynamodb.Client
var db *DB

const tableName = "ArcheryTag"

func TestMain(m *testing.M) {
	ctx := context.Background()

	err := setupDynamo(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer shutdownDynamo(ctx)

	os.Exit(m.Run())
}

func setupDynamo(ctx context.Context) error {
	var err error
	dynamodbTestContainer, err = dynamodblocal.RunContainer(ctx)
	if err != nil {
		return fmt.Errorf("error starting dynamo testcontainer: %w", err)
	}

	dynamoClient, err = dynamodbTestContainer.GetDynamoDBClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting dynamo client: %w", err)
	}

	err = makeTable(ctx)
	if err != nil {
		return err
	}

	db = NewDB(tableName, dynamoClient)

	return nil
}

func makeTable(ctx context.Context) error {
	_, err := dynamoClient.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(tablekeys.PK),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(tablekeys.SK),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(tablekeys.GSI1PK),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(tablekeys.GSI1SK),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(tablekeys.PK),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String(tablekeys.SK),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String(tablekeys.GSI1_INDEX),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String(tablekeys.GSI1PK),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String(tablekeys.GSI1SK),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error making table: %w", err)
	}

	return nil
}

func resetTable(ctx context.Context) {
	_, err := dynamoClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		fmt.Printf("failed to delete table: %s", err)
	}

	err = makeTable(ctx)
	if err != nil {
		fmt.Printf("failed to remake table: %s", err)
	}
}

func shutdownDynamo(ctx context.Context) {
	if dynamodbTestContainer == nil {
		return
	}

	err := dynamodbTestContainer.Terminate(ctx)
	if err != nil {
		fmt.Printf("error terminating dynamo testcontainer: %s\n", err)
	}
}
