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

func TestRegisterUsecase_Call_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	secret := entities.Secret("supersecret")
	ctx := context.Background()

	registerRequest := RegisterRequest{Login: "newuser", Password: "password"}
	hashedPassword, _ := utils.HashPassword(registerRequest.Password)
	user := &domain.User{
		Login:    "newuser",
		Password: hashedPassword,
	}

	mockRepo.EXPECT().UserFindByLogin(gomock.Any(), "newuser").Return(nil, storage.ErrRecordNotFound)
	mockRepo.EXPECT().UserCreate(gomock.Any(), "newuser", gomock.Any()).Return(user, nil)

	uc := NewRegisterUsecase(mockStorage, mockRepo, secret, 5*time.Second)

	token, err := uc.Call(ctx, registerRequest)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestRegisterUsecase_Call_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	secret := entities.Secret("supersecret")
	ctx := context.Background()

	registerRequest := RegisterRequest{Login: "existinguser", Password: "password"}
	existingUser := &domain.User{Login: "existinguser"}

	mockRepo.EXPECT().UserFindByLogin(gomock.Any(), "existinguser").Return(existingUser, nil)

	uc := NewRegisterUsecase(mockStorage, mockRepo, secret, 5*time.Second)

	token, err := uc.Call(ctx, registerRequest)

	assert.ErrorIs(t, err, ErrUserAlreadyExists)
	assert.Empty(t, token)
}

func TestRegisterUsecase_GetUserByLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	secret := entities.Secret("supersecret")
	ctx := context.Background()

	registerRequest := RegisterRequest{Login: "testuser"}
	user := &domain.User{Login: "testuser"}
	mockRepo.EXPECT().UserFindByLogin(gomock.Any(), "testuser").Return(user, nil)

	uc := NewRegisterUsecase(mockStorage, mockRepo, secret, 5*time.Second)

	resultUser, err := uc.GetUserByLogin(ctx, registerRequest)

	assert.NoError(t, err)
	assert.Equal(t, user, resultUser)
}

func TestRegisterUsecase_GetUserByLogin_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	secret := entities.Secret("supersecret")
	ctx := context.Background()

	registerRequest := RegisterRequest{Login: "nonexistent"}

	mockRepo.EXPECT().UserFindByLogin(gomock.Any(), "nonexistent").Return(nil, storage.ErrRecordNotFound)

	uc := NewRegisterUsecase(mockStorage, mockRepo, secret, 5*time.Second)

	resultUser, err := uc.GetUserByLogin(ctx, registerRequest)

	assert.ErrorIs(t, err, storage.ErrRecordNotFound)
	assert.Nil(t, resultUser)
}

func TestRegisterUsecase_CreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	secret := entities.Secret("supersecret")
	ctx := context.Background()

	login := "newuser"
	password := "password"
	user := &domain.User{Login: "newuser", Password: password}

	mockRepo.EXPECT().UserCreate(gomock.Any(), login, password).Return(user, nil)

	uc := NewRegisterUsecase(mockStorage, mockRepo, secret, 5*time.Second)

	newUser, err := uc.CreateUser(ctx, login, password)

	assert.NoError(t, err)
	assert.Equal(t, user, newUser)
}

func TestRegisterUsecase_CreateAccessToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &domain.User{
		Login: "testuser",
	}
	secret := entities.Secret("supersecret")
	lifetime := 24 * time.Hour

	uc := NewRegisterUsecase(nil, nil, secret, 5*time.Second)

	token, err := uc.CreateAccessToken(user, secret, lifetime)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
