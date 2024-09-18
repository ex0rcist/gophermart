package entities

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrBadAddressFormat = errors.New("bad net address format")
	ErrRecordNotFound   = errors.New("record not found")

	ErrUnexpected = errors.New("unexpected error")
)

func NewStackError(err error) error {
	return errors.New(err.Error())
}

var _ error = (*RetriableError)(nil)

type RetriableError struct {
	Err        error
	RetryAfter time.Duration
}

func (e RetriableError) Error() string {
	return fmt.Sprintf("%s (retry after %v)", e.Err.Error(), e.RetryAfter)
}

func WrapError(prefix string, err error) error {
	return fmt.Errorf("%s: %s", prefix, err.Error())
}
