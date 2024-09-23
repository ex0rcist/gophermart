package entities

import (
	"errors"
	"fmt"
)

var (
	ErrBadAddressFormat = errors.New("bad net address format")
	ErrRecordNotFound   = errors.New("record not found")

	ErrUnexpected = errors.New("unexpected error")
)

func NewStackError(err error) error {
	return errors.New(err.Error())
}

func WrapError(prefix string, err error) error {
	return fmt.Errorf("%s: %s", prefix, err.Error())
}
