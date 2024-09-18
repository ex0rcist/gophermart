package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type CreateWithdrawalDTO struct {
	UserID      models.UserID
	OrderNumber string
	Amount      decimal.Decimal
}

func (s *PGXStorage) WithdrawalCreate(ctx context.Context, tx pgx.Tx, d CreateWithdrawalDTO) error {
	stmt := `INSERT INTO withdrawals (user_id, order_number, amount) VALUES ($1, $2, $3)`

	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, stmt, d.UserID, d.OrderNumber, d.Amount)
	} else {
		_, err = s.pool.Exec(ctx, stmt, d.UserID, d.OrderNumber, d.Amount)
	}
	if err != nil {
		return fmt.Errorf("PGXStorage -> WithdrawalCreate() error: %w", err)
	}

	return nil
}

func (s *PGXStorage) WithdrawalList(ctx context.Context, userID models.UserID) ([]*models.Withdrawal, error) {
	stmt := `SELECT order_number, amount, created_at FROM withdrawals WHERE user_id = $1 ORDER BY created_at DESC`
	wds := make([]*models.Withdrawal, 0)

	rows, err := s.pool.Query(ctx, stmt, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrRecordNotFound
		}

		return nil, fmt.Errorf("PGXStorage -> WithdrawalList() error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		wd := &models.Withdrawal{}
		if err = rows.Scan(&wd.OrderNumber, &wd.Amount, &wd.CreatedAt); err != nil {
			return nil, fmt.Errorf("PGXStorage -> WithdrawalList() error: %w", err)
		}
		wds = append(wds, wd)
	}

	return wds, nil
}
