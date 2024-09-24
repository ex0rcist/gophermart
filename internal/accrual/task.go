package accrual

import (
	"context"
	"errors"
	"fmt"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

type ITask interface {
	Handle() error
}

type Task struct {
	service *Service
	order   *domain.Order
}

func NewTask(service *Service, order *domain.Order) Task {
	return Task{
		service: service,
		order:   order,
	}
}

func (t Task) Handle() error {
	tCtx, cancel := context.WithTimeout(context.Background(), t.service.contextTimeout)
	defer cancel()

	// внедряем общую метку в логи запросов и логи сервиса, для облегчения чтения
	ctx := setupCtxWithRID(tCtx)

	// получаем статус и баланс из accrual
	res, err := t.service.client.GetBonuses(ctx, t.order.Number)
	if err != nil {
		return err
	}

	// проверяем
	switch res.Status {
	case StatusRegistered:
		// только что создан
		logging.LogInfoCtx(ctx, fmt.Sprintf("%s is just created, do nothing", t.order))

	case StatusProcessing:
		// в обработке; если статус в базе не совпадает, обновляем
		logging.LogInfoCtx(ctx, fmt.Sprintf("%s is still in processing", t.order))
		if t.order.Status != domain.OrderStatusProcessing {
			err := t.updateOrder(ctx, domain.OrderStatusProcessing, decimal.NewFromInt(0))
			if err != nil {
				return err
			}
		}

	case StatusInvalid:
		// invalid; обновляем статус
		logging.LogInfoCtx(ctx, fmt.Sprintf("%s is invalid", t.order))
		err := t.updateOrder(ctx, domain.OrderStatusInvalid, decimal.NewFromInt(0))
		if err != nil {
			return err
		}

	case StatusProcessed:
		// обработан; обновляем статус и сумму накоплений
		logging.LogInfoCtx(ctx, fmt.Sprintf("%s processed, accrual=%s", t.order, res.Amount))
		err := t.updateOrder(ctx, domain.OrderStatusProcessed, res.Amount)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t Task) updateOrder(ctx context.Context, status domain.OrderStatus, amount decimal.Decimal) error {
	tx, err := t.service.storage.GetPool().Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			logging.LogErrorCtx(ctx, err, "Task: updateOrder(): error rolling tx back")
		}
	}()

	err = t.service.orderRepo.OrderUpdate(ctx, tx, domain.Order{ID: t.order.ID, Status: status, Accrual: amount})
	if err != nil {
		return err
	}

	if status == domain.OrderStatusProcessed {
		err = t.service.userRepo.UserUpdateBalanceAndWithdrawals(ctx, tx, t.order.UserID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		logging.LogErrorCtx(ctx, err, "Task: updateOrder(): error commiting tx")
		return err
	}

	return nil
}

func setupCtxWithRID(ctx context.Context) context.Context {
	logger := log.Logger.With().Ctx(ctx).Str("rid", utils.GenerateRequestID()).Logger()
	return logger.WithContext(ctx)
}
