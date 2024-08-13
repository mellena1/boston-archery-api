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
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/auth"
	"github.com/mellena1/boston-archery-api/db"
	"github.com/mellena1/boston-archery-api/handlers"
	authHandler "github.com/mellena1/boston-archery-api/handlers/auth"
	"github.com/mellena1/boston-archery-api/handlers/middleware"
	"github.com/mellena1/boston-archery-api/handlers/players"
	"github.com/mellena1/boston-archery-api/handlers/seasons"
	"github.com/mellena1/boston-archery-api/handlers/teams"
	"github.com/mellena1/boston-archery-api/logging"
)

var (
	ginLambda *ginadapter.GinLambda
)

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
		appVars:   &appVars,
	}

	r := handlers.NewGin(logger)
	r.Use(cors.New(cors.Config{
		// TODO: make configurable for prod
		AllowOrigins:  []string{"http://localhost:*"},
		AllowHeaders:  []string{"Authorization", "Content-Type"},
		AllowWildcard: true,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodOptions,
			http.MethodPost,
			http.MethodPut,
		},
	}))

	apiV1Group := r.Group("/api/v1")
	adminGroup := apiV1Group.Group("", middleware.ParseJWTMiddleware(api.jwtParser), middleware.MustBeAdminMiddleware())

	api.addPlayerAPIs(apiV1Group, adminGroup)
	api.addTeamAPIs(apiV1Group, adminGroup)
	api.addSeasonAPIs(apiV1Group, adminGroup)
	api.addAuthAPIs(apiV1Group)

	ginLambda = ginadapter.New(r)
}

type API struct {
	logger    *slog.Logger
	db        *db.DB
	jwtParser *auth.JWTService
	appVars   *handlers.AppVars
}

func (api *API) addSeasonAPIs(apiV1Group, adminGroup *gin.RouterGroup) {
	seasonApi := seasons.NewAPI(api.logger, api.db)

	apiV1Group.GET("/seasons", seasonApi.GetSeasons)

	seasonAdminGroup := adminGroup.Group("/season")
	{
		seasonAdminGroup.POST("", seasonApi.PostSeason)
		seasonAdminGroup.PUT("/:id", seasonApi.PutSeason)
	}
}

func (api *API) addPlayerAPIs(apiV1Group, adminGroup *gin.RouterGroup) {
	playerApi := players.NewAPI(api.logger, api.db)

	apiV1Group.GET("/players", playerApi.GetPlayers)

	playerGroup := apiV1Group.Group("/player")
	{
		playerGroup.GET("/:id", playerApi.GetPlayer)
	}

	playerAdminGroup := adminGroup.Group("/player")
	{
		playerAdminGroup.POST("", playerApi.PostPlayer)
		playerAdminGroup.PUT("", playerApi.PutPlayer)
	}
}

func (api *API) addTeamAPIs(apiV1Group, adminGroup *gin.RouterGroup) {
	teamApi := teams.NewAPI(api.logger, api.db)

	apiV1Group.GET("/teams", teamApi.GetTeams)

	teamAdminGroup := adminGroup.Group("/team")
	{
		teamAdminGroup.POST("", teamApi.PostTeam)
		teamAdminGroup.PUT("", teamApi.PutTeam)
	}
}

func (api *API) addAuthAPIs(apiV1Group *gin.RouterGroup) {
	authApi := authHandler.NewAPI(api.logger, api.jwtParser, api.appVars)

	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.GET("/login", authApi.Login)
		authGroup.GET("/callback", authApi.Callback)
	}
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
