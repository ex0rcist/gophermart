package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	mock_storage "github.com/ex0rcist/gophermart/internal/storage/mocks"
	mock_repository "github.com/ex0rcist/gophermart/internal/storage/repository/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrderListUsecase_Call_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIOrderRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx := context.Background()
	user := &domain.User{
		ID: 1,
	}

	orders := []*domain.Order{
		{
			Number:    "123456",
			Status:    domain.OrderStatusProcessed,
			CreatedAt: time.Now(),
			Accrual:   decimal.NewFromFloat(10.0),
		},
		{
			Number:    "654321",
			Status:    domain.OrderStatusNew,
			CreatedAt: time.Now(),
		},
	}

	mockRepo.EXPECT().OrderList(gomock.Any(), user.ID).Return(orders, nil)
	uc := NewOrderListUsecase(mockStorage, mockRepo, 5*time.Second)
	result, err := uc.Call(ctx, user)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "123456", result[0].Number)
	assert.NotNil(t, result[0].Accrual)
	assert.Equal(t, "654321", result[1].Number)
	assert.Nil(t, result[1].Accrual)
}

func TestOrderListUsecase_Call_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIOrderRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx := context.Background()
	user := &domain.User{
		ID: 1,
	}

	expectedError := errors.New("database error")
	mockRepo.EXPECT().OrderList(gomock.Any(), user.ID).Return(nil, expectedError)
	uc := NewOrderListUsecase(mockStorage, mockRepo, 5*time.Second)

	result, err := uc.Call(ctx, user)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestOrderListUsecase_Call_ContextTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIOrderRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	user := &domain.User{
		ID: 1,
	}

	mockRepo.EXPECT().OrderList(gomock.Any(), user.ID).AnyTimes().DoAndReturn(func(ctx context.Context, userID domain.UserID) ([]*domain.Order, error) {
		time.Sleep(2 * time.Millisecond) // Симуляция задержки, чтобы истек контекст
		return nil, context.DeadlineExceeded
	})
	uc := NewOrderListUsecase(mockStorage, mockRepo, 5*time.Second)
	result, err := uc.Call(ctx, user)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, context.DeadlineExceeded, err)
}
