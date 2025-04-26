package auth

import (
	"errors"
	"time"

	"github.com/Anacardo89/doubleOrNothingDice/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ClientID string `json:"client_id"`
	jwt.RegisteredClaims
}

func GenerateToken(clientID string) (string, error) {
	expiration := time.Now().Add(time.Duration(config.AppConfig.JWT.ExpiryMinutes) * time.Minute)
	claims := &Claims{
		ClientID: clientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
