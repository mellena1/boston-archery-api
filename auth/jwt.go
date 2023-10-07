package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var signingMethod = jwt.SigningMethodHS256

type JWTService struct {
	key []byte
}

func NewJWTService(key []byte) *JWTService {
	return &JWTService{key: key}
}

type JWTClaims struct {
	Nickname string `json:"nickname"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

func (j *JWTService) CreateJWT(nickname string, isAdmin bool, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(signingMethod, JWTClaims{
		Nickname: nickname,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	})

	return token.SignedString(j.key)
}

func (j *JWTService) ParseJWT(token string) (*JWTClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return j.key, nil
	}, jwt.WithValidMethods([]string{signingMethod.Name}))
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if typedClaims, ok := parsedToken.Claims.(*JWTClaims); ok && parsedToken.Valid {
		return typedClaims, nil
	}
	return nil, fmt.Errorf("invalid claims from token or token invalid")
}
