package usecases

import (
	"context"
	"errors"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/models"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/utils"
	"github.com/jinzhu/copier"
)

type UserRegisterForm struct {
	Login    string `json:"login" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=3"`
}

var (
	ErrUserAlreadyExists = errors.New("login already exists")
)

// создать нового пользователя
func UserRegister(ctx context.Context, s *storage.PGXStorage, d UserRegisterForm) (*models.User, error) {
	// проверяем, что логин не занят
	existingUser, err := s.UserFindByLogin(ctx, d.Login)
	if err != nil && err != entities.ErrRecordNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// генерируем пароль
	d.Password, err = utils.HashPassword(d.Password)
	if err != nil {
		return nil, err
	}

	// копируем данные в DTO
	dto := storage.CreateUserDTO{}
	err = copier.Copy(&dto, &d)
	if err != nil {
		return nil, err
	}

	// создаем пользователя
	newUser, err := s.UserCreate(ctx, dto)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
