package sdk

import "fmt"

func (sse *ServerSentEventGenerator) ConsoleErrorf(format string, args ...any) error {
	return sse.Consolef(ConsoleLogModeError, format, args...)
}

func (sse *ServerSentEventGenerator) ConsoleError(message string) error {
	return sse.Console(ConsoleLogModeError, message)
}

func (sse *ServerSentEventGenerator) ConsoleWarnf(format string, args ...any) error {
	return sse.Consolef(ConsoleLogModeWarn, format, args...)
}

func (sse *ServerSentEventGenerator) ConsoleWarn(message string) error {
	return sse.Console(ConsoleLogModeWarn, message)
}

func (sse *ServerSentEventGenerator) ConsoleInfof(format string, args ...any) error {
	return sse.Consolef(ConsoleLogModeInfo, format, args...)
}

func (sse *ServerSentEventGenerator) ConsoleInfo(message string) error {
	return sse.Console(ConsoleLogModeInfo, message)
}

func (sse *ServerSentEventGenerator) ConsoleLogf(format string, args ...any) error {
	return sse.Consolef(ConsoleLogModeLog, format, args...)
}

func (sse *ServerSentEventGenerator) ConsoleDebugf(format string, args ...any) error {
	return sse.Consolef(ConsoleLogModeDebug, format, args...)
}

func (sse *ServerSentEventGenerator) ConsoleDebug(message string) error {
	return sse.Console(ConsoleLogModeDebug, message)
}

func (sse *ServerSentEventGenerator) ConsoleGroupf(labelFormat string, args ...any) error {
	label := fmt.Sprintf(labelFormat, args...)
	return sse.ConsoleGroup(label)
}

func (sse *ServerSentEventGenerator) ConsoleGroup(label string) error {
	return sse.Console(ConsoleLogModeGroup, label)
}

func (sse *ServerSentEventGenerator) ConsoleGroupEnd() error {
	return sse.Console(ConsoleLogModeGroupEnd, "")
}
