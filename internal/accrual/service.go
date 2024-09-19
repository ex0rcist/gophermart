package accrual

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/ex0rcist/gophermart/internal/config"
	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/utils"
)

type Service struct {
	ctx context.Context

	client  *Client
	storage *storage.PGXStorage

	taskCh chan Task

	refillInterval time.Duration
	lockedUntil    time.Time
}

func NewService(ctx context.Context, config *config.Accrual, storage *storage.PGXStorage) *Service {
	return &Service{
		ctx:            ctx,
		client:         NewClient(config.Address, utils.IntToDuration(config.Timeout)),
		storage:        storage,
		refillInterval: config.RefillInterval,
		taskCh:         make(chan Task),
	}
}

func (s *Service) Start() {
	s.spawnWorkers()
	err := s.refillChannel()
	if err != nil {
		logging.LogErrorCtx(s.ctx, err)
	}

	go func() {
		ticker := time.NewTicker(s.refillInterval)
		defer ticker.Stop()

		for range ticker.C {
			err := s.refillChannel()
			if err != nil {
				logging.LogErrorCtx(s.ctx, err)
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
		logging.LogDebugCtx(s.ctx, fmt.Sprintf("accrual service: spawning worker %d", i))

		worker := NewWorker(s)
		go worker.work()
	}
}

func (s *Service) refillChannel() error {
	logging.LogDebugCtx(s.ctx, "refilling channel...")

	orders, err := s.storage.OrderListForUpdate(s.ctx)
	if err != nil {
		return err
	}

	for _, o := range orders {
		t := Task{service: s, order: o}
		s.Push(t)
	}

	return nil
}
