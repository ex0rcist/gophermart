// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/usecase/withdrawal_list.go
//
// Generated by this command:
//
//	mockgen -source=./internal/usecase/withdrawal_list.go -destination=internal/usecase/mocks/withdrawal_list_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	domain "github.com/ex0rcist/gophermart/internal/domain"
	usecase "github.com/ex0rcist/gophermart/internal/usecase"
	gomock "go.uber.org/mock/gomock"
)

// MockIWithdrawalListUsecase is a mock of IWithdrawalListUsecase interface.
type MockIWithdrawalListUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockIWithdrawalListUsecaseMockRecorder
}

// MockIWithdrawalListUsecaseMockRecorder is the mock recorder for MockIWithdrawalListUsecase.
type MockIWithdrawalListUsecaseMockRecorder struct {
	mock *MockIWithdrawalListUsecase
}

// NewMockIWithdrawalListUsecase creates a new mock instance.
func NewMockIWithdrawalListUsecase(ctrl *gomock.Controller) *MockIWithdrawalListUsecase {
	mock := &MockIWithdrawalListUsecase{ctrl: ctrl}
	mock.recorder = &MockIWithdrawalListUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIWithdrawalListUsecase) EXPECT() *MockIWithdrawalListUsecaseMockRecorder {
	return m.recorder
}

// Call mocks base method.
func (m *MockIWithdrawalListUsecase) Call(ctx context.Context, user *domain.User) ([]*usecase.WithdrawalListResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", ctx, user)
	ret0, _ := ret[0].([]*usecase.WithdrawalListResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Call indicates an expected call of Call.
func (mr *MockIWithdrawalListUsecaseMockRecorder) Call(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockIWithdrawalListUsecase)(nil).Call), ctx, user)
}
