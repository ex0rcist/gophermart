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

type withdrawalRepository struct {
	pool storage.IPGXPool
}

func NewWithdrawalRepository(pool storage.IPGXPool) domain.IWithdrawalRepository {
	return &withdrawalRepository{pool: pool}
}

func (repo *withdrawalRepository) WithdrawalCreate(ctx context.Context, tx pgx.Tx, w domain.Withdrawal) error {
	stmt := `INSERT INTO withdrawals (user_id, order_number, amount) VALUES ($1, $2, $3)`

	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, stmt, w.UserID, w.OrderNumber, w.Amount)
	} else {
		_, err = repo.pool.Exec(ctx, stmt, w.UserID, w.OrderNumber, w.Amount)
	}
	if err != nil {
		return fmt.Errorf("PGXStorage -> WithdrawalCreate() error: %w", err)
	}

	return nil
}

func (repo *withdrawalRepository) WithdrawalList(ctx context.Context, userID domain.UserID) ([]*domain.Withdrawal, error) {
	stmt := `SELECT order_number, amount, created_at FROM withdrawals WHERE user_id = $1 ORDER BY created_at DESC`
	wds := make([]*domain.Withdrawal, 0)

	rows, err := repo.pool.Query(ctx, stmt, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrRecordNotFound
		}

		return nil, fmt.Errorf("PGXStorage -> WithdrawalList() error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		wd := &domain.Withdrawal{}
		if err = rows.Scan(&wd.OrderNumber, &wd.Amount, &wd.CreatedAt); err != nil {
			return nil, fmt.Errorf("PGXStorage -> WithdrawalList() error: %w", err)
		}
		wds = append(wds, wd)
	}

	return wds, nil
}