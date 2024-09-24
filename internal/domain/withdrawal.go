package domain

import (
	"time"

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
