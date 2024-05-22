package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

const stateCookieKey = "login-state"

func (a *API) Login(c *gin.Context) {
	state := generateState()
	setStateCookie(c, state, 300)

	c.Redirect(http.StatusTemporaryRedirect, a.oauthConf.AuthCodeURL(state))
}

func generateState() string {
	b := make([]byte, 128)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func setStateCookie(c *gin.Context, val string, maxAge int) {
	c.SetCookie(stateCookieKey, val, maxAge, "/api/v1/auth", "", true, true)
}
