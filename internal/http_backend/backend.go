package httpbackend

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ex0rcist/gophermart/internal/config"
	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/middleware"
	"github.com/ex0rcist/gophermart/internal/models"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HTTPBackend struct {
	config     *config.Server
	httpServer *http.Server
	router     *gin.Engine
	storage    *storage.PGXStorage
}

func NewHTTPBackend(
	ctx context.Context,
	config *config.Server,
	storage *storage.PGXStorage,
) *HTTPBackend {
	b := &HTTPBackend{config: config, storage: storage}
	b.setupRouter()
	b.setupRoutes()
	b.setupServer()

	return b
}

func (b *HTTPBackend) Run() error {
	return b.httpServer.ListenAndServe()
}

func (b *HTTPBackend) Shutdown(ctx context.Context) error {
	return b.httpServer.Shutdown(ctx)
}

func (b *HTTPBackend) setupRouter() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.Logger

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.RequestsLogger())

	b.router = router
}

func (b *HTTPBackend) setupRoutes() {
	b.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	public := b.router.Group("/api/user")
	public.POST("/register", b.UserRegisterHandler())
	public.POST("/login", b.UserAuthHandler())

	private := b.router.Group("/api/user")

	private.Use(middleware.Auth(b.storage, b.config.Secret))
	private.POST("/orders", b.OrderCreateHandler())
	private.GET("/orders", b.OrderListHandler())
	private.GET("/balance", b.UserGetBalanceHandler())

	private.POST("/balance/withdraw", b.UserWithdrawBalanceHandler())
	private.GET("/withdrawals", b.WithdrawalListHandler())
}

func (b *HTTPBackend) setupServer() {
	b.httpServer = &http.Server{
		Addr:    b.config.Address,
		Handler: b.router.Handler(),
	}
}

func (b *HTTPBackend) getCurrentUser(c *gin.Context) *models.User {
	user, exists := c.Get(middleware.UserContextKey)
	if !exists {
		return nil
	}

	currentUser := user.(*models.User)
	return currentUser
}

func handleInternalError(c *gin.Context, ctx context.Context, err error, ep string) {
	logging.LogErrorCtx(ctx, wrap(ep, err))
	c.Status(http.StatusInternalServerError)
}

func wrap(prefix string, err error) error {
	return fmt.Errorf("%s: %s", prefix, err.Error())
}
