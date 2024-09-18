package usecases

import (
	"context"
	"errors"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/models"
	"github.com/ex0rcist/gophermart/internal/storage"
)

type OrderFindForm struct {
	Number string `binding:"required,min=3,luhn"`
}

var (
	ErrOrderNotFound = errors.New("order not found")
)

func OrderFindByNumber(ctx context.Context, s *storage.PGXStorage, d OrderFindForm) (*models.Order, error) {
	order, err := s.OrderFindByNumber(ctx, d.Number)
	if err != nil {
		if err == entities.ErrRecordNotFound {
			return nil, ErrOrderNotFound
		}

		return nil, err
	}

	return order, nil
}
