package datastar

import "fmt"

func (sse *ServerSentEventGenerator) ConsoleError(err error, opts ...SSEEventOption) error {
	call := fmt.Sprintf("console.error('%s')", err.Error())
	return sse.Call(call, opts...)
}

// func (sse *ServerSentEventGenerator) ConsoleErrorf(format string, args ...any) error {
// 	return sse.Consolef(ConsoleModeError, format, args...)
// }

// func (sse *ServerSentEventGenerator) ConsoleError(message string, opts ...SSEEventOption) error {
// 	return sse.Console(ConsoleModeError, message, opts...)
// }

// func (sse *ServerSentEventGenerator) ConsoleWarnf(format string, args ...any) error {
// 	return sse.Consolef(ConsoleModeWarn, format, args...)
// }

// func (sse *ServerSentEventGenerator) ConsoleWarn(message string, opts ...SSEEventOption) error {
// 	return sse.Console(ConsoleModeWarn, message, opts...)
// }

// func (sse *ServerSentEventGenerator) ConsoleInfof(format string, args ...any) error {
// 	return sse.Consolef(ConsoleModeInfo, format, args...)
// }

// func (sse *ServerSentEventGenerator) ConsoleInfo(message string, opts ...SSEEventOption) error {
// 	return sse.Console(ConsoleModeInfo, message, opts...)
// }

// func (sse *ServerSentEventGenerator) ConsoleLogf(format string, args ...any) error {
// 	return sse.Consolef(ConsoleModeLog, format, args...)
// }

// func (sse *ServerSentEventGenerator) ConsoleLog(message string, opts ...SSEEventOption) error {
// 	return sse.Console(ConsoleModeLog, message, opts...)
// }

// func (sse *ServerSentEventGenerator) ConsoleDebugf(format string, args ...any) error {
// 	return sse.Consolef(ConsoleModeDebug, format, args...)
// }

// func (sse *ServerSentEventGenerator) ConsoleDebug(message string, opts ...SSEEventOption) error {
// 	return sse.Console(ConsoleModeDebug, message, opts...)
// }

// func (sse *ServerSentEventGenerator) ConsoleGroupf(labelFormat string, args ...any) error {
// 	label := fmt.Sprintf(labelFormat, args...)
// 	return sse.Consolef(ConsoleModeGroup, label)
// }

// func (sse *ServerSentEventGenerator) ConsoleGroup(label string, opts ...SSEEventOption) error {
// 	return sse.Console(ConsoleModeGroup, label, opts...)
// }

// func (sse *ServerSentEventGenerator) ConsoleGroupEnd(opts ...SSEEventOption) error {
// 	return sse.Console(ConsoleModeGroupEnd, "", opts...)
// }
