package usecase

import (
	"context"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
)

type getUserBalanceUsecase struct {
	repo           domain.IUserRepository
	contextTimeout time.Duration
}

func NewGetUserBalanceUsecase(repo domain.IUserRepository, timeout time.Duration) domain.IGetUserBalanceUsecase {
	return &getUserBalanceUsecase{repo: repo, contextTimeout: timeout}
}

func (uc *getUserBalanceUsecase) Fetch(ctx context.Context, id domain.UserID) (*domain.GetUserBalanceResult, error) {
	tCtx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	b, w, err := uc.repo.UserGetBalance(tCtx, nil, id)
	if err != nil {
		return nil, err
	}

	result := &domain.GetUserBalanceResult{
		Current:   entities.GDecimal(*b),
		Withdrawn: entities.GDecimal(*w),
	}

	return result, nil
}
