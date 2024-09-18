package handlers

import (
	"net/http"

	"github.com/ex0rcist/gophermart/internal/usecases"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) UserRegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		const ep = "UserRegisterHandler()" // error prefix
		ctx := c.Request.Context()

		var form = usecases.UserRegisterForm{}
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// регистрируем пользователя
		user, err := usecases.UserRegister(ctx, h.storage, form)

		switch {
		case err == usecases.ErrUserAlreadyExists:
			c.Status(http.StatusConflict)
			return
		case err != nil:
			handleInternalError(c, ctx, err, ep)
			return
		}

		// создаем JWT токен
		token, err := usecases.UserCreateJWT(ctx, h.config.Secret, user.Login)
		if err != nil {
			handleInternalError(c, ctx, err, ep)
			return
		}

		c.Header("Authorization", token)
		c.Status(http.StatusOK)
	}
}
