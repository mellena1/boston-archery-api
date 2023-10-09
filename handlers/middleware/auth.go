package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/auth"
	"github.com/mellena1/boston-archery-api/handlers/errors"
)

const jwtClaimsCtxKey = "jwtClaims"

var mustSetAuthHeaderError = errors.Error{Msg: "must include Authorization header"}
var invalidAuthTokenError = errors.Error{Msg: "auth token is invalid or expired"}
var mustBeAdminError = errors.Error{Msg: "must be an admin to perform this task"}

type JWTParser interface {
	ParseJWT(token string) (*auth.JWTClaims, error)
}

func ParseJWTMiddleware(parser JWTParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, mustSetAuthHeaderError)
			return
		}

		authSplit := strings.Split(authHeader, " ")
		if len(authSplit) != 2 || authSplit[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, mustSetAuthHeaderError)
			return
		}

		claims, err := parser.ParseJWT(authSplit[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, invalidAuthTokenError)
			return
		}

		SetJWTClaimsCtx(c, claims)

		c.Next()
	}
}

func MustBeAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetJWTClaimsCtx(c)

		if claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, invalidAuthTokenError)
			return
		}

		if !claims.IsAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, mustBeAdminError)
			return
		}

		c.Next()
	}
}

func GetJWTClaimsCtx(c *gin.Context) *auth.JWTClaims {
	claims, ok := c.Get(jwtClaimsCtxKey)
	if !ok {
		return nil
	}
	if claimsTyped, ok := claims.(*auth.JWTClaims); ok {
		return claimsTyped
	}
	return nil
}

func SetJWTClaimsCtx(c *gin.Context, claims *auth.JWTClaims) {
	c.Set(jwtClaimsCtxKey, claims)
}
