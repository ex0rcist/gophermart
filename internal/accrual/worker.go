package accrual

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ex0rcist/gophermart/internal/logging"
)

type Worker struct {
	service IService
}

func NewWorker(service IService) Worker {
	return Worker{service: service}
}

func (w Worker) Work(ctx context.Context, taskCh chan ITask) {
	for {
		select {
		case <-ctx.Done():
			logging.LogDebug("accrual worker stopping")
			return
		case <-time.After(time.Until(w.service.GetLockedUntil())):
			// сработает немедленно, если значение <= 0
			// т.е. продолжит выполнение итерации
		}

		task := <-taskCh
		err := task.Handle()

		var cErr *ClientError
		if errors.As(err, &cErr) {
			// если получили 429 от клиента
			if cErr.HTTPStatus == http.StatusTooManyRequests && cErr.RetryAfter > 0 {
				logging.LogInfoCtx(ctx, "got 429, setting retry-after")

				// блокируем до указанного времени
				w.service.SetLockedUntil(time.Now().Add(cErr.RetryAfter))

				// возвращаем задачу обратно в канал
				w.service.Push(task)
			}

			logging.LogErrorCtx(ctx, fmt.Errorf("task failed: %w", err))
			continue
		}

		if err != nil {
			logging.LogErrorCtx(ctx, fmt.Errorf("task failed: %w", err))
		}
	}
}
