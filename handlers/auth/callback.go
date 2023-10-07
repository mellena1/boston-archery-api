package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/handlers"
)

const (
	captainRoleID      = "1031944371026800691"
	guildMemberInfoURL = "https://discord.com/api/users/@me/guilds/1031942797105823745/member"

	jwtCookieKey = "authToken"
	jwtTTL       = time.Hour * 2
)

var failedToGetTokenError = handlers.Error{Msg: "failed to get token"}
var invalidStateError = handlers.Error{Msg: "invalid state"}

type DiscordGuildMemberResp struct {
	Nickname string   `json:"nick"`
	Roles    []string `json:"roles"`
}

func (a *API) Callback(c *gin.Context) {
	ctx := c.Request.Context()

	cookieState, err := c.Cookie(stateCookieKey)
	if err != nil || cookieState != c.Query("state") {
		c.JSON(http.StatusBadRequest, invalidStateError)
		return
	}
	setStateCookie(c, "", -1)

	token, err := a.oauthConf.Exchange(ctx, c.Query("code"))
	if err != nil {
		a.logger.Error("failed to get token", "error", err)
		c.JSON(http.StatusInternalServerError, failedToGetTokenError)
		return
	}

	res, err := a.oauthConf.Client(ctx, token).Get(guildMemberInfoURL)
	if err != nil {
		a.logger.Error("failed to get guild info", "error", err)
		c.JSON(http.StatusInternalServerError, failedToGetTokenError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, handlers.UnauthorizedError)
		return
	}

	var memberInfo DiscordGuildMemberResp
	err = json.NewDecoder(res.Body).Decode(&memberInfo)
	if err != nil {
		a.logger.Error("failed to decode discord resp", "error", err)
		c.JSON(http.StatusInternalServerError, failedToGetTokenError)
		return
	}

	jwt, err := a.jwtService.CreateJWT(memberInfo.Nickname, slices.Contains(memberInfo.Roles, captainRoleID), jwtTTL)
	if err != nil {
		a.logger.Error("failed to make jwt", "error", err)
		c.JSON(http.StatusInternalServerError, failedToGetTokenError)
		return
	}

	redirectURL := fmt.Sprintf("%s?authToken=%s", a.appVars.WebHost, jwt)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
