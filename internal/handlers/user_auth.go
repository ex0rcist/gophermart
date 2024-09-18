package handlers

import (
	"net/http"

	"github.com/ex0rcist/gophermart/internal/usecases"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) UserAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		const ep = "UserAuthHandler()" // error prefix
		ctx := c.Request.Context()

		var form = usecases.UserAuthForm{}
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// проверяем логин и пароль
		err := usecases.UserAuth(ctx, h.storage, form)
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
		token, err := usecases.UserCreateJWT(ctx, h.config.Secret, form.Login)
		if err != nil {
			handleInternalError(c, ctx, err, ep)
			return
		}

		c.Header("Authorization", token)
		c.Status(http.StatusOK)
	}
}
