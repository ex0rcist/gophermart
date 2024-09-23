package storage

import (
	"context"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

var _ IPGXPool = (*PGXPoolMock)(nil)

type PGXPoolMock struct {
	mock.Mock
}

func NewPGXPoolMock() *PGXPoolMock {
	return new(PGXPoolMock)
}

func (m *PGXPoolMock) Begin(ctx context.Context) (pgx.Tx, error) {
	mArgs := m.Called(ctx)
	return mArgs.Get(0).(pgx.Tx), mArgs.Error(1)
}

func (m *PGXPoolMock) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	mArgs := m.Called(ctx)
	return mArgs.Get(0).(pgx.Tx), mArgs.Error(1)
}

func (m *PGXPoolMock) Acquire(ctx context.Context) (c *pgxpool.Conn, err error) {
	mArgs := m.Called(ctx)
	return mArgs.Get(0).(*pgxpool.Conn), mArgs.Error(1)
}

func (m *PGXPoolMock) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	mArgs := m.Called(ctx, sql, args)
	return mArgs.Get(0).(pgconn.CommandTag), mArgs.Error(1)
}

func (m *PGXPoolMock) SendBatch(ctx context.Context, b *pgx.Batch) (br pgx.BatchResults) {
	mArgs := m.Called(ctx, b)
	return mArgs.Get(0).(pgx.BatchResults)
}

func (m *PGXPoolMock) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *PGXPoolMock) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	mArgs := m.Called(ctx, sql, args)
	return mArgs.Get(0).(pgx.Rows), mArgs.Error(1)
}

func (m *PGXPoolMock) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	mArgs := m.Called(ctx, sql, args)
	return mArgs.Get(0).(pgx.Row)
}

func (m *PGXPoolMock) Close() {
	_ = m.Called()
}

// ************** PGXRowMock ************** //

type PGXRowMock struct {
	mock.Mock
}

func (m *PGXRowMock) Scan(args ...any) error {
	mArgs := m.Called(args...)
	return mArgs.Error(0)
}

// ************** PGXRowsMock ************** //

type PGXRowsMock struct {
	mock.Mock
}

func (m *PGXRowsMock) Close() {
	_ = m.Called()
}

func (m *PGXRowsMock) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *PGXRowsMock) Scan(args ...any) error {
	mArgs := m.Called(args...)
	return mArgs.Error(0)
}

func (m *PGXRowsMock) Err() error {
	mArgs := m.Called()
	return mArgs.Error(0)
}

func (m *PGXRowsMock) CommandTag() pgconn.CommandTag {
	mArgs := m.Called()
	return mArgs.Get(0).(pgconn.CommandTag)
}

func (m *PGXRowsMock) FieldDescriptions() []pgconn.FieldDescription {
	mArgs := m.Called()
	return mArgs.Get(0).([]pgconn.FieldDescription)
}

func (m *PGXRowsMock) Values() ([]any, error) {
	mArgs := m.Called()
	return mArgs.Get(0).([]any), mArgs.Error(1)
}

func (m *PGXRowsMock) RawValues() [][]byte {
	mArgs := m.Called()
	return mArgs.Get(0).([][]byte)
}

func (m *PGXRowsMock) Conn() *pgx.Conn {
	mArgs := m.Called()
	return mArgs.Get(0).(*pgx.Conn)
}

// ************** PGXTxMock ************** //

type PGXTxMock struct {
	mock.Mock
}

func (m *PGXTxMock) Begin(ctx context.Context) (pgx.Tx, error) {
	mArgs := m.Called(ctx)
	return mArgs.Get(0).(pgx.Tx), mArgs.Error(1)
}

func (m *PGXTxMock) Commit(ctx context.Context) error {
	mArgs := m.Called(ctx)
	return mArgs.Error(0)
}

func (m *PGXTxMock) Rollback(ctx context.Context) error {
	mArgs := m.Called(ctx)
	return mArgs.Error(0)
}

func (m *PGXTxMock) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	mArgs := m.Called(ctx, tableName)
	return mArgs.Get(0).(int64), mArgs.Error(1)
}

func (m *PGXTxMock) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	mArgs := m.Called(ctx, b)
	return mArgs.Get(0).(pgx.BatchResults)
}

func (m *PGXTxMock) LargeObjects() pgx.LargeObjects {
	mArgs := m.Called()
	return mArgs.Get(0).(pgx.LargeObjects)
}

func (m *PGXTxMock) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	mArgs := m.Called(ctx, name, sql)
	return mArgs.Get(0).(*pgconn.StatementDescription), mArgs.Error(1)
}

func (m *PGXTxMock) Exec(ctx context.Context, sql string, args ...any) (commandTag pgconn.CommandTag, err error) {
	varargs := append([]any{ctx, sql}, args...)
	mArgs := m.Called(varargs...)
	return mArgs.Get(0).(pgconn.CommandTag), mArgs.Error(1)
}
func (m *PGXTxMock) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	mArgs := m.Called(ctx, sql, args)
	return mArgs.Get(0).(pgx.Rows), mArgs.Error(1)
}
func (m *PGXTxMock) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	mArgs := m.Called(ctx, sql, args)
	return mArgs.Get(0).(pgx.Row)
}
func (m *PGXTxMock) Conn() *pgx.Conn {
	mArgs := m.Called()
	return mArgs.Get(0).(*pgx.Conn)
}

// ************** PGXStorage ************** //

type MockPGXStorage struct {
	mock.Mock
}

func (m *MockPGXStorage) Close() {
	_ = m.Called()
}

func (m *MockPGXStorage) OrderCreate(ctx context.Context, order domain.Order) (*domain.Order, error) {
	mArgs := m.Called(ctx, order)
	return mArgs.Get(0).(*domain.Order), mArgs.Error(1)
}

func (m *MockPGXStorage) OrderFindByNumber(ctx context.Context, number string) (*domain.Order, error) {
	mArgs := m.Called(ctx, number)
	return mArgs.Get(0).(*domain.Order), mArgs.Error(1)
}

func (m *MockPGXStorage) OrderList(ctx context.Context, userID domain.UserID) ([]*domain.Order, error) {
	mArgs := m.Called(ctx, userID)
	return mArgs.Get(0).([]*domain.Order), mArgs.Error(1)
}

func (m *MockPGXStorage) OrderListForUpdate(ctx context.Context) ([]*domain.Order, error) {
	mArgs := m.Called(ctx)
	return mArgs.Get(0).([]*domain.Order), mArgs.Error(1)
}

func (m *MockPGXStorage) OrderUpdate(ctx context.Context, tx pgx.Tx, order domain.Order) error {
	mArgs := m.Called(ctx, tx, order)
	return mArgs.Error(0)
}

func (m *MockPGXStorage) StartTx(ctx context.Context) (*PGXTx, error) {
	mArgs := m.Called(ctx)
	return mArgs.Get(0).(*PGXTx), mArgs.Error(1)
}

func (m *MockPGXStorage) UserCreate(ctx context.Context, login string, password string) (*domain.User, error) {
	mArgs := m.Called(ctx, login, password)
	return mArgs.Get(0).(*domain.User), mArgs.Error(1)
}

func (m *MockPGXStorage) UserFindByLogin(ctx context.Context, login string) (*domain.User, error) {
	mArgs := m.Called(ctx, login)
	return mArgs.Get(0).(*domain.User), mArgs.Error(1)
}

func (m *MockPGXStorage) UserGetBalance(ctx context.Context, tx pgx.Tx, id domain.UserID) (*decimal.Decimal, *decimal.Decimal, error) {
	mArgs := m.Called(ctx, tx, id)
	return mArgs.Get(0).(*decimal.Decimal), mArgs.Get(1).(*decimal.Decimal), mArgs.Error(2)
}

func (m *MockPGXStorage) UserUpdateBalanceAndWithdrawals(ctx context.Context, tx pgx.Tx, id domain.UserID) error {
	mArgs := m.Called(ctx, tx, id)
	return mArgs.Error(0)
}

func (m *MockPGXStorage) WithdrawalCreate(ctx context.Context, tx pgx.Tx, w domain.Withdrawal) error {
	mArgs := m.Called(ctx, tx, w)
	return mArgs.Error(0)
}

func (m *MockPGXStorage) WithdrawalList(ctx context.Context, userID domain.UserID) ([]*domain.Withdrawal, error) {
	mArgs := m.Called(ctx, userID)
	return mArgs.Get(0).([]*domain.Withdrawal), mArgs.Error(1)
}
