// Code generated by MockGen. DO NOT EDIT.
// Source: internal/accrual/client.go
//
// Generated by this command:
//
//      mockgen -source=internal/accrual/client.go
//

// Package mock_accrual is a generated GoMock package.
package mock_accrual

import (
        context "context"
        reflect "reflect"

        accrual "github.com/ex0rcist/gophermart/internal/accrual"
        gomock "go.uber.org/mock/gomock"
)

// MockIClient is a mock of IClient interface.
type MockIClient struct {
        ctrl     *gomock.Controller
        recorder *MockIClientMockRecorder
}

// MockIClientMockRecorder is the mock recorder for MockIClient.
type MockIClientMockRecorder struct {
        mock *MockIClient
}

// NewMockIClient creates a new mock instance.
func NewMockIClient(ctrl *gomock.Controller) *MockIClient {
        mock := &MockIClient{ctrl: ctrl}
        mock.recorder = &MockIClientMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIClient) EXPECT() *MockIClientMockRecorder {
        return m.recorder
}

// GetBonuses mocks base method.
func (m *MockIClient) GetBonuses(ctx context.Context, orderNumber string) (*accrual.Response, *accrual.ClientError) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "GetBonuses", ctx, orderNumber)
        ret0, _ := ret[0].(*accrual.Response)
        ret1, _ := ret[1].(*accrual.ClientError)
        return ret0, ret1
}

// GetBonuses indicates an expected call of GetBonuses.
func (mr *MockIClientMockRecorder) GetBonuses(ctx, orderNumber any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBonuses", reflect.TypeOf((*MockIClient)(nil).GetBonuses), ctx, orderNumber)
}