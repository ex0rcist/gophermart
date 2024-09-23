package controller

import (
	"net/http"

	"github.com/ex0rcist/gophermart/internal/config"
	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	LoginUsecase           domain.ILoginUsecase
	RegisterUsecase        domain.IRegisterUsecase
	GetUserBalanceUsecase  domain.IGetUserBalanceUsecase
	WithdrawBalanceUsecase domain.IWithdrawBalanceUsecase

	config *config.Server
}

func (uc *UserController) Login(c *gin.Context) {
	var form domain.LoginRequest
	ctx := c.Request.Context()
	const errorPrefix = "UserController -> Login()" // error prefix

	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// находим пользователя
	user, err := uc.LoginUsecase.GetUserByLogin(ctx, form)
	if err != nil {
		if err == domain.ErrInvalidLoginOrPassword {
			c.Status(http.StatusUnauthorized)
			return
		} else {
			handleInternalError(c, ctx, err, errorPrefix)
			return
		}
	}

	// создаем JWT токен
	token, err := uc.LoginUsecase.CreateAccessToken(user, uc.config.Secret, domain.LoginTokenLifetime)
	if err != nil {
		handleInternalError(c, ctx, err, errorPrefix)
		return
	}

	c.Header("Authorization", token)
	c.Status(http.StatusOK)
}

func (uc *UserController) Register(c *gin.Context) {
	var form domain.RegisterRequest
	ctx := c.Request.Context()
	const errorPrefix = "UserController -> Register()" // error prefix

	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// находим пользователя
	existingUser, err := uc.RegisterUsecase.GetUserByLogin(ctx, form)
	if err != nil && err != entities.ErrRecordNotFound {
		handleInternalError(c, ctx, err, errorPrefix)
		return
	}
	if existingUser != nil {
		c.Status(http.StatusConflict)
		return
	}

	// генерируем пароль
	form.Password, err = utils.HashPassword(form.Password)
	if err != nil {
		handleInternalError(c, ctx, err, errorPrefix)
		return
	}

	// создаем пользователя
	newUser, err := uc.RegisterUsecase.CreateUser(ctx, form.Login, form.Password)
	if err != nil {
		handleInternalError(c, ctx, err, errorPrefix)
		return
	}

	// создаем JWT токен
	token, err := uc.RegisterUsecase.CreateAccessToken(newUser, uc.config.Secret, domain.LoginTokenLifetime)
	if err != nil {
		handleInternalError(c, ctx, err, errorPrefix)
		return
	}

	c.Header("Authorization", token)
	c.Status(http.StatusOK)
}

func (uc *UserController) GetUserBalance(c *gin.Context) {
	const ep = "GetUserBalance()" // error prefix
	ctx := c.Request.Context()
	currentUser := getCurrentUser(c)

	bl, err := uc.GetUserBalanceUsecase.Fetch(ctx, currentUser.ID)
	if err != nil {
		handleInternalError(c, ctx, err, ep)
		return
	}

	c.JSON(http.StatusOK, bl)
}

func (uc *UserController) WithdrawBalance(c *gin.Context) {
	const ep = "WithdrawBalance()" // error prefix
	ctx := c.Request.Context()
	currentUser := getCurrentUser(c)

	var form = domain.WithdrawBalanceRequest{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	err := uc.WithdrawBalanceUsecase.Call(ctx, currentUser, form)
	switch {
	case err == domain.ErrInsufficientUserBalance:
		c.Status(http.StatusPaymentRequired)
		return
	case err != nil:
		handleInternalError(c, ctx, err, ep)
		return
	}

	c.Status(http.StatusOK)
}
