package usecase

import (
	"context"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/storage"
)

type withdrawBalanceUsecase struct {
	storage        storage.IPGXStorage
	userRepo       domain.IUserRepository
	wdrwRepo       domain.IWithdrawalRepository
	contextTimeout time.Duration
}

func NewWithdrawBalanceUsecase(userRepo domain.IUserRepository, wdrwRepo domain.IWithdrawalRepository, timeout time.Duration) domain.IWithdrawBalanceUsecase {
	return &withdrawBalanceUsecase{userRepo: userRepo, wdrwRepo: wdrwRepo, contextTimeout: timeout}
}

func (uc *withdrawBalanceUsecase) Call(c context.Context, user *domain.User, form domain.WithdrawBalanceRequest) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// стартуем транзакцию
	tx, err := uc.storage.StartTx(ctx)
	if err != nil {
		logging.LogErrorCtx(ctx, err, "withdrawBalanceUsecase(): error starting tx")
		return err
	}
	defer func() {
		err := tx.Rollback()
		if err != nil {
			logging.LogErrorCtx(ctx, err, "withdrawBalanceUsecase(): error rolling tx back")
		}
	}()

	// получаем активный баланс, транзакция блокирует user.balance и user.withdrawn
	b, _, err := uc.userRepo.UserGetBalance(ctx, tx.Tx, user.ID)
	if err != nil {
		return err
	}

	// убеждаемся что баланса достаточно
	if b.Cmp(form.Amount) == -1 {
		return domain.ErrInsufficientUserBalance
	}

	// создаем списание
	err = uc.wdrwRepo.WithdrawalCreate(ctx, tx.Tx, domain.Withdrawal{UserID: user.ID, OrderNumber: form.OrderNumber, Amount: form.Amount})
	if err != nil {
		logging.LogErrorCtx(ctx, err, "UserWithdrawBalance(): error creating withdrawal")
		return err
	}

	// актуализируем user.balance и user.withdrawn
	err = uc.userRepo.UserUpdateBalanceAndWithdrawals(ctx, tx.Tx, user.ID)
	if err != nil {
		logging.LogErrorCtx(ctx, err, "UserWithdrawBalance(): error recalculating balance/withdrawn")
		return err
	}

	// завершаем транзакцию
	err = tx.Commit()
	if err != nil {
		logging.LogErrorCtx(ctx, err, "UserWithdrawBalance(): error commiting tx")
		return err
	}

	return nil
}
