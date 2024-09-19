package models

import (
	"fmt"
	"strings"
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

func (o *Order) String() string {
	str := []string{
		fmt.Sprintf("user_id=%d", o.UserID),
		fmt.Sprintf("number=%s", o.Number),
		fmt.Sprintf("status=%s", o.Status),
		fmt.Sprintf("accrual=%s", o.Accrual),
	}

	return fmt.Sprintf("order(id=%d)[%s]", o.ID, strings.Join(str, ";"))
}
