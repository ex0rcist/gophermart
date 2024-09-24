package app

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ex0rcist/gophermart/internal/accrual"
	"github.com/ex0rcist/gophermart/internal/config"
	httpbackend "github.com/ex0rcist/gophermart/internal/http_backend"
	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/storage"
)

type App struct {
	ctx         context.Context
	config      *config.Config
	storage     storage.IPGXStorage
	httpBackend httpbackend.IHTTPBackend
	accrService accrual.IService
	cancelFunc  context.CancelFunc
}

func New(
	config *config.Config,
	pgxStorage storage.IPGXStorage,
	accrService accrual.IService,
	httpBackend httpbackend.IHTTPBackend,
) (*App, error) {
	var err error
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	if pgxStorage == nil {
		pgxStorage, err = storage.NewPGXStorage(config.DB, nil, true)
		if err != nil {
			return nil, fmt.Errorf("NewPGXStorage() failed: %w", err)
		}
	}

	if accrService == nil {
		accrService = accrual.NewService(ctx, &config.Accrual, nil, pgxStorage, nil, nil)
	}

	if httpBackend == nil {
		httpBackend = httpbackend.NewHTTPBackend(ctx, &config.Server, pgxStorage)
	}

	return &App{
		ctx:         ctx,
		config:      config,
		storage:     pgxStorage,
		httpBackend: httpBackend,
		accrService: accrService,
		cancelFunc:  cancel,
	}, nil
}

func (a *App) Run() error {
	logging.LogInfo(a.String())
	logging.LogInfo("app ready")

	// стартуем http backend
	go func() {
		err := a.httpBackend.Run()
		if err != nil && err != http.ErrServerClosed {
			logging.LogError(err, "httpServer error")
			a.cancelFunc()
		}
	}()

	// стартуем интеграцию с accrual
	// NB: остановка a.accrService не требуется, т.к. он слушает a.ctx
	go func() {
		err := a.accrService.Run()
		if err != nil {
			logging.LogError(err, "accrualService error")
			a.cancelFunc()
		}
	}()

	<-a.ctx.Done() // ждем сигнал от NotifyContext

	// останавливаем сервер
	logging.LogInfo("stopping server... ")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := a.httpBackend.Shutdown(ctx); err != nil {
		logging.LogError(err, "err stopping server")
	}

	// ждем остановку http.Server
	<-ctx.Done()
	logging.LogInfo("server stopped")

	// закрываем коннекты к БД
	a.storage.Close()
	logging.LogInfo("storage closed")

	return nil
}

func (a *App) String() string {
	str := []string{
		fmt.Sprintf("gophermart-address=%s", a.config.Server.Address),
		fmt.Sprintf("accrual-address=%s", a.config.Accrual.Address),
		fmt.Sprintf("database=%s", a.config.DB.DSN),
		fmt.Sprintf("secret=%s", a.config.Server.Secret),
	}

	return "app config: " + strings.Join(str, "; ")
}
