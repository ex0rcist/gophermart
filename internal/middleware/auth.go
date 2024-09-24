package middleware

import (
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/storage/repository"
	"github.com/ex0rcist/gophermart/pkg/jwt"
	"github.com/gin-gonic/gin"

	"net/http"
	"time"
)

const UserContextKey = "currentUser"

func Auth(
	repo repository.IUserRepository,
	key entities.Secret,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		token := c.Request.Header.Get("Authorization")
		if token == "" {
			logging.LogInfoCtx(ctx, "auth: no token provided")
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		login, ttl, err := jwt.ParseJWT(key, token)
		if err != nil {
			logging.LogErrorCtx(ctx, err, "auth: jwt parsing err")
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		if time.Now().After(ttl) {
			logging.LogInfoCtx(ctx, "auth: jwt token expired")
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		user, err := repo.UserFindByLogin(ctx, login)
		if err != nil {
			if err == storage.ErrRecordNotFound {
				logging.LogInfoCtx(ctx, "auth: login not found")
				c.Status(http.StatusUnauthorized)
			} else {
				logging.LogErrorCtx(ctx, err, "auth: middleware err")
				c.Status(http.StatusInternalServerError)
			}

			c.Abort()
			return
		}

		c.Set(UserContextKey, user)
		c.Next()
	}
}
