package handlers

import (
	"net/http"

	"github.com/ex0rcist/gophermart/internal/usecases"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) UserGetBalanceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		const ep = "UserGetBalanceHandler()" // error prefix
		ctx := c.Request.Context()
		currentUser := h.getCurrentUser(c)

		// получаем баланс пользователя
		userBalance, err := usecases.UserGetBalance(ctx, h.storage, currentUser)
		if err != nil {
			handleInternalError(c, ctx, err, ep)
			return
		}

		c.JSON(http.StatusOK, userBalance)
	}
}
