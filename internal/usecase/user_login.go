package usecase

import (
	"context"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/utils"
	"github.com/ex0rcist/gophermart/pkg/jwt"
)

type loginUsecase struct {
	storage        storage.IPGXStorage
	repo           domain.IUserRepository
	secret         entities.Secret
	contextTimeout time.Duration
}

func NewLoginUsecase(storage storage.IPGXStorage, repo domain.IUserRepository, secret entities.Secret, timeout time.Duration) domain.ILoginUsecase {
	return &loginUsecase{storage: storage, repo: repo, secret: secret, contextTimeout: timeout}
}

func (uc *loginUsecase) Call(ctx context.Context, form domain.LoginRequest) (string, error) {
	// находим пользователя
	user, err := uc.GetUserByLogin(ctx, form)
	if err != nil {
		return "", err
	}

	// создаем JWT токен
	token, err := uc.CreateAccessToken(user, uc.secret, domain.LoginTokenLifetime)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *loginUsecase) GetUserByLogin(ctx context.Context, req domain.LoginRequest) (*domain.User, error) {
	tCtx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	user, err := uc.repo.UserFindByLogin(tCtx, req.Login)
	if err != nil {
		if err == storage.ErrRecordNotFound {
			return nil, domain.ErrInvalidLoginOrPassword
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
