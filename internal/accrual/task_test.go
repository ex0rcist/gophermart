package accrual_test

import (
	"context"
	"testing"

	"github.com/ex0rcist/gophermart/internal/accrual"
	mock_accrual "github.com/ex0rcist/gophermart/internal/accrual/mocks"
	"github.com/ex0rcist/gophermart/internal/config"
	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/mock/gomock"

	mock_storage "github.com/ex0rcist/gophermart/internal/storage/mocks"
	mock_repository "github.com/ex0rcist/gophermart/internal/storage/repository/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTask_Handle_StatusRegistered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	order := &domain.Order{ID: 1, Number: "12345", Status: domain.OrderStatusNew}

	mockClient := mock_accrual.NewMockIClient(ctrl)
	mockOrderRepo := mock_repository.NewMockIOrderRepository(ctrl)
	mockUserRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	mockPool := mock_storage.NewMockIPGXPool(ctrl)

	mockStorage.EXPECT().GetPool().Return(mockPool).AnyTimes()

	mockClient.EXPECT().
		GetBonuses(gomock.Any(), "12345").
		Return(&accrual.Response{
			OrderNumber: "12345",
			Status:      accrual.StatusRegistered,
			Amount:      decimal.Zero,
		}, nil)

	mockOrderRepo.EXPECT().OrderUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockUserRepo.EXPECT().UserUpdateBalanceAndWithdrawals(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	cfg, _ := config.NewDefault(&config.Config{})
	service := accrual.NewService(
		context.Background(),
		&cfg.Accrual,
		mockClient,
		mockStorage,
		mockUserRepo,
		mockOrderRepo,
	)

	task := accrual.NewTask(service, order)
	err := task.Handle()
	assert.NoError(t, err)
}

func TestTask_Handle_StatusProcessing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Моки зависимостей
	mockClient := mock_accrual.NewMockIClient(ctrl)
	mockOrderRepo := mock_repository.NewMockIOrderRepository(ctrl)
	mockUserRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	mockPool := mock_storage.NewMockIPGXPool(ctrl)
	mockTx := new(storage.PGXTxMock)

	mockStorage.EXPECT().GetPool().Return(mockPool).AnyTimes()

	// Тестовый заказ
	order := &domain.Order{ID: 1, Number: "12345", Status: domain.OrderStatusNew}

	// Настройка клиента
	mockClient.EXPECT().
		GetBonuses(gomock.Any(), "12345").
		Return(&accrual.Response{
			OrderNumber: "12345",
			Status:      accrual.StatusProcessing,
			Amount:      decimal.Zero,
		}, nil)

	// Настройка транзакции
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)
	mockTx.
		On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, nil)

	mockOrderRepo.EXPECT().OrderUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
	mockUserRepo.EXPECT().UserUpdateBalanceAndWithdrawals(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	cfg, _ := config.NewDefault(&config.Config{})
	service := accrual.NewService(
		context.Background(),
		&cfg.Accrual,
		mockClient,
		mockStorage,
		mockUserRepo,
		mockOrderRepo,
	)

	task := accrual.NewTask(service, order)
	err := task.Handle()
	assert.NoError(t, err)
}

func TestTask_Handle_StatusInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Моки зависимостей
	mockClient := mock_accrual.NewMockIClient(ctrl)
	mockOrderRepo := mock_repository.NewMockIOrderRepository(ctrl)
	mockUserRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	mockPool := mock_storage.NewMockIPGXPool(ctrl)
	mockTx := new(storage.PGXTxMock)

	// Тестовый заказ
	order := &domain.Order{ID: 1, Number: "12345", Status: domain.OrderStatusNew}

	mockStorage.EXPECT().GetPool().Return(mockPool).AnyTimes()
	mockClient.EXPECT().
		GetBonuses(gomock.Any(), "12345").
		Return(&accrual.Response{
			OrderNumber: "12345",
			Status:      accrual.StatusInvalid,
			Amount:      decimal.Zero,
		}, nil)

	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)
	mockTx.
		On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, nil)

	mockOrderRepo.EXPECT().OrderUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
	mockUserRepo.EXPECT().UserUpdateBalanceAndWithdrawals(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	cfg, _ := config.NewDefault(&config.Config{})

	service := accrual.NewService(
		context.Background(),
		&cfg.Accrual,
		mockClient,
		mockStorage,
		mockUserRepo,
		mockOrderRepo,
	)

	task := accrual.NewTask(service, order)
	err := task.Handle()
	assert.NoError(t, err)
}

func TestTask_Handle_StatusProcessed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Моки зависимостей
	mockClient := mock_accrual.NewMockIClient(ctrl)
	mockOrderRepo := mock_repository.NewMockIOrderRepository(ctrl)
	mockUserRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	mockPool := mock_storage.NewMockIPGXPool(ctrl)
	mockTx := new(storage.PGXTxMock)

	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockStorage.EXPECT().GetPool().Return(mockPool).AnyTimes()

	// Тестовый заказ
	order := &domain.Order{ID: 1, Number: "12345", Status: domain.OrderStatusNew}

	mockClient.EXPECT().
		GetBonuses(gomock.Any(), "12345").
		Return(&accrual.Response{
			OrderNumber: "12345",
			Status:      accrual.StatusProcessed,
			Amount:      decimal.NewFromFloat(150.50),
		}, nil)

	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)
	mockTx.
		On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, nil)

	mockOrderRepo.EXPECT().OrderUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockUserRepo.EXPECT().UserUpdateBalanceAndWithdrawals(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	cfg, _ := config.NewDefault(&config.Config{})

	service := accrual.NewService(
		context.Background(),
		&cfg.Accrual,
		mockClient,
		mockStorage,
		mockUserRepo,
		mockOrderRepo,
	)

	task := accrual.NewTask(service, order)

	err := task.Handle()
	assert.NoError(t, err)
}
