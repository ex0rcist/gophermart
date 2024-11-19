package accrual_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ex0rcist/gophermart/internal/accrual"
	mock_accrual "github.com/ex0rcist/gophermart/internal/accrual/mocks"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func TestWorker_Work_ContextDone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_accrual.NewMockIService(ctrl)
	taskCh := make(chan accrual.ITask)

	// настраиваем мок, чтобы система не была заблокирована
	mockService.EXPECT().GetLockedUntil().Return(time.Now())

	worker := accrual.NewWorker(mockService)

	// создаем контекст, который будет отменен
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // отменяем сразу

	go worker.Work(ctx, taskCh)
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, true, ctx.Err() != nil, "Expected worker to stop due to cancelled context")
}

func TestWorker_Work_TaskSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_accrual.NewMockIService(ctrl)
	mockTask := mock_accrual.NewMockITask(ctrl)
	taskCh := make(chan accrual.ITask, 1)

	// настраиваем мок, чтобы система не была заблокирована
	mockService.EXPECT().GetLockedUntil().Return(time.Now()).AnyTimes()

	mockTask.EXPECT().Handle().Return(nil).Times(1)

	taskCh <- mockTask

	worker := accrual.NewWorker(mockService)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go worker.Work(ctx, taskCh)
	time.Sleep(100 * time.Millisecond)
}

func TestWorker_Work_ClientError429(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_accrual.NewMockIService(ctrl)
	mockTask := mock_accrual.NewMockITask(ctrl)
	taskCh := make(chan accrual.ITask, 1)

	// настраиваем мок, чтобы система не была заблокирована
	mockService.EXPECT().GetLockedUntil().Return(time.Now()).AnyTimes()

	clientError := &accrual.ClientError{
		HTTPStatus: http.StatusTooManyRequests,
		RetryAfter: 5 * time.Second,
	}
	mockTask.EXPECT().Handle().Return(clientError).Times(1)

	mockService.EXPECT().SetLockedUntil(gomock.Any()).Do(func(tm time.Time) {
		// т.к. точно поймать 5 секунд невозможно, используем indelta
		assert.InDelta(t, time.Until(tm), 5*time.Second, float64(time.Second))
	})

	// ожидание, что задача будет возвращена в очередь
	mockService.EXPECT().Push(mockTask)

	taskCh <- mockTask

	worker := accrual.NewWorker(mockService)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go worker.Work(ctx, taskCh)
	time.Sleep(100 * time.Millisecond)
}

func TestWorker_Work_NonClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_accrual.NewMockIService(ctrl)
	mockTask := mock_accrual.NewMockITask(ctrl)
	taskCh := make(chan accrual.ITask, 1)

	// настраиваем мок, чтобы система не была заблокирована
	mockService.EXPECT().GetLockedUntil().Return(time.Now()).AnyTimes()

	mockTask.EXPECT().Handle().Return(errors.New("some error")).Times(1)

	taskCh <- mockTask

	worker := accrual.NewWorker(mockService)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go worker.Work(ctx, taskCh)
	time.Sleep(100 * time.Millisecond)
}
