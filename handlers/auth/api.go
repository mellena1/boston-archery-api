package auth

import (
	"fmt"
	"log/slog"

	"github.com/mellena1/boston-archery-api/auth"
	"github.com/mellena1/boston-archery-api/handlers"
	"golang.org/x/oauth2"
)

type API struct {
	logger     *slog.Logger
	oauthConf  *oauth2.Config
	jwtService *auth.JWTService
	appVars    *handlers.AppVars
}

func NewAPI(logger *slog.Logger, jwtService *auth.JWTService, appVars *handlers.AppVars) *API {
	return &API{
		logger:     logger,
		jwtService: jwtService,
		appVars:    appVars,
		oauthConf: &oauth2.Config{
			RedirectURL:  fmt.Sprintf("%s/api/v1/auth/callback", appVars.APIHost),
			ClientID:     appVars.DiscordClientID,
			ClientSecret: appVars.DiscordClientSecret,
			Scopes:       []string{"guilds.members.read"},
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://discord.com/api/oauth2/authorize",
				TokenURL:  "https://discord.com/api/oauth2/token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
		},
	}
}
