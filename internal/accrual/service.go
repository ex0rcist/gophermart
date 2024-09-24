package accrual

import (
	"context"
	"runtime"
	"time"

	"github.com/ex0rcist/gophermart/internal/config"
	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/storage"
	"github.com/ex0rcist/gophermart/internal/storage/repository"
)

type IService interface {
	Push(t ITask)
	Run() error
	SetLockedUntil(lockedUntil time.Time)
	GetLockedUntil() time.Time
}

type Service struct {
	ctx context.Context

	client    IClient
	storage   storage.IPGXStorage
	userRepo  repository.IUserRepository
	orderRepo repository.IOrderRepository

	taskCh chan ITask

	contextTimeout time.Duration
	refillInterval time.Duration

	lockedUntil time.Time
}

func NewService(
	ctx context.Context,
	config *config.Accrual,
	client IClient,
	storage storage.IPGXStorage,
	userRepo repository.IUserRepository,
	orderRepo repository.IOrderRepository,
) *Service {
	if client == nil {
		client = NewClient(config.Address, config.Timeout)
	}

	if userRepo == nil {
		userRepo = repository.NewUserRepository(storage.GetPool())
	}

	if orderRepo == nil {
		orderRepo = repository.NewOrderRepository(storage.GetPool())
	}

	return &Service{
		ctx: ctx,

		client:    client,
		storage:   storage,
		userRepo:  userRepo,
		orderRepo: orderRepo,

		taskCh: make(chan ITask),

		contextTimeout: config.Timeout,
		refillInterval: config.RefillInterval,
	}
}

func (s *Service) Run() error {
	logging.LogInfoF("starting accrual service, spawning %d workers", runtime.NumCPU())
	s.spawnWorkers()

	err := s.refillChannel()
	if err != nil {
		logging.LogError(err, "error refilling accrual channel")
		return err
	}

	go func() {
		for {
			select {
			case <-s.ctx.Done(): // проверка на остановку приложения
				logging.LogInfo("accrual refilling stopped")
				return
			case <-time.After(s.refillInterval): // если не остановлено, выполняем заправку канала
				err := s.refillChannel()
				if err != nil {
					logging.LogError(err, "err refilling channel")
				}
			}
		}
	}()

	return nil
}

func (s *Service) Push(t ITask) {
	s.taskCh <- t
}

func (s *Service) SetLockedUntil(lockedUntil time.Time) {
	if time.Now().After(lockedUntil) {
		return
	}

	s.lockedUntil = lockedUntil
}

func (s *Service) GetLockedUntil() time.Time {
	return s.lockedUntil
}

func (s *Service) spawnWorkers() {
	for i := 0; i < runtime.NumCPU(); i++ {
		worker := NewWorker(s)
		go worker.Work(s.ctx, s.taskCh)
	}
}

func (s *Service) refillChannel() error {
	logging.LogDebug("refilling channel...")

	// частично исключаем дублирование задач
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
