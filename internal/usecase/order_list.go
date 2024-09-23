package usecase

import (
	"context"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
)

type orderListUsecase struct {
	repo           domain.IOrderRepository
	contextTimeout time.Duration
}

func NewOrderListUsecase(repo domain.IOrderRepository, timeout time.Duration) domain.IOrderListUsecase {
	return &orderListUsecase{repo: repo, contextTimeout: timeout}
}

func (uc *orderListUsecase) Call(ctx context.Context, user *domain.User) ([]*domain.OrderListResult, error) {
	orders, err := uc.repo.OrderList(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.OrderListResult, 0)
	for _, o := range orders {
		el := domain.OrderListResult{Number: o.Number, Status: o.Status, CreatedAt: entities.RFC3339Time(o.CreatedAt)}

		if o.Status == domain.OrderStatusProcessed {
			val := entities.GDecimal(o.Accrual)
			el.Accrual = &val
		}

		result = append(result, &el)
	}

	return result, nil
}
