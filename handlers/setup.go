package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/db"
	"github.com/mellena1/boston-archery-api/logging"
)

func NewGin(logger *slog.Logger) *gin.Engine {
	engine := gin.New()

	engine.Use(logging.GinMiddlewareLogger(logger), gin.Recovery())

	return engine
}

func NewDB(ctx context.Context) (*db.DB, error) {
	dynamoClient, err := db.CreateLocalClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamo client: %w", err)
	}

	database := db.NewDB(os.Getenv("ARCHERY_TABLE_NAME"), "EntityTypeGSI", dynamoClient)
	return database, nil
}

func IsLocal() bool {
	isOffline := os.Getenv("IS_OFFLINE")
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
