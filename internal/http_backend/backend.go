package httpbackend

import (
	"context"
	"net/http"
	"time"

	"github.com/ex0rcist/gophermart/internal/config"
	"github.com/ex0rcist/gophermart/internal/controller"
	"github.com/ex0rcist/gophermart/internal/middleware"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/storage/repository"
	"github.com/ex0rcist/gophermart/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const defaultTimeout = 5 * time.Second

type HTTPBackend struct {
	config     *config.Server
	httpServer *http.Server
	router     *gin.Engine
	storage    storage.IPGXStorage
}

func NewHTTPBackend(
	ctx context.Context,
	config *config.Server,
	storage storage.IPGXStorage,
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

	publicRouter := b.router.Group("")

	privateRouter := b.router.Group("")
	privateRouter.Use(middleware.Auth(
		repository.NewUserRepository(b.storage.GetPool()),
		b.config.Secret,
	))

	b.setupUserController(publicRouter, privateRouter)
	b.setupOrderController(publicRouter, privateRouter)
	b.setupWithdrawalController(publicRouter, privateRouter)

}

func (b *HTTPBackend) setupUserController(publicRouter *gin.RouterGroup, privateRouter *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(b.storage.GetPool())
	wdrwRepo := repository.NewWithdrawalRepository(b.storage.GetPool())

	ctrl := &controller.UserController{
		LoginUsecase:           usecase.NewLoginUsecase(b.storage, userRepo, b.config.Secret, defaultTimeout),
		RegisterUsecase:        usecase.NewRegisterUsecase(b.storage, userRepo, b.config.Secret, defaultTimeout),
		GetUserBalanceUsecase:  usecase.NewGetUserBalanceUsecase(b.storage, userRepo, defaultTimeout),
		WithdrawBalanceUsecase: usecase.NewWithdrawBalanceUsecase(b.storage, userRepo, wdrwRepo, defaultTimeout),
	}

	publicRouter.POST("/api/user/register", ctrl.Register)
	publicRouter.POST("/api/user/login", ctrl.Login)

	privateRouter.GET("/api/user/balance", ctrl.GetUserBalance)
	privateRouter.POST("/api/user/balance/withdraw", ctrl.WithdrawBalance)
}

func (b *HTTPBackend) setupOrderController(_ *gin.RouterGroup, privateRouter *gin.RouterGroup) {
	repo := repository.NewOrderRepository(b.storage.GetPool())

	ctrl := &controller.OrderController{
		OrderCreateUsecase: usecase.NewOrderCreateUsecase(b.storage, repo, defaultTimeout),
		OrderListUsecase:   usecase.NewOrderListUsecase(b.storage, repo, defaultTimeout),
	}

	privateRouter.POST("/api/user/orders", ctrl.CreateOrder)
	privateRouter.GET("/api/user/orders", ctrl.OrderList)
}

func (b *HTTPBackend) setupWithdrawalController(_ *gin.RouterGroup, privateRouter *gin.RouterGroup) {
	repo := repository.NewWithdrawalRepository(b.storage.GetPool())

	ctrl := &controller.WithdrawalController{
		WithdrawalListUsecase: usecase.NewWithdrawalListUsecase(b.storage, repo, defaultTimeout),
	}

	privateRouter.GET("/api/user/withdrawals", ctrl.WithdrawalList)
}

func (b *HTTPBackend) setupServer() {
	b.httpServer = &http.Server{
		Addr:    b.config.Address,
		Handler: b.router.Handler(),
	}
}
