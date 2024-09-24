package usecase

import (
	"context"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/storage/repository"
)

type IOrderListUsecase interface {
	Call(ctx context.Context, user *domain.User) ([]*OrderListResult, error)
}

type orderListUsecase struct {
	storage        storage.IPGXStorage
	repo           repository.IOrderRepository
	contextTimeout time.Duration
}

type OrderListResult struct {
	Number    string               `json:"number"`
	Status    domain.OrderStatus   `json:"status"`
	Accrual   *entities.GDecimal   `json:"accrual,omitempty"` // без использования указателя omitempty не считает значение пустым
	CreatedAt entities.RFC3339Time `json:"uploaded_at"`
}

func NewOrderListUsecase(storage storage.IPGXStorage, repo repository.IOrderRepository, timeout time.Duration) IOrderListUsecase {
	return &orderListUsecase{storage: storage, repo: repo, contextTimeout: timeout}
}

func (uc *orderListUsecase) Call(ctx context.Context, user *domain.User) ([]*OrderListResult, error) {
	tCtx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	orders, err := uc.repo.OrderList(tCtx, user.ID)
	if err != nil {
		return nil, err
	}

	result := make([]*OrderListResult, 0)
	for _, o := range orders {
		el := OrderListResult{Number: o.Number, Status: o.Status, CreatedAt: entities.RFC3339Time(o.CreatedAt)}

		if o.Status == domain.OrderStatusProcessed {
			val := entities.GDecimal(o.Accrual)
			el.Accrual = &val
		}

		result = append(result, &el)
	}

	return result, nil
}
