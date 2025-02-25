package controller

import (
	"net/http"

	"github.com/ex0rcist/gophermart/internal/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	LoginUsecase           usecase.ILoginUsecase
	RegisterUsecase        usecase.IRegisterUsecase
	GetUserBalanceUsecase  usecase.IGetUserBalanceUsecase
	WithdrawBalanceUsecase usecase.IWithdrawBalanceUsecase
}

func (ctrl *UserController) Login(c *gin.Context) {
	const errorPrefix = "UserController -> Login()"
	var form usecase.LoginRequest
	ctx := c.Request.Context()

	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ctrl.LoginUsecase.Call(ctx, form)
	if err != nil {
		if err == usecase.ErrInvalidLoginOrPassword {
			c.Status(http.StatusUnauthorized)
			return
		}

		handleInternalError(c, ctx, err, errorPrefix)
		return
	}

	c.Header("Authorization", token)
	c.Status(http.StatusOK)
}

func (ctrl *UserController) Register(c *gin.Context) {
	const errorPrefix = "UserController -> Register()"
	var form usecase.RegisterRequest
	ctx := c.Request.Context()

	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ctrl.RegisterUsecase.Call(ctx, form)
	if err != nil {
		if err == usecase.ErrUserAlreadyExists {
			c.Status(http.StatusConflict)
			return
		}

		handleInternalError(c, ctx, err, errorPrefix)
		return
	}

	c.Header("Authorization", token)
	c.Status(http.StatusOK)
}

func (ctrl *UserController) GetUserBalance(c *gin.Context) {
	const errorPrefix = "UserController -> GetUserBalance()"
	ctx := c.Request.Context()
	currentUser := getCurrentUser(c)

	bl, err := ctrl.GetUserBalanceUsecase.Call(ctx, currentUser)
	if err != nil {
		handleInternalError(c, ctx, err, errorPrefix)
		return
	}

	c.JSON(http.StatusOK, bl)
}

func (ctrl *UserController) WithdrawBalance(c *gin.Context) {
	const ep = "UserController -> WithdrawBalance()"
	ctx := c.Request.Context()
	currentUser := getCurrentUser(c)

	var form = usecase.WithdrawBalanceRequest{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.WithdrawBalanceUsecase.Call(ctx, currentUser, form)
	switch {
	case err == usecase.ErrInvalidOrderNumber:
		// тут уже вторая валидация на Luhn,
		// первая в form (для общего развития)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	case err == usecase.ErrInsufficientUserBalance:
		c.Status(http.StatusPaymentRequired)
		return
	case err != nil:
		handleInternalError(c, ctx, err, ep)
		return
	}

	c.Status(http.StatusOK)
}
