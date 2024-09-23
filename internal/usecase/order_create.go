package usecase

import (
	"context"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/utils"
)

type orderCreateUsecase struct {
	storage        storage.IPGXStorage
	repo           domain.IOrderRepository
	contextTimeout time.Duration
}

func NewOrderCreateUsecase(storage storage.IPGXStorage, repo domain.IOrderRepository, timeout time.Duration) domain.IOrderCreateUsecase {
	return &orderCreateUsecase{storage: storage, repo: repo, contextTimeout: timeout}
}

func (uc *orderCreateUsecase) Create(ctx context.Context, user *domain.User, number string) (*domain.Order, error) {
	tCtx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	// валидируем номер заказа
	if !utils.LuhnCheck(number) {
		return nil, domain.ErrInvalidOrderNumber
	}

	// ищем заказ по номеру
	existingOrder, err := uc.OrderFindByNumber(tCtx, number)
	if err != nil && err != domain.ErrOrderNotFound {
		return nil, err
	}

	// проверяем существующий заказ
	if existingOrder != nil {
		if existingOrder.UserID == user.ID {
			return nil, domain.ErrOrderAlreadyRegistered
		}

		return nil, domain.ErrOrderConflict
	}

	order, err := uc.repo.OrderCreate(ctx, domain.Order{UserID: user.ID, Number: number, Status: domain.OrderStatusNew})
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (uc *orderCreateUsecase) OrderFindByNumber(ctx context.Context, number string) (*domain.Order, error) {
	tCtx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	order, err := uc.repo.OrderFindByNumber(tCtx, number)

	if err != nil {
		if err == storage.ErrRecordNotFound {
			return nil, domain.ErrOrderNotFound
		}

		return nil, err
	}

	return order, nil
}
