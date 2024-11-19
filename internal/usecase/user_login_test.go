package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/storage"
	mock_storage "github.com/ex0rcist/gophermart/internal/storage/mocks"
	mock_repository "github.com/ex0rcist/gophermart/internal/storage/repository/mocks"
	"github.com/ex0rcist/gophermart/internal/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLoginUsecase_Call_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	secret := entities.Secret("supersecret")
	ctx := context.Background()
	loginRequest := LoginRequest{
		Login:    "testuser",
		Password: "password",
	}

	p, _ := utils.HashPassword("password")
	user := &domain.User{
		Login:    "testuser",
		Password: p,
	}

	mockRepo.EXPECT().UserFindByLogin(gomock.Any(), "testuser").Return(user, nil)

	uc := NewLoginUsecase(mockStorage, mockRepo, secret, 5*time.Second)

	token, err := uc.Call(ctx, loginRequest)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLoginUsecase_Call_InvalidLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	secret := entities.Secret("supersecret")
	ctx := context.Background()
	loginRequest := LoginRequest{
		Login:    "wronguser",
		Password: "password",
	}

	mockRepo.EXPECT().UserFindByLogin(gomock.Any(), "wronguser").Return(nil, storage.ErrRecordNotFound)

	uc := NewLoginUsecase(mockStorage, mockRepo, secret, 5*time.Second)

	token, err := uc.Call(ctx, loginRequest)

	assert.ErrorIs(t, err, ErrInvalidLoginOrPassword)
	assert.Empty(t, token)
}

func TestLoginUsecase_GetUserByLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	secret := entities.Secret("supersecret")
	ctx := context.Background()
	loginRequest := LoginRequest{
		Login: "testuser",
	}

	p, _ := utils.HashPassword("password")
	user := &domain.User{
		Login:    "testuser",
		Password: p,
	}

	mockRepo.EXPECT().UserFindByLogin(gomock.Any(), "testuser").Return(user, nil)

	uc := NewLoginUsecase(mockStorage, mockRepo, secret, 5*time.Second)

	resultUser, err := uc.GetUserByLogin(ctx, loginRequest)

	assert.NoError(t, err)
	assert.Equal(t, user, resultUser)
}

func TestLoginUsecase_GetUserByLogin_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	secret := entities.Secret("supersecret")
	ctx := context.Background()
	loginRequest := LoginRequest{
		Login: "nonexistent",
	}

	mockRepo.EXPECT().UserFindByLogin(gomock.Any(), "nonexistent").Return(nil, storage.ErrRecordNotFound)

	uc := NewLoginUsecase(mockStorage, mockRepo, secret, 5*time.Second)

	resultUser, err := uc.GetUserByLogin(ctx, loginRequest)

	assert.ErrorIs(t, err, ErrInvalidLoginOrPassword)
	assert.Nil(t, resultUser)
}

func TestLoginUsecase_ComparePassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p, _ := utils.HashPassword("password")
	user := &domain.User{
		Login:    "testuser",
		Password: p,
	}
	password := "password"

	uc := NewLoginUsecase(nil, nil, "", 5*time.Second)

	err := uc.ComparePassword(user, password)

	assert.NoError(t, err)
}

func TestLoginUsecase_ComparePassword_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p, _ := utils.HashPassword("password")
	user := &domain.User{
		Login:    "testuser",
		Password: p,
	}
	wrongPassword := "wrongpassword"

	uc := NewLoginUsecase(nil, nil, "", 5*time.Second)

	err := uc.ComparePassword(user, wrongPassword)

	assert.Error(t, err)
}

func TestLoginUsecase_CreateAccessToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &domain.User{
		Login: "testuser",
	}
	secret := entities.Secret("supersecret")
	lifetime := 24 * time.Hour

	uc := NewLoginUsecase(nil, nil, secret, 5*time.Second)

	token, err := uc.CreateAccessToken(user, secret, lifetime)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
