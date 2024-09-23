package domain

import (
	"context"
	"time"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type WithdrawalID int32

type Withdrawal struct {
	ID          WithdrawalID
	UserID      UserID
	OrderNumber string
	Amount      decimal.Decimal
	CreatedAt   time.Time
}

type IWithdrawalRepository interface {
	WithdrawalCreate(ctx context.Context, tx pgx.Tx, w Withdrawal) error
	WithdrawalList(ctx context.Context, userID UserID) ([]*Withdrawal, error)
}

// ========== WithdrawalList ========== //

type IWithdrawalListUsecase interface {
	Call(ctx context.Context, user *User) ([]*WithdrawalListResult, error)
}

type WithdrawalListResult struct {
	OrderNumber string               `json:"order"`
	Amount      entities.GDecimal    `json:"sum"`
	CreatedAt   entities.RFC3339Time `json:"processed_at"`
}
