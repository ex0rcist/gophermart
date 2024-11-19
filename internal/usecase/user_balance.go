package usecase

import (
	"context"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/storage/repository"
)

type IGetUserBalanceUsecase interface {
	Call(ctx context.Context, user *domain.User) (*GetUserBalanceResult, error)
}

type GetUserBalanceResult struct {
	Current   entities.GDecimal `json:"current"`
	Withdrawn entities.GDecimal `json:"withdrawn"`
}

type getUserBalanceUsecase struct {
	storage        storage.IPGXStorage
	repo           repository.IUserRepository
	contextTimeout time.Duration
}

func NewGetUserBalanceUsecase(storage storage.IPGXStorage, repo repository.IUserRepository, timeout time.Duration) IGetUserBalanceUsecase {
	return &getUserBalanceUsecase{storage: storage, repo: repo, contextTimeout: timeout}
}

func (uc *getUserBalanceUsecase) Call(ctx context.Context, user *domain.User) (*GetUserBalanceResult, error) {
	tCtx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	b, w, err := uc.repo.UserGetBalance(tCtx, nil, user.ID)
	if err != nil {
		return nil, err
	}

	result := &GetUserBalanceResult{
		Current:   entities.GDecimal(*b),
		Withdrawn: entities.GDecimal(*w),
	}

	return result, nil
}
