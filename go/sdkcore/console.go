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

func (sse *ServerSentEventsHandler) Console(Mode ConsoleLogMode, messageFormat string, args ...any) {
	message := fmt.Sprintf(messageFormat, args...)
	sse.send(
		EventTypeConsole,
		[]string{
			fmt.Sprintf("%s %s", Mode, message),
		},
	)
}
