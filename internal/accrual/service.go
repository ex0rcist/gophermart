package accrual

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/ex0rcist/gophermart/internal/config"
	"github.com/ex0rcist/gophermart/internal/domain"
	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/storage/repository"
	"github.com/ex0rcist/gophermart/internal/utils"
)

type Service struct {
	ctx context.Context

	client    *Client
	storage   *storage.PGXStorage
	userRepo  domain.IUserRepository
	orderRepo domain.IOrderRepository

	taskCh chan Task

	refillInterval time.Duration
	lockedUntil    time.Time
}

func NewService(ctx context.Context, config *config.Accrual, storage *storage.PGXStorage) *Service {
	return &Service{
		ctx: ctx,

		client:    NewClient(config.Address, utils.IntToDuration(config.Timeout)),
		storage:   storage,
		userRepo:  repository.NewUserRepository(storage.GetPool()),
		orderRepo: repository.NewOrderRepository(storage.GetPool()),

		taskCh: make(chan Task),

		refillInterval: config.RefillInterval,
	}
}

func (s *Service) Run() {
	logging.LogInfoF("starting accrual service, spawning %d workers", runtime.NumCPU())
	s.spawnWorkers()

	err := s.refillChannel()
	if err != nil {
		logging.LogError(err, "error refilling accrual channel")
	}

	go func() {
		ticker := time.NewTicker(s.refillInterval)
		defer ticker.Stop()

		for range ticker.C {
			// проверка на остановку приложения
			select {
			case <-s.ctx.Done():
				logging.LogInfo("accrual refilling stopped")
				return
			default:
			}

			err := s.refillChannel()
			if err != nil {
				logging.LogError(err, "err refilling channel")
			}
		}
	}()
}

func (s *Service) Push(t Task) {
	logging.LogDebugCtx(s.ctx, fmt.Sprintf("accrual service: pushing task with %s", t.order))

	s.taskCh <- t
}

func (s *Service) SetLockedUntil(lockedUntil time.Time) {
	if time.Now().After(lockedUntil) {
		return
	}

	s.lockedUntil = lockedUntil
}

func (s *Service) spawnWorkers() {
	for i := 0; i < runtime.NumCPU(); i++ {
		worker := NewWorker(s)
		go worker.work()
	}
}

func (s *Service) refillChannel() error {
	logging.LogDebug("refilling channel...")

	if len(s.taskCh) > 0 {
		logging.LogDebug("channel still has unread tasks, skipping")
		return nil
	}

	orders, err := s.orderRepo.OrderListForUpdate(s.ctx)
	if err != nil {
		return err
	}

	for _, o := range orders {
		t := Task{service: s, order: o}
		s.Push(t)
	}

	return nil
}
