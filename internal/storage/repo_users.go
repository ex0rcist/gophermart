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

type CreateUserDTO struct {
	Login    string
	Password string
}

func (s *PGXStorage) UserCreate(ctx context.Context, d CreateUserDTO) (*models.User, error) {
	stmt := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id, login, balance, created_at, updated_at`

	rows, err := s.pool.Query(ctx, stmt, d.Login, d.Password)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(models.User)
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Login, &user.Balance, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("PGXStorage -> UserCreate() error: %w", err)
	}

	return user, nil
}

func (s *PGXStorage) UserFindByLogin(ctx context.Context, login string) (*models.User, error) {
	stmt := `SELECT id, login, password, balance, created_at, updated_at FROM users WHERE login = $1`
	user := new(models.User)

	err := s.pool.QueryRow(ctx, stmt, login).Scan(
		&user.ID, &user.Login, &user.Password,
		&user.Balance, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrRecordNotFound
		}
		return nil, fmt.Errorf("PGXStorage -> UserFindByLogin() error: %w", err)
	}

	return user, nil
}

type UserGetBalanceDTO struct {
	Balance   decimal.Decimal
	Withdrawn decimal.Decimal
}

func (s *PGXStorage) UserGetBalance(ctx context.Context, tx pgx.Tx, id models.UserID) (*UserGetBalanceDTO, error) {
	stmt := `SELECT balance, withdrawn FROM users WHERE id = $1 FOR UPDATE`
	result := new(UserGetBalanceDTO)

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, stmt, id)
	} else {
		row = s.pool.QueryRow(ctx, stmt, id)
	}

	err := row.Scan(&result.Balance, &result.Withdrawn)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrRecordNotFound
		}
		return nil, fmt.Errorf("PGXStorage -> UserGetBalance() error: %w", err)
	}

	return result, nil
}

func (s *PGXStorage) UserUpdateBalanceAndWithdrawals(ctx context.Context, tx pgx.Tx, id models.UserID) error {
	stmt := `
	WITH
    	accruals AS (        
			SELECT o.user_id, COALESCE(SUM(o.accrual), 0) AS total_accrual
        	FROM orders o
        	WHERE o.status = 'PROCESSED'
        	GROUP BY o.user_id),
    	withdrawals AS (
        	SELECT w.user_id, COALESCE(SUM(w.amount), 0) AS total_withdrawn
        	FROM withdrawals w       
        	GROUP BY w.user_id)
	UPDATE users u
	SET
    	balance = COALESCE(a.total_accrual, 0) - COALESCE(w.total_withdrawn, 0),
    	withdrawn = COALESCE(w.total_withdrawn, 0)
	FROM
    	accruals a
    LEFT JOIN 
		withdrawals w ON a.user_id = w.user_id
	WHERE
    	u.id = $1 AND a.user_id = $1`

	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, stmt, id)
	} else {
		_, err = s.pool.Exec(ctx, stmt, id)
	}

	if err != nil {
		return fmt.Errorf("PGXStorage -> UserUpdateBalanceAndWithdrawals() error: %w", err)
	}

	return nil
}
