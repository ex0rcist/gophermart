// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/storage.go
//
// Generated by this command:
//
//      mockgen -source=internal/storage/storage.go
//

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
        reflect "reflect"

        storage "github.com/ex0rcist/gophermart/internal/storage"
        gomock "go.uber.org/mock/gomock"
)

// MockIPGXStorage is a mock of IPGXStorage interface.
type MockIPGXStorage struct {
        ctrl     *gomock.Controller
        recorder *MockIPGXStorageMockRecorder
}

// MockIPGXStorageMockRecorder is the mock recorder for MockIPGXStorage.
type MockIPGXStorageMockRecorder struct {
        mock *MockIPGXStorage
}

// NewMockIPGXStorage creates a new mock instance.
func NewMockIPGXStorage(ctrl *gomock.Controller) *MockIPGXStorage {
        mock := &MockIPGXStorage{ctrl: ctrl}
        mock.recorder = &MockIPGXStorageMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPGXStorage) EXPECT() *MockIPGXStorageMockRecorder {
        return m.recorder
}

// Close mocks base method.
func (m *MockIPGXStorage) Close() {
        m.ctrl.T.Helper()
        m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockIPGXStorageMockRecorder) Close() *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockIPGXStorage)(nil).Close))
}

// GetPool mocks base method.
func (m *MockIPGXStorage) GetPool() storage.IPGXPool {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "GetPool")
        ret0, _ := ret[0].(storage.IPGXPool)
        return ret0
}

// GetPool indicates an expected call of GetPool.
func (mr *MockIPGXStorageMockRecorder) GetPool() *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPool", reflect.TypeOf((*MockIPGXStorage)(nil).GetPool))
}