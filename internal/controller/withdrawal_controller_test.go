package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/usecase"
	mock_usecase "github.com/ex0rcist/gophermart/internal/usecase/mocks"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func TestWithdrawalController_WithdrawalList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawalListUsecase := mock_usecase.NewMockIWithdrawalListUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	withdrawalController := &WithdrawalController{
		WithdrawalListUsecase: mockWithdrawalListUsecase,
	}

	r.GET("/withdrawals", withdrawalController.WithdrawalList)

	withdrawals := []*usecase.WithdrawalListResult{
		{
			OrderNumber: "12345678903",
			Amount:      entities.GDecimal(decimal.NewFromFloat(100.0)),
			CreatedAt:   entities.RFC3339Time(time.Now()),
		},
		{
			OrderNumber: "98765432100",
			Amount:      entities.GDecimal(decimal.NewFromFloat(50.0)),
			CreatedAt:   entities.RFC3339Time(time.Now()),
		},
	}

	mockWithdrawalListUsecase.EXPECT().Call(gomock.Any(), gomock.Any()).Return(withdrawals, nil)

	req := httptest.NewRequest(http.MethodGet, "/withdrawals", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "12345678903")
	assert.Contains(t, w.Body.String(), "98765432100")
}

func TestWithdrawalController_WithdrawalList_NoWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawalListUsecase := mock_usecase.NewMockIWithdrawalListUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	withdrawalController := &WithdrawalController{
		WithdrawalListUsecase: mockWithdrawalListUsecase,
	}

	r.GET("/withdrawals", withdrawalController.WithdrawalList)

	mockWithdrawalListUsecase.EXPECT().Call(gomock.Any(), gomock.Any()).Return([]*usecase.WithdrawalListResult{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/withdrawals", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestWithdrawalController_WithdrawalList_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawalListUsecase := mock_usecase.NewMockIWithdrawalListUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	withdrawalController := &WithdrawalController{
		WithdrawalListUsecase: mockWithdrawalListUsecase,
	}

	r.GET("/withdrawals", withdrawalController.WithdrawalList)

	expectedError := errors.New("database error")

	mockWithdrawalListUsecase.EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil, expectedError)

	req := httptest.NewRequest(http.MethodGet, "/withdrawals", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestWithdrawalController_WithdrawalList_RecordNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawalListUsecase := mock_usecase.NewMockIWithdrawalListUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	withdrawalController := &WithdrawalController{
		WithdrawalListUsecase: mockWithdrawalListUsecase,
	}

	r.GET("/withdrawals", withdrawalController.WithdrawalList)

	mockWithdrawalListUsecase.EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil, storage.ErrRecordNotFound)

	req := httptest.NewRequest(http.MethodGet, "/withdrawals", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
