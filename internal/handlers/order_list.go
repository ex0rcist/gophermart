package handlers

import (
	"net/http"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/usecases"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) OrderListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		const ep = "OrderListHandler()" // error prefix
		ctx := c.Request.Context()
		currentUser := h.getCurrentUser(c)

		orders, err := usecases.OrderList(ctx, h.storage, currentUser)
		if err != nil && err != entities.ErrRecordNotFound {
			handleInternalError(c, ctx, err, ep)
			return
		}

		if len(orders) == 0 {
			c.Status(http.StatusNoContent)
		} else {
			c.JSON(http.StatusOK, orders)
		}
	}
}
