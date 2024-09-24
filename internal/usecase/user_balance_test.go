package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	mock_storage "github.com/ex0rcist/gophermart/internal/storage/mocks"
	mock_repository "github.com/ex0rcist/gophermart/internal/storage/repository/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserBalanceUsecase_Call_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx := context.Background()
	user := domain.User{ID: 1}

	balance := decimal.NewFromFloat(float64(100.50))
	withdrawn := decimal.NewFromFloat(float64(50.25))

	mockRepo.EXPECT().UserGetBalance(gomock.Any(), nil, gomock.Any()).Return(&balance, &withdrawn, nil)

	uc := NewGetUserBalanceUsecase(mockStorage, mockRepo, 5*time.Second)

	result, err := uc.Call(ctx, &user)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entities.GDecimal(balance), result.Current)
	assert.Equal(t, entities.GDecimal(withdrawn), result.Withdrawn)
}

func TestGetUserBalanceUsecase_Call_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx := context.Background()
	user := domain.User{ID: 1}

	expectedError := errors.New("database error")

	mockRepo.EXPECT().UserGetBalance(gomock.Any(), nil, gomock.Any()).Return(nil, nil, expectedError)

	uc := NewGetUserBalanceUsecase(mockStorage, mockRepo, 5*time.Second)

	result, err := uc.Call(ctx, &user)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestGetUserBalanceUsecase_Call_ContextTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	user := domain.User{ID: 1}

	mockRepo.EXPECT().UserGetBalance(gomock.Any(), nil, gomock.Any()).AnyTimes().DoAndReturn(func(ctx context.Context, tx interface{}, id domain.UserID) (*float64, *float64, error) {
		time.Sleep(2 * time.Millisecond) // симуляция задержки, чтобы истек контекст
		return nil, nil, context.DeadlineExceeded
	})

	uc := NewGetUserBalanceUsecase(mockStorage, mockRepo, 5*time.Second)

	result, err := uc.Call(ctx, &user)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, context.DeadlineExceeded, err)
}
