package handlers

import (
	"net/http"

	"github.com/ex0rcist/gophermart/internal/usecases"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) UserWithdrawBalanceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		const ep = "UserWithdrawBalanceHandler()" // error prefix
		ctx := c.Request.Context()
		currentUser := h.getCurrentUser(c)

		var form = usecases.UserWithdrawBalanceForm{UserID: currentUser.ID}
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// списываем баланс
		err := usecases.UserWithdrawBalance(ctx, h.storage, form)

		switch {
		case err == usecases.ErrInsufficientUserBalance:
			c.Status(http.StatusPaymentRequired)
			return
		case err != nil:
			handleInternalError(c, ctx, err, ep)
			return
		}

		c.Status(http.StatusOK)
	}
}
