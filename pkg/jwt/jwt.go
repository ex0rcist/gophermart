package jwt

import (
	"errors"
	"time"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/golang-jwt/jwt/v4"
)

var ErrInvalidToken = errors.New("invalid JWT token")
var tokenLifetime = time.Minute * 5000

type GMClaims struct {
	jwt.RegisteredClaims
	Login string
}

func CreateJWT(key entities.Secret, login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, GMClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
		},
		Login: login,
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(key entities.Secret, rawToken string) (string, time.Time, error) {
	claims := new(GMClaims)
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}

	token, err := jwt.ParseWithClaims(rawToken, claims, keyFunc)
	if err != nil {
		return "", time.Time{}, err
	}

	if !token.Valid {
		return "", time.Time{}, ErrInvalidToken
	}

	if claims.ExpiresAt == nil {
		return "", time.Time{}, ErrInvalidToken
	}

	return claims.Login, claims.ExpiresAt.Time, nil
}
