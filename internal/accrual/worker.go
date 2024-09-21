package accrual

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ex0rcist/gophermart/internal/logging"
)

type Worker struct {
	service *Service
}

func NewWorker(service *Service) Worker {
	return Worker{service: service}
}

func (w Worker) work() {
	for {
		// если сервис заблокирован (429), ждем до разблокировки
		if time.Now().Before(w.service.lockedUntil) {
			time.Sleep(time.Until(w.service.lockedUntil))
		}

		select {
		case <-w.service.ctx.Done():
			logging.LogDebug("accrual worker stopping")
			return
		default:
		}

		task := <-w.service.taskCh
		err := task.Handle()

		var cErr *ClientError
		if errors.As(err, &cErr) {
			// если получили 429 от клиента
			if cErr.HTTPStatus == http.StatusTooManyRequests && cErr.RetryAfter > 0 {
				logging.LogInfoCtx(w.service.ctx, "got 429, setting retry-after")

				// блокируем до указанного времени
				w.service.SetLockedUntil(time.Now().Add(cErr.RetryAfter))

				// возвращаем задачу обратно в канал
				w.service.Push(task)
			}

			logging.LogErrorCtx(w.service.ctx, fmt.Errorf("task failed: %w", err))
			continue
		}

		if err != nil {
			logging.LogErrorCtx(w.service.ctx, fmt.Errorf("task failed: %w", err))
		}
	}
}
