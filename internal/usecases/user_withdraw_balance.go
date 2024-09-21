package usecases

import (
	"context"
	"errors"

	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/models"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/shopspring/decimal"
)

type UserWithdrawBalanceForm struct {
	UserID      models.UserID
	OrderNumber string          `json:"order" binding:"required,luhn"`
	Amount      decimal.Decimal `json:"sum" binding:"required"`
}

var (
	ErrInsufficientUserBalance = errors.New("insufficient user balance")
)

func UserWithdrawBalance(ctx context.Context, s *storage.PGXStorage, form UserWithdrawBalanceForm) error {
	// стартуем транзакцию, детали реализации скрыты в PGXTx
	tx, err := s.StartTx(ctx)
	if err != nil {
		logging.LogErrorCtx(ctx, err, "UserWithdrawBalance(): error starting tx")
		return err
	}
	defer func() {
		err := tx.Rollback()
		if err != nil {
			logging.LogErrorCtx(ctx, err, "UserWithdrawBalance(): error rolling tx back")
		}
	}()

	// получаем активный баланс, транзакция блокирует user.balance и user.withdrawn
	res, err := s.UserGetBalance(ctx, tx.Tx, form.UserID)
	if err != nil {
		return err
	}

	// убеждаемся что баланса достаточно
	if res.Balance.Cmp(form.Amount) == -1 {
		return ErrInsufficientUserBalance
	}

	// создаем списание
	data := storage.CreateWithdrawalDTO{UserID: form.UserID, OrderNumber: form.OrderNumber, Amount: form.Amount}
	err = s.WithdrawalCreate(ctx, tx.Tx, data)
	if err != nil {
		logging.LogErrorCtx(ctx, err, "UserWithdrawBalance(): error creating withdrawal")
		return err
	}

	// актуализируем user.balance и user.withdrawn
	err = s.UserUpdateBalanceAndWithdrawals(ctx, tx.Tx, form.UserID)
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
