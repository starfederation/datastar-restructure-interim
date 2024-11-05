package sdk

import "fmt"

type ConsoleLogMode string

const (
	ConsoleLogModeLog      ConsoleLogMode = "log"
	ConsoleLogModeError    ConsoleLogMode = "error"
	ConsoleLogModeWarn     ConsoleLogMode = "warn"
	ConsoleLogModeInfo     ConsoleLogMode = "info"
	ConsoleLogModeDebug    ConsoleLogMode = "debug"
	ConsoleLogModeGroup    ConsoleLogMode = "group"
	ConsoleLogModeGroupEnd ConsoleLogMode = "groupEnd"
)

func (sse *ServerSentEventGenerator) Consolef(mode ConsoleLogMode, messageFormat string, args ...any) error {
	message := fmt.Sprintf(messageFormat, args...)
	return sse.Console(mode, message)
}

func (sse *ServerSentEventGenerator) Console(mode ConsoleLogMode, message string) error {
	return sse.send(
		EventTypeConsole,
		[]string{string(mode) + message},
	)
}
