package models

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() error {
	return errors.New(fmt.Sprintf("code: %d, message: %s", e.Code, e.Message))
}
