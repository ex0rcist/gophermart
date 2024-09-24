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

func TestWithdrawalListUsecase_Call_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIWithdrawalRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx := context.Background()
	user := &domain.User{ID: 1}

	withdrawals := []*domain.Withdrawal{
		{
			OrderNumber: "12345678903",
			Amount:      decimal.NewFromFloat(100.0),
			CreatedAt:   time.Now(),
		},
		{
			OrderNumber: "98765432100",
			Amount:      decimal.NewFromFloat(50.0),
			CreatedAt:   time.Now(),
		},
	}

	mockRepo.EXPECT().WithdrawalList(gomock.Any(), user.ID).Return(withdrawals, nil)

	uc := NewWithdrawalListUsecase(mockStorage, mockRepo, 5*time.Second)

	result, err := uc.Call(ctx, user)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "12345678903", result[0].OrderNumber)
	assert.Equal(t, entities.GDecimal(decimal.NewFromFloat(100.0)), result[0].Amount)
	assert.Equal(t, "98765432100", result[1].OrderNumber)
	assert.Equal(t, entities.GDecimal(decimal.NewFromFloat(50.0)), result[1].Amount)
}

func TestWithdrawalListUsecase_Call_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIWithdrawalRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx := context.Background()
	user := &domain.User{
		ID: 1,
	}

	expectedError := errors.New("database error")

	mockRepo.EXPECT().WithdrawalList(gomock.Any(), user.ID).Return(nil, expectedError)

	uc := NewWithdrawalListUsecase(mockStorage, mockRepo, 5*time.Second)

	result, err := uc.Call(ctx, user)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestWithdrawalListUsecase_Call_EmptyList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIWithdrawalRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx := context.Background()
	user := &domain.User{ID: 1}

	mockRepo.EXPECT().WithdrawalList(gomock.Any(), user.ID).Return([]*domain.Withdrawal{}, nil)

	uc := NewWithdrawalListUsecase(mockStorage, mockRepo, 5*time.Second)

	result, err := uc.Call(ctx, user)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestWithdrawalListUsecase_Call_ContextTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIWithdrawalRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	user := &domain.User{ID: 1}

	mockRepo.EXPECT().WithdrawalList(gomock.Any(), user.ID).AnyTimes().DoAndReturn(func(ctx context.Context, userID domain.UserID) ([]*domain.Withdrawal, error) {
		time.Sleep(2 * time.Millisecond)
		return nil, context.DeadlineExceeded
	})

	uc := NewWithdrawalListUsecase(mockStorage, mockRepo, 5*time.Second)

	result, err := uc.Call(ctx, user)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, context.DeadlineExceeded, err)
}
