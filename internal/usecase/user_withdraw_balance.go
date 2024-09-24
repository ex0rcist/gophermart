package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/storage/repository"
	"github.com/ex0rcist/gophermart/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

var ErrInsufficientUserBalance = errors.New("insufficient user balance")

type WithdrawBalanceRequest struct {
	OrderNumber string          `json:"order" binding:"required,luhn"`
	Amount      decimal.Decimal `json:"sum" binding:"required"`
}

type IWithdrawBalanceUsecase interface {
	Call(ctx context.Context, user *domain.User, req WithdrawBalanceRequest) error
}

type withdrawBalanceUsecase struct {
	storage        storage.IPGXStorage
	userRepo       repository.IUserRepository
	wdrwRepo       repository.IWithdrawalRepository
	contextTimeout time.Duration
}

func NewWithdrawBalanceUsecase(
	storage storage.IPGXStorage,
	userRepo repository.IUserRepository,
	wdrwRepo repository.IWithdrawalRepository,
	timeout time.Duration,
) IWithdrawBalanceUsecase {
	return &withdrawBalanceUsecase{storage: storage, userRepo: userRepo, wdrwRepo: wdrwRepo, contextTimeout: timeout}
}

func (uc *withdrawBalanceUsecase) Call(ctx context.Context, user *domain.User, form WithdrawBalanceRequest) error {
	tCtx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	// валидируем номер заказа
	if !utils.LuhnCheck(form.OrderNumber) {
		return ErrInvalidOrderNumber
	}

	// стартуем транзакцию
	tx, err := uc.storage.GetPool().Begin(ctx)
	if err != nil {
		logging.LogErrorCtx(ctx, err, "withdrawBalanceUsecase(): error starting tx")
		return err
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(pgx.ErrTxClosed, err) {
			logging.LogErrorCtx(ctx, err, "withdrawBalanceUsecase(): error rolling tx back")
		}
	}()

	// получаем активный баланс, транзакция блокирует user.balance и user.withdrawn
	b, _, err := uc.userRepo.UserGetBalance(tCtx, tx, user.ID)
	if err != nil {
		return err
	}

	// убеждаемся что баланса достаточно
	if b.Cmp(form.Amount) == -1 {
		return ErrInsufficientUserBalance
	}

	// создаем списание
	err = uc.wdrwRepo.WithdrawalCreate(tCtx, tx, domain.Withdrawal{UserID: user.ID, OrderNumber: form.OrderNumber, Amount: form.Amount})
	if err != nil {
		logging.LogErrorCtx(ctx, err, "UserWithdrawBalance(): error creating withdrawal")
		return err
	}

	// актуализируем user.balance и user.withdrawn
	err = uc.userRepo.UserUpdateBalanceAndWithdrawals(ctx, tx, user.ID)
	if err != nil {
		logging.LogErrorCtx(ctx, err, "UserWithdrawBalance(): error recalculating balance/withdrawn")
		return err
	}

	// завершаем транзакцию
	err = tx.Commit(ctx)
	if err != nil {
		logging.LogErrorCtx(ctx, err, "UserWithdrawBalance(): error commiting tx")
		return err
	}

	return nil
}
