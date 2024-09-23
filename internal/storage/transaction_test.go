package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для pgx.Tx
type MockTx struct {
	mock.Mock
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestPGXTx_Rollback_Success(t *testing.T) {
	mockTx := new(PGXTxMock)
	ctx := context.Background()
	tx := &PGXTx{Tx: mockTx, ctx: ctx}

	// успешный откат транзакции
	mockTx.On("Rollback", ctx).Return(nil)

	err := tx.Rollback()
	assert.NoError(t, err)
	assert.True(t, tx.rolledBack)
	assert.False(t, tx.committed)

	// повторный откат должен пропускаться
	err = tx.Rollback()
	assert.NoError(t, err)
	mockTx.AssertNumberOfCalls(t, "Rollback", 1) // вызвался только один раз
}

func TestPGXTx_Commit_Success(t *testing.T) {
	mockTx := new(PGXTxMock)
	ctx := context.Background()
	tx := &PGXTx{Tx: mockTx, ctx: ctx}

	// успешный коммит транзакции
	mockTx.On("Commit", ctx).Return(nil)

	err := tx.Commit()
	assert.NoError(t, err)
	assert.True(t, tx.committed)
	assert.False(t, tx.rolledBack)

	// повторный коммит должен пропускаться
	err = tx.Commit()
	assert.NoError(t, err)
	mockTx.AssertNumberOfCalls(t, "Commit", 1) //  вызвался только один раз
}

func TestPGXTx_Rollback_After_Commit(t *testing.T) {
	mockTx := new(PGXTxMock)
	ctx := context.Background()
	tx := &PGXTx{Tx: mockTx, ctx: ctx}

	// успешный коммит
	mockTx.On("Commit", ctx).Return(nil)

	err := tx.Commit()
	assert.NoError(t, err)
	assert.True(t, tx.committed)

	// попытка отката после коммита
	err = tx.Rollback()
	assert.NoError(t, err)
	mockTx.AssertNotCalled(t, "Rollback") // откат не должен вызываться
}

func TestPGXTx_Commit_After_Rollback(t *testing.T) {
	mockTx := new(PGXTxMock)
	ctx := context.Background()
	tx := &PGXTx{Tx: mockTx, ctx: ctx}

	// Успешный откат
	mockTx.On("Rollback", ctx).Return(nil)

	err := tx.Rollback()
	assert.NoError(t, err)
	assert.True(t, tx.rolledBack)

	// попытка коммита после отката
	err = tx.Commit()
	assert.NoError(t, err)
	mockTx.AssertNotCalled(t, "Commit") // коммит не должен вызываться
}

func TestPGXTx_Rollback_Error(t *testing.T) {
	mockTx := new(PGXTxMock)
	ctx := context.Background()
	tx := &PGXTx{Tx: mockTx, ctx: ctx}

	// ошибка при откате
	mockTx.On("Rollback", ctx).Return(errors.New("rollback error"))

	err := tx.Rollback()
	assert.Error(t, err)
	assert.False(t, tx.rolledBack) // поле rolledBack не должно быть true при ошибке
}

func TestPGXTx_Commit_Error(t *testing.T) {
	mockTx := new(PGXTxMock)
	ctx := context.Background()
	tx := &PGXTx{Tx: mockTx, ctx: ctx}

	// ошибка при коммите
	mockTx.On("Commit", ctx).Return(errors.New("commit error"))

	err := tx.Commit()
	assert.Error(t, err)
	assert.False(t, tx.committed) // поле committed не должно быть true при ошибке
}
