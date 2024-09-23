package usecase

import (
	"context"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
)

type withdrawalListUsecase struct {
	repo           domain.IWithdrawalRepository
	contextTimeout time.Duration
}

func NewWithdrawalListUsecase(repo domain.IWithdrawalRepository, timeout time.Duration) domain.IWithdrawalListUsecase {
	return &withdrawalListUsecase{repo: repo, contextTimeout: timeout}
}

func (uc *withdrawalListUsecase) Call(ctx context.Context, u *domain.User) ([]*domain.WithdrawalListResult, error) {
	wds, err := uc.repo.WithdrawalList(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.WithdrawalListResult, 0)
	for _, w := range wds {
		el := domain.WithdrawalListResult{
			OrderNumber: w.OrderNumber,
			Amount:      entities.GDecimal(w.Amount),
			CreatedAt:   entities.RFC3339Time(w.CreatedAt),
		}
		result = append(result, &el)
	}

	return result, nil
}
