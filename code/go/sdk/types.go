package sdk

import (
	"errors"
	"fmt"
)

var (
	ErrEventTypeError = fmt.Errorf("event type is required")
	ErrMissingInput   = errors.New("missing datastar input")

	newLineBuf       = []byte("\n")
	doubleNewLineBuf = []byte("\n\n")
)

type EventType []byte

var (
	EventTypeFragment = EventType("datastar-fragment")
	EventTypeSignal   = EventType("datastar-signal")
	EventTypeDelete   = EventType("datastar-delete")
	EventTypeRedirect = EventType("datastar-redirect")
	EventTypeConsole  = EventType("datastar-console")
)
