package datastar

import "fmt"

func (sse *ServerSentEventGenerator) Consolef(mode ConsoleMode, messageFormat string, args ...any) error {
	message := fmt.Sprintf(messageFormat, args...)
	return sse.Console(mode, message)
}

func (sse *ServerSentEventGenerator) Console(mode ConsoleMode, message string, opts ...SSEEventOption) error {
	return sse.Send(
		EventTypeConsole,
		[]string{string(mode) + " " + message},
		opts...,
	)
}
