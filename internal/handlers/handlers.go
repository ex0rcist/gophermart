package handlers

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

type Handlers struct {
	config  *config.Server
	engine  *gin.Engine
	storage *storage.PGXStorage
}

func NewHandlers(
	ctx context.Context,
	config *config.Server,
	storage *storage.PGXStorage,
) *Handlers {

	h := &Handlers{config: config, storage: storage}
	h.setupEngine()
	h.setupRoutes()

	return h
}

func (h *Handlers) Run() error {
	return h.engine.Run(h.config.Address)
}

func (h *Handlers) setupEngine() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.Logger

	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestsLogger())
	// engine.Use(middleware.JSONOnly())

	h.engine = engine
}

func (h *Handlers) setupRoutes() {
	h.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	public := h.engine.Group("/api/user")
	public.POST("/register", h.UserRegisterHandler())
	public.POST("/login", h.UserAuthHandler())

	private := h.engine.Group("/api/user")

	private.Use(middleware.Auth(h.storage, h.config.Secret))
	private.POST("/orders", h.OrderCreateHandler())
	private.GET("/orders", h.OrderListHandler())
	private.GET("/balance", h.UserGetBalanceHandler())

	private.POST("/balance/withdraw", h.UserWithdrawBalanceHandler())
	private.GET("/withdrawals", h.WithdrawalListHandler())
}

func (h *Handlers) getCurrentUser(c *gin.Context) *models.User {
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
