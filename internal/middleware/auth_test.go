package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	mockStorage := new(storage.MockPGXStorage)
	secret := entities.Secret("test-secret")
	r.Use(Auth(mockStorage, secret))

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	mockStorage := new(storage.MockPGXStorage)
	secret := entities.Secret("test-secret")
	r.Use(Auth(mockStorage, secret))

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "invalid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	mockStorage := new(storage.MockPGXStorage)
	secret := entities.Secret("test-secret")

	expiredToken, _ := jwt.CreateJWT(secret, "test-login", -1*time.Minute)

	r.Use(Auth(mockStorage, secret))

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", expiredToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	mockStorage := new(storage.MockPGXStorage)
	secret := entities.Secret("test-secret")
	dur := 1 * time.Hour

	validToken, _ := jwt.CreateJWT(secret, "test-login", dur)

	user := &domain.User{ID: 1, Login: "test-login"}
	mockStorage.On("UserFindByLogin", mock.Anything, "test-login").Return(user, nil)

	r.Use(Auth(mockStorage, secret))
	r.GET("/test", func(c *gin.Context) {
		userFromContext, exists := c.Get(UserContextKey)

		assert.True(t, exists)
		assert.Equal(t, user.Login, userFromContext.(*domain.User).Login)

		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", validToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
