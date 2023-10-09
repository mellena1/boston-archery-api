package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/handlers/errors"
)

const (
	captainRoleID      = "1031944371026800691"
	noStringsID        = "1047314424928686271"
	guildMemberInfoURL = "https://discord.com/api/users/@me/guilds/1031942797105823745/member"

	jwtCookieKey = "authToken"
	jwtTTL       = time.Hour * 2
)

var failedToGetTokenError = errors.Error{Msg: "failed to get token"}
var invalidStateError = errors.Error{Msg: "invalid state"}

type discordGuildMemberResp struct {
	Nickname string   `json:"nick"`
	Roles    []string `json:"roles"`
}

func (a *API) Callback(c *gin.Context) {
	ctx := c.Request.Context()

	cookieState, err := c.Cookie(stateCookieKey)
	if err != nil || cookieState != c.Query("state") {
		c.AbortWithStatusJSON(http.StatusBadRequest, invalidStateError)
		return
	}
	setStateCookie(c, "", -1)

	token, err := a.oauthConf.Exchange(ctx, c.Query("code"))
	if err != nil {
		a.logger.Error("failed to get token", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, failedToGetTokenError)
		return
	}

	res, err := a.oauthConf.Client(ctx, token).Get(guildMemberInfoURL)
	if err != nil {
		a.logger.Error("failed to get guild info", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, failedToGetTokenError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.UnauthorizedError)
		return
	}

	var memberInfo discordGuildMemberResp
	err = json.NewDecoder(res.Body).Decode(&memberInfo)
	if err != nil {
		a.logger.Error("failed to decode discord resp", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, failedToGetTokenError)
		return
	}

	// TODO: figure out what id to actually use
	jwt, err := a.jwtService.CreateJWT(memberInfo.Nickname, slices.Contains(memberInfo.Roles, noStringsID), jwtTTL)
	if err != nil {
		a.logger.Error("failed to make jwt", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, failedToGetTokenError)
		return
	}

	redirectURL := fmt.Sprintf("%s?authToken=%s", a.appVars.WebHost, jwt)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
