package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type IUserRepository interface {
	UserCreate(ctx context.Context, login string, password string) (*domain.User, error)
	UserFindByLogin(ctx context.Context, login string) (*domain.User, error)
	UserGetBalance(ctx context.Context, tx pgx.Tx, id domain.UserID) (*decimal.Decimal, *decimal.Decimal, error)
	UserUpdateBalanceAndWithdrawals(ctx context.Context, tx pgx.Tx, id domain.UserID) error
}

type userRepository struct {
	pool storage.IPGXPool
}

func NewUserRepository(pool storage.IPGXPool) IUserRepository {
	return &userRepository{pool: pool}
}

func (repo *userRepository) UserCreate(ctx context.Context, login string, password string) (*domain.User, error) {
	stmt := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id, login, balance, created_at, updated_at`

	rows, err := repo.pool.Query(ctx, stmt, login, password)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(domain.User)
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Login, &user.Balance, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("userRepository -> UserCreate() error: %w", err)
	}

	return user, nil
}

func (repo *userRepository) UserFindByLogin(ctx context.Context, login string) (*domain.User, error) {
	stmt := `SELECT id, login, password, balance, created_at, updated_at FROM users WHERE login = $1`
	user := new(domain.User)

	err := repo.pool.QueryRow(ctx, stmt, login).Scan(
		&user.ID, &user.Login, &user.Password,
		&user.Balance, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrRecordNotFound
		}
		return nil, fmt.Errorf("userRepository -> UserFindByLogin() error: %w", err)
	}

	return user, nil
}

func (repo *userRepository) UserGetBalance(ctx context.Context, tx pgx.Tx, id domain.UserID) (*decimal.Decimal, *decimal.Decimal, error) {
	stmt := `SELECT balance, withdrawn FROM users WHERE id = $1 FOR UPDATE`

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, stmt, id)
	} else {
		row = repo.pool.QueryRow(ctx, stmt, id)
	}

	var b, w decimal.Decimal
	err := row.Scan(&b, &w)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, storage.ErrRecordNotFound
		}
		return nil, nil, fmt.Errorf("userRepository -> UserGetBalance() error: %w", err)
	}

	return &b, &w, nil
}

func (repo *userRepository) UserUpdateBalanceAndWithdrawals(ctx context.Context, tx pgx.Tx, id domain.UserID) error {
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
		_, err = repo.pool.Exec(ctx, stmt, id)
	}

	if err != nil {
		return fmt.Errorf("userRepository -> UserUpdateBalanceAndWithdrawals() error: %w", err)
	}

	return nil
}
