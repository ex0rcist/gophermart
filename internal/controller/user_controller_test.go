package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/usecase"
	mock_usecase "github.com/ex0rcist/gophermart/internal/usecase/mocks"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func TestUserController_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoginUsecase := mock_usecase.NewMockILoginUsecase(ctrl)
	mockRegisterUsecase := mock_usecase.NewMockIRegisterUsecase(ctrl)
	mockGetUserBalanceUsecase := mock_usecase.NewMockIGetUserBalanceUsecase(ctrl)
	mockWithdrawBalanceUsecase := mock_usecase.NewMockIWithdrawBalanceUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	userController := &UserController{
		LoginUsecase:           mockLoginUsecase,
		RegisterUsecase:        mockRegisterUsecase,
		GetUserBalanceUsecase:  mockGetUserBalanceUsecase,
		WithdrawBalanceUsecase: mockWithdrawBalanceUsecase,
	}

	r.POST("/login", userController.Login)

	loginRequest := `{"login":"testuser","password":"password"}`
	mockLoginUsecase.EXPECT().Call(gomock.Any(), gomock.Any()).Return("test-token", nil)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(loginRequest)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test-token", w.Header().Get("Authorization"))
}

func TestUserController_Login_InvalidLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoginUsecase := mock_usecase.NewMockILoginUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	userController := &UserController{
		LoginUsecase: mockLoginUsecase,
	}

	r.POST("/login", userController.Login)

	loginRequest := `{"login":"testuser","password":"wrongpassword"}`
	mockLoginUsecase.EXPECT().Call(gomock.Any(), gomock.Any()).Return("", usecase.ErrInvalidLoginOrPassword)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(loginRequest)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserController_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegisterUsecase := mock_usecase.NewMockIRegisterUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	userController := &UserController{
		RegisterUsecase: mockRegisterUsecase,
	}

	r.POST("/register", userController.Register)

	registerRequest := `{"login":"newuser","password":"password"}`
	mockRegisterUsecase.EXPECT().Call(gomock.Any(), gomock.Any()).Return("test-token", nil)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte(registerRequest)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test-token", w.Header().Get("Authorization"))
}

func TestUserController_Register_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegisterUsecase := mock_usecase.NewMockIRegisterUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	userController := &UserController{
		RegisterUsecase: mockRegisterUsecase,
	}

	r.POST("/register", userController.Register)

	registerRequest := `{"login":"existinguser","password":"password"}`
	mockRegisterUsecase.EXPECT().Call(gomock.Any(), gomock.Any()).Return("", usecase.ErrUserAlreadyExists)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte(registerRequest)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUserController_GetUserBalance_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetUserBalanceUsecase := mock_usecase.NewMockIGetUserBalanceUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	userController := &UserController{
		GetUserBalanceUsecase: mockGetUserBalanceUsecase,
	}

	r.GET("/balance", userController.GetUserBalance)

	userBalance := &usecase.GetUserBalanceResult{
		Current:   entities.GDecimal(decimal.NewFromFloat(100.0)),
		Withdrawn: entities.GDecimal(decimal.NewFromFloat(50.0)),
	}

	mockGetUserBalanceUsecase.EXPECT().Call(gomock.Any(), gomock.Any()).Return(userBalance, nil)

	req := httptest.NewRequest(http.MethodGet, "/balance", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "100")
}

func TestUserController_WithdrawBalance_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawBalanceUsecase := mock_usecase.NewMockIWithdrawBalanceUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	userController := &UserController{
		WithdrawBalanceUsecase: mockWithdrawBalanceUsecase,
	}

	r.POST("/withdraw", userController.WithdrawBalance)

	withdrawRequest := `{"order":"79927398713","sum":50.0}`
	mockWithdrawBalanceUsecase.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer([]byte(withdrawRequest)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_WithdrawBalance_InsufficientBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWithdrawBalanceUsecase := mock_usecase.NewMockIWithdrawBalanceUsecase(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	userController := &UserController{
		WithdrawBalanceUsecase: mockWithdrawBalanceUsecase,
	}

	r.POST("/withdraw", userController.WithdrawBalance)

	withdrawRequest := `{"order":"12345678903","sum":150.0}`
	mockWithdrawBalanceUsecase.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(usecase.ErrInsufficientUserBalance)

	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer([]byte(withdrawRequest)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusPaymentRequired, w.Code)
}
