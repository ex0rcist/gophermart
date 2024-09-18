package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/models"
	"github.com/jackc/pgx/v5"
)

type OrderCreateDTO struct {
	UserID models.UserID
	Number string
	Status models.OrderStatus
}

func (s *PGXStorage) OrderCreate(ctx context.Context, d OrderCreateDTO) (*models.Order, error) {
	stmt := `INSERT INTO orders (user_id, number, status) VALUES ($1, $2, $3) RETURNING id, user_id, number, status, accrual, created_at, updated_at`

	rows, err := s.pool.Query(ctx, stmt, d.UserID, d.Number, d.Status)
	if err != nil {
		return nil, fmt.Errorf("PGXStorage -> OrderCreate() error: %w", err)
	}
	defer rows.Close()

	order := new(models.Order)
	for rows.Next() {
		err = rows.Scan(
			&order.ID, &order.UserID, &order.Number, &order.Status,
			&order.Accrual, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("PGXStorage -> OrderCreate() error: %w", err)
		}
	}

	return order, nil
}

func (s *PGXStorage) OrderList(ctx context.Context, userID models.UserID) ([]*models.Order, error) {
	stmt := `SELECT number, status, accrual, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`
	orders := make([]*models.Order, 0)

	rows, err := s.pool.Query(ctx, stmt, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrRecordNotFound
		}

		return nil, fmt.Errorf("PGXStorage -> OrderList() error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		order := &models.Order{}
		if err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.CreatedAt); err != nil {
			return nil, fmt.Errorf("PGXStorage -> OrderList() error: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *PGXStorage) OrderFindByNumber(ctx context.Context, number string) (*models.Order, error) {
	stmt := `SELECT id, user_id, number, status, accrual, created_at, updated_at FROM orders WHERE number = $1`
	order := new(models.Order)

	err := s.pool.QueryRow(ctx, stmt, number).Scan(
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
