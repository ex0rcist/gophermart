package usecase

import (
	"context"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/utils"
	"github.com/ex0rcist/gophermart/pkg/jwt"
)

type registerUsecase struct {
	repo           domain.IUserRepository
	contextTimeout time.Duration
	secret         entities.Secret
}

func NewRegisterUsecase(repo domain.IUserRepository, secret entities.Secret, timeout time.Duration) domain.IRegisterUsecase {
	return &registerUsecase{repo: repo, contextTimeout: timeout, secret: secret}
}

func (uc *registerUsecase) Call(ctx context.Context, form domain.RegisterRequest) (string, error) {
	// находим пользователя
	existingUser, err := uc.GetUserByLogin(ctx, form)
	if err != nil && err != entities.ErrRecordNotFound {
		return "", err
	}
	if existingUser != nil {
		return "", domain.ErrUserAlreadyExists
	}

	// генерируем пароль
	form.Password, err = utils.HashPassword(form.Password)
	if err != nil {
		return "", err
	}

	// создаем пользователя
	newUser, err := uc.CreateUser(ctx, form.Login, form.Password)
	if err != nil {
		return "", err
	}

	// создаем JWT токен
	token, err := uc.CreateAccessToken(newUser, uc.secret, domain.LoginTokenLifetime)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *registerUsecase) GetUserByLogin(ctx context.Context, req domain.RegisterRequest) (*domain.User, error) {
	tCtx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	return uc.repo.UserFindByLogin(tCtx, req.Login)
}

func (uc *registerUsecase) CreateUser(ctx context.Context, login string, password string) (*domain.User, error) {
	tCtx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	return uc.repo.UserCreate(tCtx, login, password)
}

func (uc *registerUsecase) CreateAccessToken(user *domain.User, secret entities.Secret, lifetime time.Duration) (string, error) {
	token, err := jwt.CreateJWT(secret, user.Login, lifetime)
	if err != nil {
		return "", err
	}

	return token, nil
}
