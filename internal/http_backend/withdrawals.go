package httpbackend

import (
	"net/http"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/usecases"
	"github.com/gin-gonic/gin"
)

func (b *HTTPBackend) WithdrawalListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		const ep = "UserGetWithdrawalsHandler()" // error prefix
		ctx := c.Request.Context()
		currentUser := b.getCurrentUser(c)

		wds, err := usecases.WithdrawalList(ctx, b.storage, currentUser)
		if err != nil && err != entities.ErrRecordNotFound {
			handleInternalError(c, ctx, err, ep)
			return
		}

		if len(wds) == 0 {
			c.Status(http.StatusNoContent)
		} else {
			c.JSON(http.StatusOK, wds)
		}
	}
}
