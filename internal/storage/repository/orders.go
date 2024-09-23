package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/jackc/pgx/v5"
)

type orderRepository struct {
	pool storage.IPGXPool
}

func NewOrderRepository(pool storage.IPGXPool) domain.IOrderRepository {
	return &orderRepository{pool: pool}
}

func (repo *orderRepository) OrderCreate(ctx context.Context, order domain.Order) (*domain.Order, error) {
	stmt := `INSERT INTO orders (user_id, number, status) VALUES ($1, $2, $3) RETURNING id, user_id, number, status, accrual, created_at, updated_at`

	rows, err := repo.pool.Query(ctx, stmt, order.UserID, order.Number, order.Status)
	if err != nil {
		return nil, fmt.Errorf("PGXStorage -> OrderCreate() error: %w", err)
	}
	defer rows.Close()

	newOrder := new(domain.Order)
	for rows.Next() {
		err = rows.Scan(
			&newOrder.ID, &newOrder.UserID, &newOrder.Number, &newOrder.Status,
			&newOrder.Accrual, &newOrder.CreatedAt, &newOrder.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("PGXStorage -> OrderCreate() error: %w", err)
		}
	}

	return newOrder, nil
}

func (repo *orderRepository) OrderList(ctx context.Context, userID domain.UserID) ([]*domain.Order, error) {
	stmt := `SELECT number, status, accrual, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`
	orders := make([]*domain.Order, 0)

	rows, err := repo.pool.Query(ctx, stmt, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrRecordNotFound
		}

		return nil, fmt.Errorf("PGXStorage -> OrderList() error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		order := &domain.Order{}
		if err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.CreatedAt); err != nil {
			return nil, fmt.Errorf("PGXStorage -> OrderList() error: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (repo *orderRepository) OrderFindByNumber(ctx context.Context, number string) (*domain.Order, error) {
	stmt := `SELECT id, user_id, number, status, accrual, created_at, updated_at FROM orders WHERE number = $1`
	order := new(domain.Order)

	err := repo.pool.QueryRow(ctx, stmt, number).Scan(
		&order.ID, &order.UserID, &order.Number, &order.Status,
		&order.Accrual, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrRecordNotFound
		}
		return nil, fmt.Errorf("PGXStorage -> OrderFindByNumber() error: %w", err)
	}

	return order, nil
}

func (repo *orderRepository) OrderListForUpdate(ctx context.Context) ([]*domain.Order, error) {
	stmt := `SELECT id, user_id, number, status, created_at FROM orders WHERE status IN ('NEW', 'PROCESSING');`
	orders := make([]*domain.Order, 0)

	rows, err := repo.pool.Query(ctx, stmt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrRecordNotFound
		}
		return nil, err
	}

	for rows.Next() {
		order := &domain.Order{}
		if err = rows.Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.CreatedAt); err != nil {
			return nil, fmt.Errorf("PGXStorage -> OrderListForUpdate() error: %w", err)
		}
		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("PGXStorage -> OrderListForUpdate() error: %w", err)
	}

	return orders, nil
}

func (repo *orderRepository) OrderUpdate(ctx context.Context, tx pgx.Tx, order domain.Order) error {
	stmt := `UPDATE orders SET status = $1, accrual = $2 WHERE id = $3`

	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, stmt, order.Status, order.Accrual, order.ID)
	} else {
		_, err = repo.pool.Exec(ctx, stmt, order.Status, order.Accrual, order.ID)
	}

	if err != nil {
		return fmt.Errorf("PGXStorage -> OrderUpdate() error: %w", err)
	}

	return nil
}
