// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/repository/orders.go
//
// Generated by this command:
//
//      mockgen -source=internal/storage/repository/orders.go
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
        context "context"
        reflect "reflect"

        domain "github.com/ex0rcist/gophermart/internal/domain"
        pgx "github.com/jackc/pgx/v5"
        gomock "go.uber.org/mock/gomock"
)

// MockIOrderRepository is a mock of IOrderRepository interface.
type MockIOrderRepository struct {
        ctrl     *gomock.Controller
        recorder *MockIOrderRepositoryMockRecorder
}

// MockIOrderRepositoryMockRecorder is the mock recorder for MockIOrderRepository.
type MockIOrderRepositoryMockRecorder struct {
        mock *MockIOrderRepository
}

// NewMockIOrderRepository creates a new mock instance.
func NewMockIOrderRepository(ctrl *gomock.Controller) *MockIOrderRepository {
        mock := &MockIOrderRepository{ctrl: ctrl}
        mock.recorder = &MockIOrderRepositoryMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIOrderRepository) EXPECT() *MockIOrderRepositoryMockRecorder {
        return m.recorder
}

// OrderCreate mocks base method.
func (m *MockIOrderRepository) OrderCreate(ctx context.Context, o domain.Order) (*domain.Order, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "OrderCreate", ctx, o)
        ret0, _ := ret[0].(*domain.Order)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// OrderCreate indicates an expected call of OrderCreate.
func (mr *MockIOrderRepositoryMockRecorder) OrderCreate(ctx, o any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrderCreate", reflect.TypeOf((*MockIOrderRepository)(nil).OrderCreate), ctx, o)
}

// OrderFindByNumber mocks base method.
func (m *MockIOrderRepository) OrderFindByNumber(ctx context.Context, number string) (*domain.Order, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "OrderFindByNumber", ctx, number)
        ret0, _ := ret[0].(*domain.Order)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// OrderFindByNumber indicates an expected call of OrderFindByNumber.
func (mr *MockIOrderRepositoryMockRecorder) OrderFindByNumber(ctx, number any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrderFindByNumber", reflect.TypeOf((*MockIOrderRepository)(nil).OrderFindByNumber), ctx, number)
}

// OrderList mocks base method.
func (m *MockIOrderRepository) OrderList(ctx context.Context, userID domain.UserID) ([]*domain.Order, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "OrderList", ctx, userID)
        ret0, _ := ret[0].([]*domain.Order)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// OrderList indicates an expected call of OrderList.
func (mr *MockIOrderRepositoryMockRecorder) OrderList(ctx, userID any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrderList", reflect.TypeOf((*MockIOrderRepository)(nil).OrderList), ctx, userID)
}

// OrderListForUpdate mocks base method.
func (m *MockIOrderRepository) OrderListForUpdate(ctx context.Context) ([]*domain.Order, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "OrderListForUpdate", ctx)
        ret0, _ := ret[0].([]*domain.Order)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// OrderListForUpdate indicates an expected call of OrderListForUpdate.
func (mr *MockIOrderRepositoryMockRecorder) OrderListForUpdate(ctx any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrderListForUpdate", reflect.TypeOf((*MockIOrderRepository)(nil).OrderListForUpdate), ctx)
}

// OrderUpdate mocks base method.
func (m *MockIOrderRepository) OrderUpdate(ctx context.Context, tx pgx.Tx, o domain.Order) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "OrderUpdate", ctx, tx, o)
        ret0, _ := ret[0].(error)
        return ret0
}

// OrderUpdate indicates an expected call of OrderUpdate.
func (mr *MockIOrderRepositoryMockRecorder) OrderUpdate(ctx, tx, o any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrderUpdate", reflect.TypeOf((*MockIOrderRepository)(nil).OrderUpdate), ctx, tx, o)
}