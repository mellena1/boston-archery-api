// Package classification Boston Archery API
//
// API for Boston Archery
// Version: 1.0.0
// Schemes: http
// Host: localhost:3000
// BasePath: /api/v1
//
// Consumes:
//   - application/json
//
// Produces:
//   - application/json
//
// swagger:meta
package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/mellena1/boston-archery-api/auth"
	"github.com/mellena1/boston-archery-api/db"
	"github.com/mellena1/boston-archery-api/handlers"
	"github.com/mellena1/boston-archery-api/handlers/middleware"
	"github.com/mellena1/boston-archery-api/logging"
)

var (
	ginLambda *ginadapter.GinLambda
)

type SeasonDB interface {
	AddSeason(ctx context.Context, newSeason db.SeasonInput) (*db.Season, error)
	GetAllSeasons(ctx context.Context) ([]db.Season, error)
	GetSeasonByName(ctx context.Context, name string) (*db.Season, error)
}

type API struct {
	logger    *slog.Logger
	db        SeasonDB
	jwtParser middleware.JWTParser
}

func init() {
	logger := logging.NewLogger(slog.LevelDebug)
	logger.Info("Gin cold start")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	appVars, err := handlers.GetAppVars()
	if err != nil {
		logger.Error("failed to get app vars", "error", err)
		panic(err)
	}

	database, err := handlers.NewDB(ctx)
	if err != nil {
		logger.Error("failed to create database", "error", err)
		panic(err)
	}

	api := API{
		logger:    logger,
		db:        database,
		jwtParser: auth.NewJWTService(appVars.JWTKey),
	}

	r := handlers.NewGin(logger)
	r.Use(cors.New(cors.Config{
		// TODO: make configurable for prod
		AllowOrigins:  []string{"http://localhost:*"},
		AllowWildcard: true,
	}))
	group := r.Group("/api/v1/seasons")
	{
		group.GET("", api.GetSeasons)

		adminGroup := group.Group("", middleware.ParseJWTMiddleware(api.jwtParser), middleware.MustBeAdminMiddleware())
		{
			adminGroup.POST("", api.PostSeason)
		}
	}

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
