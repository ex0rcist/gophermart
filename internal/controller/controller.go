package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/middleware"
	"github.com/gin-gonic/gin"
)

func handleInternalError(c *gin.Context, ctx context.Context, err error, ep string) {
	logging.LogErrorCtx(ctx, wrap(ep, err))
	c.Status(http.StatusInternalServerError)
}

func wrap(prefix string, err error) error {
	return fmt.Errorf("%s: %s", prefix, err.Error())
}

func getCurrentUser(c *gin.Context) *domain.User {
	user, exists := c.Get(middleware.UserContextKey)
	if !exists {
		return nil
	}

	currentUser := user.(*domain.User)
	return currentUser
}
