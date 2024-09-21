package httpbackend

import (
	"net/http"
	"strings"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/models"
	"github.com/ex0rcist/gophermart/internal/usecases"
	"github.com/gin-gonic/gin"
)

func (b *HTTPBackend) OrderCreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		const ep = "OrderCreateHandler()" // error prefix
		ctx := c.Request.Context()
		currentUser := b.getCurrentUser(c)

		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
			return
		}

		// подставляем номер заказа в форму напрямую из body, валидируем
		searchForm := usecases.OrderFindForm{Number: strings.TrimSpace(string(body))}
		if err := c.ShouldBindQuery(&searchForm); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// ищем заказ по номеру
		existingOrder, err := usecases.OrderFindByNumber(ctx, b.storage, searchForm)
		if err != nil && err != usecases.ErrOrderNotFound {
			handleInternalError(c, ctx, err, ep)
			return
		}

		if existingOrder != nil {
			if existingOrder.UserID == currentUser.ID {
				c.Status(http.StatusOK)
				return
			}

			c.Status(http.StatusConflict)
			return
		}

		createForm := usecases.OrderCreateForm{
			UserID: currentUser.ID,
			Number: searchForm.Number,
			Status: models.OrderStatusNew,
		}

		_, err = usecases.OrderCreate(ctx, b.storage, createForm)
		if err != nil {
			handleInternalError(c, ctx, err, ep)
			return
		}

		// приняли в обработку
		c.Status(http.StatusAccepted)
	}
}

func (b *HTTPBackend) OrderListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		const ep = "OrderListHandler()" // error prefix
		ctx := c.Request.Context()
		currentUser := b.getCurrentUser(c)

		orders, err := usecases.OrderList(ctx, b.storage, currentUser)
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
