package usecases

import (
	"context"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/models"
	"github.com/ex0rcist/gophermart/internal/storage"
)

type WithdrawalListResult struct {
	OrderNumber string               `json:"order"`
	Amount      entities.GDecimal    `json:"sum"`
	CreatedAt   entities.RFC3339Time `json:"processed_at"`
}

func WithdrawalList(ctx context.Context, s *storage.PGXStorage, u *models.User) ([]*WithdrawalListResult, error) {
	wds, err := s.WithdrawalList(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	result := make([]*WithdrawalListResult, 0)
	for _, w := range wds {
		el := WithdrawalListResult{
			OrderNumber: w.OrderNumber,
			Amount:      entities.GDecimal(w.Amount),
			CreatedAt:   entities.RFC3339Time(w.CreatedAt),
		}
		result = append(result, &el)
	}

	return result, nil
}
