package usecases

import (
	"context"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/pkg/jwt"
)

// создать JWT токен для логина
func UserCreateJWT(ctx context.Context, secret entities.Secret, login string) (string, error) {
	token, err := jwt.CreateJWT(secret, login)
	if err != nil {
		return "", err
	}

	return token, nil
}
