package usecases

import (
	"context"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/models"
	"github.com/ex0rcist/gophermart/internal/storage"
)

type UserGetBalanceResult struct {
	Current   entities.GDecimal `json:"current"`
	Withdrawn entities.GDecimal `json:"withdrawn"`
}

func UserGetBalance(ctx context.Context, s *storage.PGXStorage, u *models.User) (*UserGetBalanceResult, error) {
	res, err := s.UserGetBalance(ctx, nil, u.ID)
	if err != nil {
		return nil, err
	}

	result := &UserGetBalanceResult{
		Current:   entities.GDecimal(res.Balance),
		Withdrawn: entities.GDecimal(res.Withdrawn),
	}

	return result, nil
}
