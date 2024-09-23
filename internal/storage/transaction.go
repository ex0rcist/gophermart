package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PGXTx struct {
	Tx         pgx.Tx
	ctx        context.Context
	committed  bool
	rolledBack bool
}

func (t *PGXTx) Rollback() error {
	if t.rolledBack || t.committed {
		return nil
	}

	err := t.Tx.Rollback(t.ctx)
	if err == nil {
		t.rolledBack = true
	}

	return err
}

func (t *PGXTx) Commit() error {
	if t.rolledBack || t.committed {
		return nil
	}

	err := t.Tx.Commit(t.ctx)
	if err == nil {
		t.committed = true
	}

	return err
}
