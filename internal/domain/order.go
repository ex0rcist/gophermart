package domain

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/jackc/pgx/v5"
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

type IOrderRepository interface {
	OrderCreate(ctx context.Context, o Order) (*Order, error)
	OrderFindByNumber(ctx context.Context, number string) (*Order, error)
	OrderList(ctx context.Context, userID UserID) ([]*Order, error)
	OrderListForUpdate(ctx context.Context) ([]*Order, error)
	OrderUpdate(ctx context.Context, tx pgx.Tx, o Order) error
}

// ========== OrderList ========== //

type IOrderListUsecase interface {
	Call(ctx context.Context, user *User) ([]*OrderListResult, error)
}

type OrderListResult struct {
	Number    string               `json:"number"`
	Status    OrderStatus          `json:"status"`
	Accrual   *entities.GDecimal   `json:"accrual,omitempty"` // без использования указателя omitempty не считает значение пустым
	CreatedAt entities.RFC3339Time `json:"uploaded_at"`
}

// ========== OrderCreate ========== //

var ErrOrderNotFound = errors.New("order not found")
var ErrOrderAlreadyRegistered = errors.New("order already registered")
var ErrOrderConflict = errors.New("order number already registered by another user")
var ErrInvalidOrderNumber = errors.New("invalid order number")

type IOrderCreateUsecase interface {
	Create(ctx context.Context, user *User, number string) (*Order, error)
	OrderFindByNumber(ctx context.Context, number string) (*Order, error)
}

// type OrderCreateDTO struct {
// 	UserID UserID
// 	Number string
// 	Status OrderStatus
// }

// type OrderUpdateDTO struct {
// 	ID      OrderID
// 	Status  OrderStatus
// 	Accrual decimal.Decimal
// }
