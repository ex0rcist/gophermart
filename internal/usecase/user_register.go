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

var ErrUserAlreadyExists = errors.New("login already exists")

type IRegisterUsecase interface {
	Call(ctx context.Context, form RegisterRequest) (string, error)
	GetUserByLogin(ctx context.Context, form RegisterRequest) (*domain.User, error)
	CreateUser(ctx context.Context, login string, password string) (*domain.User, error)
	CreateAccessToken(user *domain.User, secret entities.Secret, lifetime time.Duration) (string, error)
}

type registerUsecase struct {
	storage        storage.IPGXStorage
	repo           repository.IUserRepository
	secret         entities.Secret
	contextTimeout time.Duration
}
type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=3"`
}

func NewRegisterUsecase(storage storage.IPGXStorage, repo repository.IUserRepository, secret entities.Secret, timeout time.Duration) IRegisterUsecase {
	return &registerUsecase{storage: storage, repo: repo, secret: secret, contextTimeout: timeout}
}

func (uc *registerUsecase) Call(ctx context.Context, form RegisterRequest) (string, error) {
	// находим пользователя
	existingUser, err := uc.GetUserByLogin(ctx, form)
	if err != nil && err != storage.ErrRecordNotFound {
		return "", err
	}
	if existingUser != nil {
		return "", ErrUserAlreadyExists
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
	token, err := uc.CreateAccessToken(newUser, uc.secret, jwt.LoginTokenLifetime)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *registerUsecase) GetUserByLogin(ctx context.Context, req RegisterRequest) (*domain.User, error) {
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
