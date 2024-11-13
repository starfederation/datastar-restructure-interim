package datastar

import (
	"errors"
)

var (
	ErrEventTypeError = errors.New("event type is required")

	newLineBuf       = []byte("\n")
	doubleNewLineBuf = []byte("\n\n")
)
