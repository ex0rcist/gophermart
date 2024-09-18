package usecases

import (
	"context"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/models"
	"github.com/ex0rcist/gophermart/internal/storage"
)

type OrderListForm struct {
	UserID models.UserID `binding:"required"`
	Number string        `binding:"required,min=3"`
	Status models.OrderStatus
}

type OrderListResult struct {
	Number    string               `json:"number"`
	Status    models.OrderStatus   `json:"status"`
	Accrual   *entities.GDecimal   `json:"accrual,omitempty"` // без использования указателя omitempty не считает значение пустым
	CreatedAt entities.RFC3339Time `json:"uploaded_at"`
}

func OrderList(ctx context.Context, s *storage.PGXStorage, u *models.User) ([]*OrderListResult, error) {
	orders, err := s.OrderList(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	result := make([]*OrderListResult, 0)
	for _, o := range orders {
		el := OrderListResult{Number: o.Number, Status: o.Status, CreatedAt: entities.RFC3339Time(o.CreatedAt)}

		if o.Status == models.OrderStatusProcessed {
			val := entities.GDecimal(o.Accrual)
			el.Accrual = &val
		}

		result = append(result, &el)
	}

	return result, nil
}
