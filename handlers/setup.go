package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/db"
	"github.com/mellena1/boston-archery-api/handlers/middleware"
)

func NewGin(logger *slog.Logger) *gin.Engine {
	if !IsLocal() {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()

	engine.Use(middleware.LoggingMiddleware(logger), requestid.New(), gin.Recovery())

	return engine
}

func NewDB(ctx context.Context) (*db.DB, error) {
	var dynamoClient *dynamodb.Client
	var err error
	if IsLocal() {
		dynamoClient, err = createLocalDynamoClient(ctx)
	} else {
		dynamoClient, err = createProdDynamoClient(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamo client: %w", err)
	}

	database := db.NewDB(os.Getenv("ARCHERY_TABLE_NAME"), dynamoClient)
	return database, nil
}

func IsLocal() bool {
	isOffline := os.Getenv("AWS_SAM_LOCAL")
	return isOffline == "true"
}

type AppVars struct {
	WebHost             *url.URL
	APIHost             *url.URL
	DiscordClientID     string
	DiscordClientSecret string
	JWTKey              []byte
}

func GetAppVars() (AppVars, error) {
	if IsLocal() {
		jwtKey, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_KEY"))
		if err != nil {
			return AppVars{}, err
		}
		webHost, err := url.Parse(os.Getenv("WEB_HOST"))
		if err != nil {
			return AppVars{}, err
		}
		apiHost, err := url.Parse(os.Getenv("API_HOST"))
		if err != nil {
			return AppVars{}, err
		}
		return AppVars{
			WebHost:             webHost,
			APIHost:             apiHost,
			DiscordClientID:     os.Getenv("DISCORD_CLIENT_ID"),
			DiscordClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
			JWTKey:              jwtKey,
		}, nil
	}
	// TODO: get prod secrets
	return AppVars{}, nil
}

func createLocalDynamoClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("localhost"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://dynamodb:8000"}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "local", SecretAccessKey: "local", SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}

func createProdDynamoClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg), nil
}
