package usecase

import (
	"context"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/pkg/jwt"
)

type registerUsecase struct {
	repo           domain.IUserRepository
	contextTimeout time.Duration
}

func NewRegisterUsecase(repo domain.IUserRepository, timeout time.Duration) domain.IRegisterUsecase {
	return &registerUsecase{repo: repo, contextTimeout: timeout}
}

func (uc *registerUsecase) GetUserByLogin(c context.Context, req domain.RegisterRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	return uc.repo.UserFindByLogin(ctx, req.Login)
}

func (uc *registerUsecase) CreateUser(c context.Context, login string, password string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	return uc.repo.UserCreate(ctx, login, password)
}

func (uc *registerUsecase) CreateAccessToken(user *domain.User, secret entities.Secret, lifetime time.Duration) (string, error) {
	token, err := jwt.CreateJWT(secret, user.Login, lifetime)
	if err != nil {
		return "", err
	}

	return token, nil
}
