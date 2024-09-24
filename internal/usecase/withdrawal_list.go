package usecase

import (
	"context"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/storage/repository"
)

type IWithdrawalListUsecase interface {
	Call(ctx context.Context, user *domain.User) ([]*WithdrawalListResult, error)
}

type WithdrawalListResult struct {
	OrderNumber string               `json:"order"`
	Amount      entities.GDecimal    `json:"sum"`
	CreatedAt   entities.RFC3339Time `json:"processed_at"`
}

type withdrawalListUsecase struct {
	storage        storage.IPGXStorage
	repo           repository.IWithdrawalRepository
	contextTimeout time.Duration
}

func NewWithdrawalListUsecase(storage storage.IPGXStorage, repo repository.IWithdrawalRepository, timeout time.Duration) IWithdrawalListUsecase {
	return &withdrawalListUsecase{storage: storage, repo: repo, contextTimeout: timeout}
}

func (uc *withdrawalListUsecase) Call(ctx context.Context, u *domain.User) ([]*WithdrawalListResult, error) {
	tCtx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	wds, err := uc.repo.WithdrawalList(tCtx, u.ID)
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
