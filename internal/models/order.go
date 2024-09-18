package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderID int32
type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID        OrderID
	UserID    UserID
	Number    string
	Status    OrderStatus
	Accrual   decimal.Decimal
	CreatedAt time.Time
	UpdatedAt time.Time
}
