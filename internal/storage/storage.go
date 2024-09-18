package storage

import (
	"context"
	"fmt"

	"github.com/ex0rcist/gophermart/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXStorage struct {
	pool *pgxpool.Pool
}

type PGXTx struct {
	Tx         pgx.Tx
	committed  bool
	rolledBack bool
}

func NewPGXStorage(config config.DB) (*PGXStorage, error) {
	ctx := context.Background()

	err := runMigrations(config)
	if err != nil {
		return nil, fmt.Errorf("migrations failed: %w", err)
	}

	pgxConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, fmt.Errorf("pgxpool parsse config failed: %w", err)
	}

	pgxConfig.ConnConfig.Tracer = &dbQueryTracer{} // TODO: debug only
	// use with pgbouncer: pgxConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("pgxpool init failed: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("db connect failed: %w", err)
	}

	return &PGXStorage{pool: pool}, nil
}

func (s *PGXStorage) StartTx(ctx context.Context) (*PGXTx, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	return &PGXTx{Tx: tx, committed: false, rolledBack: false}, err
}

func (s *PGXStorage) RollbackTx(ctx context.Context, w *PGXTx) error {
	if w.rolledBack || w.committed {
		return nil
	}

	err := w.Tx.Rollback(ctx)
	if err == nil {
		w.rolledBack = true
	}

	return err
}

func (s *PGXStorage) CommitTx(ctx context.Context, w *PGXTx) error {
	if w.rolledBack || w.committed {
		return nil
	}

	err := w.Tx.Commit(ctx)
	if err == nil {
		w.committed = true
	}

	return err
}

func runMigrations(config config.DB) error {
	migrator := NewDatabaseMigrator(config.DSN, "file://internal/storage/migrations", 3)
	if err := migrator.Run(); err != nil {
		return err
	}

	return nil
}

func (s *PGXStorage) Close(ctx context.Context) error {
	s.pool.Close()
	return nil
}
