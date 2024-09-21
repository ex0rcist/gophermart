package httpbackend

import (
	"net/http"

	"github.com/ex0rcist/gophermart/internal/usecases"
	"github.com/gin-gonic/gin"
)

func (b *HTTPBackend) UserAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		const ep = "UserAuthHandler()" // error prefix
		ctx := c.Request.Context()

		var form = usecases.UserAuthForm{}
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// проверяем логин и пароль
		err := usecases.UserAuth(ctx, b.storage, form)
		if err != nil {
			if err == usecases.ErrInvalidLoginOrPassword {
				c.Status(http.StatusUnauthorized)
				return
			} else {
				handleInternalError(c, ctx, err, ep)
				return
			}
		}

		// создаем JWT токен
		token, err := usecases.UserCreateJWT(ctx, b.config.Secret, form.Login)
		if err != nil {
			handleInternalError(c, ctx, err, ep)
			return
		}

		c.Header("Authorization", token)
		c.Status(http.StatusOK)
	}
}

func (b *HTTPBackend) UserRegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		const ep = "UserRegisterHandler()" // error prefix
		ctx := c.Request.Context()

		var form = usecases.UserRegisterForm{}
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// регистрируем пользователя
		user, err := usecases.UserRegister(ctx, b.storage, form)

		switch {
		case err == usecases.ErrUserAlreadyExists:
			c.Status(http.StatusConflict)
			return
		case err != nil:
			handleInternalError(c, ctx, err, ep)
			return
		}

		// создаем JWT токен
		token, err := usecases.UserCreateJWT(ctx, b.config.Secret, user.Login)
		if err != nil {
			handleInternalError(c, ctx, err, ep)
			return
		}

		c.Header("Authorization", token)
		c.Status(http.StatusOK)
	}
}

func (b *HTTPBackend) UserGetBalanceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		const ep = "UserGetBalanceHandler()" // error prefix
		ctx := c.Request.Context()
		currentUser := b.getCurrentUser(c)

		// получаем баланс пользователя
		userBalance, err := usecases.UserGetBalance(ctx, b.storage, currentUser)
		if err != nil {
			handleInternalError(c, ctx, err, ep)
			return
		}

		c.JSON(http.StatusOK, userBalance)
	}
}

func (b *HTTPBackend) UserWithdrawBalanceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		const ep = "UserWithdrawBalanceHandler()" // error prefix
		ctx := c.Request.Context()
		currentUser := b.getCurrentUser(c)

		var form = usecases.UserWithdrawBalanceForm{UserID: currentUser.ID}
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// списываем баланс
		err := usecases.UserWithdrawBalance(ctx, b.storage, form)

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
