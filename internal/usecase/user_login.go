package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/storage/repository"
	"github.com/ex0rcist/gophermart/internal/utils"
	"github.com/ex0rcist/gophermart/pkg/jwt"
)

var ErrInvalidLoginOrPassword = errors.New("invalid login or password")

type ILoginUsecase interface {
	Call(ctx context.Context, form LoginRequest) (string, error)
	GetUserByLogin(ctx context.Context, req LoginRequest) (*domain.User, error)
	ComparePassword(user *domain.User, password string) error
	CreateAccessToken(user *domain.User, secret entities.Secret, lifetime time.Duration) (string, error)
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=3"`
}

type loginUsecase struct {
	storage        storage.IPGXStorage
	repo           repository.IUserRepository
	secret         entities.Secret
	contextTimeout time.Duration
}

func NewLoginUsecase(storage storage.IPGXStorage, repo repository.IUserRepository, secret entities.Secret, timeout time.Duration) ILoginUsecase {
	return &loginUsecase{storage: storage, repo: repo, secret: secret, contextTimeout: timeout}
}

func (uc *loginUsecase) Call(ctx context.Context, form LoginRequest) (string, error) {
	// находим пользователя
	user, err := uc.GetUserByLogin(ctx, form)
	if err != nil {
		return "", err
	}

	// создаем JWT токен
	token, err := uc.CreateAccessToken(user, uc.secret, jwt.LoginTokenLifetime)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *loginUsecase) GetUserByLogin(ctx context.Context, req LoginRequest) (*domain.User, error) {
	tCtx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	user, err := uc.repo.UserFindByLogin(tCtx, req.Login)
	if err != nil {
		if err == storage.ErrRecordNotFound {
			return nil, ErrInvalidLoginOrPassword
		}

		return nil, err
	}

	return user, nil
}

func (uc *loginUsecase) ComparePassword(user *domain.User, password string) error {
	return utils.ComparePassword(user.Password, password)
}

func (uc *loginUsecase) CreateAccessToken(user *domain.User, secret entities.Secret, lifetime time.Duration) (string, error) {
	token, err := jwt.CreateJWT(secret, user.Login, lifetime)
	if err != nil {
		return "", err
	}

	return token, nil
}
