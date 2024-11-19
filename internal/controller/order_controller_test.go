package controller

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/usecase"
	mock_usecase "github.com/ex0rcist/gophermart/internal/usecase/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrderController_CreateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCreateUsecase := mock_usecase.NewMockIOrderCreateUsecase(ctrl)
	mockListUsecase := mock_usecase.NewMockIOrderListUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	orderController := &OrderController{
		OrderCreateUsecase: mockCreateUsecase,
		OrderListUsecase:   mockListUsecase,
	}

	// фактические маршруты не важны
	r.POST("/orders", orderController.CreateOrder)

	mockCreateUsecase.EXPECT().Create(gomock.Any(), gomock.Any(), "12345678903").Return(&domain.Order{}, nil)

	reqBody := []byte("12345678903")
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestOrderController_CreateOrder_InvalidOrderNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCreateUsecase := mock_usecase.NewMockIOrderCreateUsecase(ctrl)
	mockListUsecase := mock_usecase.NewMockIOrderListUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	orderController := &OrderController{
		OrderCreateUsecase: mockCreateUsecase,
		OrderListUsecase:   mockListUsecase,
	}

	r.POST("/orders", orderController.CreateOrder)

	mockCreateUsecase.EXPECT().Create(gomock.Any(), gomock.Any(), "invalid").Return(nil, usecase.ErrInvalidOrderNumber)

	reqBody := []byte("invalid")
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Contains(t, w.Body.String(), usecase.ErrInvalidOrderNumber.Error())
}

func TestOrderController_CreateOrder_OrderAlreadyRegistered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCreateUsecase := mock_usecase.NewMockIOrderCreateUsecase(ctrl)
	mockListUsecase := mock_usecase.NewMockIOrderListUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	orderController := &OrderController{
		OrderCreateUsecase: mockCreateUsecase,
		OrderListUsecase:   mockListUsecase,
	}

	r.POST("/orders", orderController.CreateOrder)

	mockCreateUsecase.EXPECT().Create(gomock.Any(), gomock.Any(), "12345678903").Return(nil, usecase.ErrOrderAlreadyRegistered)

	reqBody := []byte("12345678903")
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderController_CreateOrder_OrderConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCreateUsecase := mock_usecase.NewMockIOrderCreateUsecase(ctrl)
	mockListUsecase := mock_usecase.NewMockIOrderListUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	orderController := &OrderController{
		OrderCreateUsecase: mockCreateUsecase,
		OrderListUsecase:   mockListUsecase,
	}

	r.POST("/orders", orderController.CreateOrder)

	mockCreateUsecase.EXPECT().Create(gomock.Any(), gomock.Any(), "12345678903").Return(nil, usecase.ErrOrderConflict)

	reqBody := []byte("12345678903")
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestOrderController_OrderList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCreateUsecase := mock_usecase.NewMockIOrderCreateUsecase(ctrl)
	mockListUsecase := mock_usecase.NewMockIOrderListUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	orderController := &OrderController{
		OrderCreateUsecase: mockCreateUsecase,
		OrderListUsecase:   mockListUsecase,
	}

	r.GET("/orders", orderController.OrderList)

	orders := []*usecase.OrderListResult{
		{Number: "12345678903", Status: domain.OrderStatusNew},
	}

	mockListUsecase.EXPECT().Call(gomock.Any(), gomock.Any()).Return(orders, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "12345678903")
}

func TestOrderController_OrderList_NoOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCreateUsecase := mock_usecase.NewMockIOrderCreateUsecase(ctrl)
	mockListUsecase := mock_usecase.NewMockIOrderListUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	orderController := &OrderController{
		OrderCreateUsecase: mockCreateUsecase,
		OrderListUsecase:   mockListUsecase,
	}

	r.GET("/orders", orderController.OrderList)

	mockListUsecase.EXPECT().Call(gomock.Any(), gomock.Any()).Return([]*usecase.OrderListResult{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestOrderController_OrderList_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCreateUsecase := mock_usecase.NewMockIOrderCreateUsecase(ctrl)
	mockListUsecase := mock_usecase.NewMockIOrderListUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	orderController := &OrderController{
		OrderCreateUsecase: mockCreateUsecase,
		OrderListUsecase:   mockListUsecase,
	}

	r.GET("/orders", orderController.OrderList)

	expectedError := errors.New("database error")

	mockListUsecase.EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil, expectedError)

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
