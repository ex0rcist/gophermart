package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/ex0rcist/gophermart/internal/config"
	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/storage/tracer"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var _ IPGXStorage = (*PGXStorage)(nil)

type PGXStorage struct {
	pool IPGXPool
}

func NewPGXStorage(config config.DB, pool IPGXPool, migrate bool) (*PGXStorage, error) {
	var err error

	if migrate {
		if err := runMigrations(config); err != nil {
			return nil, fmt.Errorf("runMigrations() failed: %w", err)
		}
	}

	if pool == nil {
		pool, err = createPool(context.Background(), config)
		if err != nil {
			return nil, fmt.Errorf("create pool failed: %w", err)
		}
	}

	return &PGXStorage{pool: pool}, err
}

func (s *PGXStorage) StartTx(ctx context.Context) (*PGXTx, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	return &PGXTx{Tx: tx, ctx: ctx, committed: false, rolledBack: false}, err
}

func (s *PGXStorage) GetPool() IPGXPool {
	return s.pool
}

func (s *PGXStorage) Close() {
	s.pool.Close()
}

func createPool(ctx context.Context, config config.DB) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, fmt.Errorf("pgxpool parse config failed: %w", err)
	}

	pgxConfig.ConnConfig.Tracer = tracer.NewDBQueryTracer()

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("pgxpool init failed: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("db connect failed: %w", err)
	}

	return pool, nil
}

func runMigrations(config config.DB) error {
	migrator, err := migrate.New(config.MigrationsSource, config.DSN)
	if err != nil {
		return fmt.Errorf("migrate.New() failed: %w", err)
	}

	err = migrator.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logging.LogInfo("migrations: no change")
			return nil
		}
		return fmt.Errorf("migrations failed: %w", err)
	}

	defer func() {
		srcErr, dbErr := migrator.Close()
		if srcErr != nil {
			logging.LogError(srcErr, "failed closing migrator", srcErr.Error())
		}
		if dbErr != nil {
			logging.LogError(dbErr, "failed closing migrator", dbErr.Error())
		}
	}()

	return nil
}
