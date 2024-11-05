package sdk

import "fmt"

func (sse *ServerSentEventGenerator) ConsoleErrorf(format string, args ...any) error {
	return sse.Consolef(ConsoleModeError, format, args...)
}

func (sse *ServerSentEventGenerator) ConsoleError(message string) error {
	return sse.Console(ConsoleModeError, message)
}

func (sse *ServerSentEventGenerator) ConsoleWarnf(format string, args ...any) error {
	return sse.Consolef(ConsoleModeWarn, format, args...)
}

func (sse *ServerSentEventGenerator) ConsoleWarn(message string) error {
	return sse.Console(ConsoleModeWarn, message)
}

func (sse *ServerSentEventGenerator) ConsoleInfof(format string, args ...any) error {
	return sse.Consolef(ConsoleModeInfo, format, args...)
}

func (sse *ServerSentEventGenerator) ConsoleInfo(message string) error {
	return sse.Console(ConsoleModeInfo, message)
}

func (sse *ServerSentEventGenerator) ConsoleLogf(format string, args ...any) error {
	return sse.Consolef(ConsoleModeLog, format, args...)
}

func (sse *ServerSentEventGenerator) ConsoleDebugf(format string, args ...any) error {
	return sse.Consolef(ConsoleModeDebug, format, args...)
}

func (sse *ServerSentEventGenerator) ConsoleDebug(message string) error {
	return sse.Console(ConsoleModeDebug, message)
}

func (sse *ServerSentEventGenerator) ConsoleGroupf(labelFormat string, args ...any) error {
	label := fmt.Sprintf(labelFormat, args...)
	return sse.Consolef(ConsoleModeGroup, label)
}

func (sse *ServerSentEventGenerator) ConsoleGroup(label string) error {
	return sse.Console(ConsoleModeGroup, label)
}

func (sse *ServerSentEventGenerator) ConsoleGroupEnd() error {
	return sse.Console(ConsoleModeGroupEnd, "")
}
