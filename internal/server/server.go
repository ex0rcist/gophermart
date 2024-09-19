package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/ex0rcist/gophermart/internal/accrual"
	"github.com/ex0rcist/gophermart/internal/config"
	"github.com/ex0rcist/gophermart/internal/handlers"
	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/storage"
)

type Server struct {
	config   *config.Config
	handlers *handlers.Handlers
	storage  *storage.PGXStorage
}

func New(config *config.Config) (*Server, error) {
	pgxStorage, err := storage.NewPGXStorage(config.DB)
	if err != nil {
		return nil, err
	}

	//ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	//defer cancel()

	ctx := context.Background()

	acr := accrual.NewService(ctx, &config.Accrual, pgxStorage)
	acr.Start()

	// todo: different context?
	hdl := handlers.NewHandlers(ctx, &config.Server, pgxStorage)

	return &Server{
		config:   config,
		handlers: hdl,
		storage:  pgxStorage,
	}, nil
}

func (s *Server) Run() error {
	logging.LogInfo(s.String())
	logging.LogInfo("server ready")

	return s.handlers.Run()
}

func (s *Server) String() string {
	str := []string{
		fmt.Sprintf("gophermart-address=%s", s.config.Server.Address),
		fmt.Sprintf("accrual-address=%s", s.config.Accrual.Address),
		fmt.Sprintf("database=%s", s.config.DB.DSN),
		fmt.Sprintf("secret=%s", s.config.Server.Secret),
	}

	return "server config: " + strings.Join(str, "; ")
}
