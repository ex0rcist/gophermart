package usecases

import (
	"context"
	"errors"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/utils"
)

type UserAuthForm struct {
	Login    string `json:"login" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=3"`
}

var (
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
)

// проверить корректность логина и пароля
func UserAuth(ctx context.Context, s *storage.PGXStorage, d UserAuthForm) error {
	user, err := s.UserFindByLogin(ctx, d.Login)
	if err != nil {
		if err == entities.ErrRecordNotFound {
			return ErrInvalidLoginOrPassword
		}

		return err
	}

	err = utils.ComparePassword(user.Password, d.Password)
	if err != nil {
		return ErrInvalidLoginOrPassword
	}

	return nil
}
