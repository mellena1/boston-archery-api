package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/mellena1/boston-archery-api/db"
	"github.com/mellena1/boston-archery-api/handlers"
	"github.com/mellena1/boston-archery-api/logging"
)

var (
	ginLambda *ginadapter.GinLambda
)

type SeasonDB interface {
	AddSeason(ctx context.Context, newSeason db.SeasonInput) error
	GetAllSeasons(ctx context.Context) ([]db.Season, error)
}

type API struct {
	logger *slog.Logger
	db     SeasonDB
}

func init() {
	logger := logging.NewLogger(slog.LevelDebug)
	logger.Info("Gin cold start")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	database, err := handlers.NewDB(ctx)
	if err != nil {
		logger.Error("failed to create database", "error", err)
		panic(err)
	}

	api := API{
		logger: logger,
		db:     database,
	}

	r := handlers.NewGin(logger)
	group := r.Group("/api/v1/seasons")
	{
		group.GET("", api.GetSeasons)
		group.POST("", api.PostSeason)
	}

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
