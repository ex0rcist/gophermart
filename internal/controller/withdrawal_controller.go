package controller

import (
	"net/http"

	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/gin-gonic/gin"
)

type WithdrawalController struct {
	WithdrawalListUsecase domain.IWithdrawalListUsecase
}

func (ctrl *WithdrawalController) WithdrawalList(c *gin.Context) {
	const errorPrefix = "WithdrawalController -> WithdrawalList()"

	ctx := c.Request.Context()
	currentUser := getCurrentUser(c)

	wds, err := ctrl.WithdrawalListUsecase.Call(ctx, currentUser)
	if err != nil && err != entities.ErrRecordNotFound {
		handleInternalError(c, ctx, err, errorPrefix)
		return
	}

	if len(wds) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.JSON(http.StatusOK, wds)
	}
}