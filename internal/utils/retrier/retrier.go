package retrier

import (
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/ex0rcist/gophermart/internal/logging"
)

type Retrier struct {
	payloadFn retry.RetryableFunc
	retryIfFn retry.RetryIfFunc
	delays    []time.Duration
}

var _ error = (*RetriableError)(nil)

type RetriableError struct {
	Err        error
	RetryAfter time.Duration
}

func (e RetriableError) Error() string {
	return fmt.Sprintf("%s (retry after %v)", e.Err.Error(), e.RetryAfter)
}

func (r Retrier) Run() error {
	return retry.Do(
		r.payloadFn,
		retry.RetryIf(r.retryIfFn),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			logging.LogWarnF("will retry after %v", r.delays[n])
			return r.delays[n]
		}),
		retry.Attempts(uint(len(r.delays))+1),
	)
}

func NewRetrier(payloadFn func() error, retryIfFn func(err error) bool, delays []time.Duration) Retrier {
	return Retrier{
		payloadFn: payloadFn,
		retryIfFn: retryIfFn,
		delays:    delays,
	}
}
