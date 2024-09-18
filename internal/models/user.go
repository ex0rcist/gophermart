package models

import (
	"time"

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
