package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/storage"
	mock_storage "github.com/ex0rcist/gophermart/internal/storage/mocks"
	mock_repository "github.com/ex0rcist/gophermart/internal/storage/repository/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/mock/gomock"
)

func TestWithdrawBalanceUsecase_Call_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockWdrwRepo := mock_repository.NewMockIWithdrawalRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	mockPool := mock_storage.NewMockIPGXPool(ctrl)
	mockTx := new(storage.PGXTxMock)

	ctx := context.Background()
	user := &domain.User{ID: 1}
	form := WithdrawBalanceRequest{
		OrderNumber: "12345678903", // валидный номер Luhn
		Amount:      decimal.NewFromFloat(100),
	}

	balance := decimal.NewFromFloat(200) // баланс больше чем сумма списания

	mockStorage.EXPECT().GetPool().Return(mockPool).AnyTimes()
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockUserRepo.EXPECT().UserGetBalance(gomock.Any(), mockTx, user.ID).Return(&balance, nil, nil)
	mockWdrwRepo.EXPECT().WithdrawalCreate(gomock.Any(), mockTx, gomock.Any()).Return(nil)
	mockUserRepo.EXPECT().UserUpdateBalanceAndWithdrawals(gomock.Any(), mockTx, user.ID).Return(nil)
	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)

	uc := NewWithdrawBalanceUsecase(mockStorage, mockUserRepo, mockWdrwRepo, 5*time.Second)

	err := uc.Call(ctx, user, form)

	assert.NoError(t, err)
}

func TestWithdrawBalanceUsecase_Call_InvalidOrderNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockWdrwRepo := mock_repository.NewMockIWithdrawalRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)

	ctx := context.Background()
	user := &domain.User{ID: 1}
	invalidForm := WithdrawBalanceRequest{
		OrderNumber: "invalid", // невалидный номер Luhn
		Amount:      decimal.NewFromFloat(100),
	}

	uc := NewWithdrawBalanceUsecase(mockStorage, mockUserRepo, mockWdrwRepo, 5*time.Second)

	err := uc.Call(ctx, user, invalidForm)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidOrderNumber, err)
}

func TestWithdrawBalanceUsecase_Call_InsufficientBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockWdrwRepo := mock_repository.NewMockIWithdrawalRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	mockPool := mock_storage.NewMockIPGXPool(ctrl)
	mockTx := new(storage.PGXTxMock)

	ctx := context.Background()
	user := &domain.User{ID: 1}
	form := WithdrawBalanceRequest{
		OrderNumber: "12345678903", // валидный номер Luhn
		Amount:      decimal.NewFromFloat(300.0),
	}

	balance := decimal.NewFromFloat(100.0) // баланс меньше чем сумма списания

	mockStorage.EXPECT().GetPool().Return(mockPool).AnyTimes()
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockUserRepo.EXPECT().UserGetBalance(gomock.Any(), mockTx, user.ID).Return(&balance, nil, nil)

	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)

	uc := NewWithdrawBalanceUsecase(mockStorage, mockUserRepo, mockWdrwRepo, 5*time.Second)

	err := uc.Call(ctx, user, form)

	assert.Error(t, err)
	assert.Equal(t, ErrInsufficientUserBalance, err)
}

func TestWithdrawBalanceUsecase_Call_TransactionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockWdrwRepo := mock_repository.NewMockIWithdrawalRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	mockPool := mock_storage.NewMockIPGXPool(ctrl)

	ctx := context.Background()
	user := &domain.User{ID: 1}
	form := WithdrawBalanceRequest{
		OrderNumber: "12345678903", // валидный номер Luhn
		Amount:      decimal.NewFromFloat(100.0),
	}

	expectedError := errors.New("transaction error")

	mockStorage.EXPECT().GetPool().Return(mockPool).AnyTimes()
	mockPool.EXPECT().Begin(gomock.Any()).Return(nil, expectedError)

	uc := NewWithdrawBalanceUsecase(mockStorage, mockUserRepo, mockWdrwRepo, 5*time.Second)

	err := uc.Call(ctx, user, form)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestWithdrawBalanceUsecase_Call_CommitError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockIUserRepository(ctrl)
	mockWdrwRepo := mock_repository.NewMockIWithdrawalRepository(ctrl)
	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	mockPool := mock_storage.NewMockIPGXPool(ctrl)
	mockTx := new(storage.PGXTxMock)

	ctx := context.Background()
	user := &domain.User{ID: 1}
	form := WithdrawBalanceRequest{
		OrderNumber: "12345678903", // валидный номер Luhn
		Amount:      decimal.NewFromFloat(100.0),
	}

	balance := decimal.NewFromFloat(200.0) // достаточно средств

	commitError := errors.New("commit error")

	mockStorage.EXPECT().GetPool().Return(mockPool).AnyTimes()
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockUserRepo.EXPECT().UserGetBalance(gomock.Any(), mockTx, user.ID).Return(&balance, nil, nil)
	mockWdrwRepo.EXPECT().WithdrawalCreate(gomock.Any(), mockTx, gomock.Any()).Return(nil)
	mockUserRepo.EXPECT().UserUpdateBalanceAndWithdrawals(gomock.Any(), mockTx, user.ID).Return(nil)

	mockTx.On("Commit", mock.Anything).Return(commitError)
	mockTx.On("Rollback", mock.Anything).Return(nil)

	uc := NewWithdrawBalanceUsecase(mockStorage, mockUserRepo, mockWdrwRepo, 5*time.Second)

	err := uc.Call(ctx, user, form)

	assert.Error(t, err)
	assert.Equal(t, commitError, err)
}
