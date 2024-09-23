package domain

import (
	"context"
	"errors"
	"time"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type UserID int32

type User struct {
	ID        UserID
	Login     string
	Password  string
	Balance   decimal.Decimal
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IUserRepository interface {
	UserCreate(ctx context.Context, login string, password string) (*User, error)
	UserFindByLogin(ctx context.Context, login string) (*User, error)
	UserGetBalance(ctx context.Context, tx pgx.Tx, id UserID) (*decimal.Decimal, *decimal.Decimal, error)
	UserUpdateBalanceAndWithdrawals(ctx context.Context, tx pgx.Tx, id UserID) error
}

// ========== Login ========== //

const LoginTokenLifetime = 1 * time.Hour

var ErrInvalidLoginOrPassword = errors.New("invalid login or password")

type LoginRequest struct {
	Login    string `json:"login" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=3"`
}

type ILoginUsecase interface {
	Call(ctx context.Context, form LoginRequest) (string, error)
	GetUserByLogin(ctx context.Context, req LoginRequest) (*User, error)
	ComparePassword(user *User, password string) error
	CreateAccessToken(user *User, secret entities.Secret, lifetime time.Duration) (string, error)
}

// ========== Register ========== //

var ErrUserAlreadyExists = errors.New("login already exists")

type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=3"`
}

type IRegisterUsecase interface {
	Call(ctx context.Context, form RegisterRequest) (string, error)
	GetUserByLogin(ctx context.Context, form RegisterRequest) (*User, error)
	CreateUser(ctx context.Context, login string, password string) (*User, error)
	CreateAccessToken(user *User, secret entities.Secret, lifetime time.Duration) (string, error)
}

// ========== GetUserBalance ========== //

type IGetUserBalanceUsecase interface {
	Fetch(ctx context.Context, id UserID) (*GetUserBalanceResult, error)
}

type GetUserBalanceResult struct {
	Current   entities.GDecimal `json:"current"`
	Withdrawn entities.GDecimal `json:"withdrawn"`
}

// ========== WithdrawBalance ========== //

var ErrInsufficientUserBalance = errors.New("insufficient user balance")

type WithdrawBalanceRequest struct {
	OrderNumber string          `json:"order" binding:"required,luhn"`
	Amount      decimal.Decimal `json:"sum" binding:"required"`
}

type IWithdrawBalanceUsecase interface {
	Call(ctx context.Context, user *User, req WithdrawBalanceRequest) error
}
