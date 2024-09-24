package app

import (
	"testing"

	mock_accrual "github.com/ex0rcist/gophermart/internal/accrual/mocks"
	mock_httpbackend "github.com/ex0rcist/gophermart/internal/http_backend/mocks"

	"github.com/ex0rcist/gophermart/internal/config"
	// mock_httpbackend "github.com/ex0rcist/gophermart/internal/httpbackend/mocks"
	mock_storage "github.com/ex0rcist/gophermart/internal/storage/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestApp_New_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testConfig := &config.Config{
		Server:  config.Server{Address: "127.0.0.1:8080", Secret: "secret_key"},
		Accrual: config.Accrual{Address: "127.0.0.1:8181"},
		DB:      config.DB{DSN: "postgres://user:password@localhost/dbname"},
	}

	mockStorage := mock_storage.NewMockIPGXStorage(ctrl)
	mockAccrualService := mock_accrual.NewMockIService(ctrl)
	mockHTTPBackend := mock_httpbackend.NewMockIHTTPBackend(ctrl)

	a, err := New(testConfig, mockStorage, mockAccrualService, mockHTTPBackend)

	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, testConfig, a.config)

	expectedStr := "app config: gophermart-address=127.0.0.1:8080; accrual-address=127.0.0.1:8181; database=postgres://user:password@localhost/dbname; secret=s********y"
	assert.Equal(t, expectedStr, a.String())
}
