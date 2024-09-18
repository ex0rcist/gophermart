package usecases

import (
	"context"

	"github.com/ex0rcist/gophermart/internal/models"
	"github.com/ex0rcist/gophermart/internal/storage"
)

type OrderCreateForm struct {
	UserID models.UserID `binding:"required"`
	Number string        `binding:"required,min=3"`
	Status models.OrderStatus
}

func OrderCreate(ctx context.Context, s *storage.PGXStorage, d OrderCreateForm) (*models.Order, error) {
	data := storage.OrderCreateDTO{
		UserID: d.UserID,
		Number: d.Number,
		Status: d.Status,
	}

	order, err := s.OrderCreate(ctx, data)
	if err != nil {
		return nil, err
	}

	return order, nil
}
