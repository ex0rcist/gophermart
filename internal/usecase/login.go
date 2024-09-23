package usecase

import (
	"context"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/utils"
	"github.com/ex0rcist/gophermart/pkg/jwt"
)

type loginUsecase struct {
	repo           domain.IUserRepository
	contextTimeout time.Duration
}

func NewLoginUsecase(repo domain.IUserRepository, timeout time.Duration) domain.ILoginUsecase {
	return &loginUsecase{repo: repo, contextTimeout: timeout}
}

func (uc *loginUsecase) GetUserByLogin(c context.Context, req domain.LoginRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	return uc.repo.UserFindByLogin(ctx, req.Login)
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
