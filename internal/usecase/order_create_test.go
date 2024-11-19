package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/storage"
	mock_storage "github.com/ex0rcist/gophermart/internal/storage/mocks"
	mock_repository "github.com/ex0rcist/gophermart/internal/storage/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrderCreateUsecase_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIOrderRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx := context.Background()

	user := &domain.User{ID: 1}
	orderNumber := "12345678903" // валидный номер Luhn
	order := &domain.Order{UserID: user.ID, Number: orderNumber, Status: domain.OrderStatusNew}

	mockRepo.EXPECT().OrderFindByNumber(gomock.Any(), orderNumber).Return(nil, ErrOrderNotFound)
	mockRepo.EXPECT().OrderCreate(gomock.Any(), gomock.Any()).Return(order, nil)

	uc := NewOrderCreateUsecase(mockStorage, mockRepo, 5*time.Second)
	result, err := uc.Create(ctx, user, orderNumber)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, orderNumber, result.Number)
	assert.Equal(t, domain.OrderStatusNew, result.Status)
}

func TestOrderCreateUsecase_Create_InvalidOrderNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIOrderRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)

	ctx := context.Background()
	user := &domain.User{ID: 1}
	invalidOrderNumber := "1234567890" // невалидный номер Luhn

	uc := NewOrderCreateUsecase(mockStorage, mockRepo, 5*time.Second)

	result, err := uc.Create(ctx, user, invalidOrderNumber)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrInvalidOrderNumber, err)
}

func TestOrderCreateUsecase_Create_OrderAlreadyRegistered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIOrderRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx := context.Background()
	user := &domain.User{ID: 1}
	orderNumber := "12345678903"
	existingOrder := &domain.Order{UserID: user.ID, Number: orderNumber, Status: domain.OrderStatusNew}

	mockRepo.EXPECT().OrderFindByNumber(gomock.Any(), orderNumber).Return(existingOrder, nil)

	uc := NewOrderCreateUsecase(mockStorage, mockRepo, 5*time.Second)

	result, err := uc.Create(ctx, user, orderNumber)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrOrderAlreadyRegistered, err)
}

func TestOrderCreateUsecase_Create_OrderConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIOrderRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx := context.Background()
	user := &domain.User{ID: 1}
	orderNumber := "12345678903"
	existingOrder := &domain.Order{UserID: 2, Number: orderNumber, Status: domain.OrderStatusNew} // другой пользователь

	mockRepo.EXPECT().OrderFindByNumber(gomock.Any(), orderNumber).Return(existingOrder, nil)

	uc := NewOrderCreateUsecase(mockStorage, mockRepo, 5*time.Second)

	result, err := uc.Create(ctx, user, orderNumber)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrOrderConflict, err)
}

func TestOrderCreateUsecase_OrderFindByNumber_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIOrderRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx := context.Background()
	orderNumber := "12345678903"
	expectedOrder := &domain.Order{Number: orderNumber}

	mockRepo.EXPECT().OrderFindByNumber(gomock.Any(), orderNumber).Return(expectedOrder, nil)

	uc := NewOrderCreateUsecase(mockStorage, mockRepo, 5*time.Second)

	result, err := uc.OrderFindByNumber(ctx, orderNumber)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedOrder.Number, result.Number)
}

func TestOrderCreateUsecase_OrderFindByNumber_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIOrderRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	ctx := context.Background()
	orderNumber := "12345678903"

	mockRepo.EXPECT().OrderFindByNumber(gomock.Any(), orderNumber).Return(nil, storage.ErrRecordNotFound)

	uc := NewOrderCreateUsecase(mockStorage, mockRepo, 5*time.Second)

	result, err := uc.OrderFindByNumber(ctx, orderNumber)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrOrderNotFound, err)
}
