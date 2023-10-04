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
	logger    *slog.Logger = logging.NewLogger(slog.LevelDebug)
	database  *db.DB
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dynamoClient, err := db.CreateLocalClient(ctx)
	if err != nil {
		logger.Error("failed to create dynamo client", "error", err)
		panic(err)
	}

	database = db.NewDB("ArcheryTag", "EntityTypeGSI", dynamoClient)

	logger.Info("Gin cold start")
	r := handlers.NewGin(logger)
	group := r.Group("/seasons")
	{
		group.GET("", GetSeasons)
		group.POST("", PostSeason)
	}

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
