package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/mellena1/boston-archery-api/auth"
	"github.com/mellena1/boston-archery-api/handlers"
	"github.com/mellena1/boston-archery-api/logging"
	"golang.org/x/oauth2"
)

var (
	ginLambda *ginadapter.GinLambda
)

type API struct {
	logger     *slog.Logger
	oauthConf  *oauth2.Config
	jwtService *auth.JWTService
	appVars    *handlers.AppVars
}

func init() {
	logger := logging.NewLogger(slog.LevelDebug)
	logger.Info("Gin cold start")

	appVars, err := handlers.GetAppVars()
	if err != nil {
		logger.Error("failed to get app vars", "error", err)
		panic(err)
	}

	api := API{
		logger: logger,
		oauthConf: &oauth2.Config{
			RedirectURL:  fmt.Sprintf("%s/auth/callback", appVars.APIHost),
			ClientID:     appVars.DiscordClientID,
			ClientSecret: appVars.DiscordClientSecret,
			Scopes:       []string{"guilds.members.read"},
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://discord.com/api/oauth2/authorize",
				TokenURL:  "https://discord.com/api/oauth2/token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
		},
		jwtService: auth.NewJWTService(appVars.JWTKey),
		appVars:    &appVars,
	}

	r := handlers.NewGin(logger)
	group := r.Group("/auth")
	{
		group.GET("/login", api.Login)
		group.GET("/callback", api.Callback)
	}
	r.GET("/login", api.Login)
	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
