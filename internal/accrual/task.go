package accrual

import (
	"context"
	"fmt"

	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/models"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/shopspring/decimal"
)

type Task struct {
	service *Service
	order   *models.Order
}

func NewTask(service *Service, order *models.Order) Task {
	return Task{
		service: service,
		order:   order,
	}
}

func (t Task) Handle() error {
	// внедряем общую метку в логи запросов и логи сервиса, для облегчения чтения
	ctx := setupCtxWithRID(context.Background())

	// получаем статус и баланс из accrual
	res, err := t.service.client.GetBonuses(ctx, t.order.Number)
	if err != nil {
		return err
	}

	// проверяем
	switch res.Status {
	case StatusNew:
		// только что создан
		logging.LogInfoCtx(ctx, fmt.Sprintf("%s is just created, do nothing", t.order))

	case StatusProcessing:
		// в обработке; если статус в базе не совпадает, обновляем
		logging.LogInfoCtx(ctx, fmt.Sprintf("%s is still in processing", t.order))
		if t.order.Status != models.OrderStatusProcessing {
			err := t.updateOrder(ctx, models.OrderStatusProcessing, decimal.NewFromInt(0))
			if err != nil {
				return err
			}
		}

	case StatusInvalid:
		// invalid; обновляем статус
		logging.LogInfoCtx(ctx, fmt.Sprintf("%s is invalid", t.order))
		err := t.updateOrder(ctx, models.OrderStatusInvalid, decimal.NewFromInt(0))
		if err != nil {
			return err
		}

	case StatusProcessed:
		// обработан; обновляем статус и сумму накоплений
		logging.LogInfoCtx(ctx, fmt.Sprintf("%s processed, accrual=%s", t.order, res.Amount))
		err := t.updateOrder(ctx, models.OrderStatusProcessed, res.Amount)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t Task) updateOrder(ctx context.Context, status models.OrderStatus, amount decimal.Decimal) error {
	d := storage.OrderUpdateDTO{ID: t.order.ID, Status: status, Accrual: amount}

	tx, err := t.service.storage.StartTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err := t.service.storage.RollbackTx(ctx, tx)
		if err != nil {
			logging.LogErrorCtx(ctx, err, "Task: updateOrder(): error rolling tx back")
		}
	}()

	err = t.service.storage.OrderUpdate(ctx, tx.Tx, d)
	if err != nil {
		return err
	}

	err = t.service.storage.UserUpdateBalanceAndWithdrawals(ctx, tx.Tx, t.order.UserID)
	if err != nil {
		return err
	}

	err = t.service.storage.CommitTx(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}
