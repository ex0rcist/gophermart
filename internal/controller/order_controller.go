package controller

import (
	"net/http"
	"strings"

	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/usecase"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	OrderCreateUsecase usecase.IOrderCreateUsecase
	OrderListUsecase   usecase.IOrderListUsecase
}

func (ctrl *OrderController) CreateOrder(c *gin.Context) {
	const errorPrefix = "OrderController -> CreateOrder()"
	ctx := c.Request.Context()
	currentUser := getCurrentUser(c)

	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
		return
	}
	number := strings.TrimSpace(string(body))

	_, err = ctrl.OrderCreateUsecase.Create(ctx, currentUser, number)
	if err != nil {
		switch {
		case err == usecase.ErrInvalidOrderNumber:
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		case err == usecase.ErrOrderAlreadyRegistered:
			c.Status(http.StatusOK)
			return
		case err == usecase.ErrOrderConflict:
			c.Status(http.StatusConflict)
			return
		default:
			handleInternalError(c, ctx, err, errorPrefix)
			return
		}
	}

	// приняли в обработку
	c.Status(http.StatusAccepted)
}
func (ctrl *OrderController) OrderList(c *gin.Context) {
	const errorPrefix = "OrderController -> OrderList()"
	ctx := c.Request.Context()
	currentUser := getCurrentUser(c)

	orders, err := ctrl.OrderListUsecase.Call(ctx, currentUser)
	if err != nil && err != storage.ErrRecordNotFound {
		handleInternalError(c, ctx, err, errorPrefix)
		return
	}

	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, orders)
	}
}
