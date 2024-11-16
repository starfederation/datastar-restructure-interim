package datastar

import "fmt"

func (sse *ServerSentEventGenerator) ConsoleLog(msg string, opts ...ExecuteJSOption) error {
	call := fmt.Sprintf("console.log('%s')", msg)
	return sse.ExecuteJS(call, opts...)
}

func (sse *ServerSentEventGenerator) ConsoleLogf(format string, args ...any) error {
	return sse.ConsoleLog(fmt.Sprintf(format, args...))
}

func (sse *ServerSentEventGenerator) ConsoleError(err error, opts ...ExecuteJSOption) error {
	call := fmt.Sprintf("console.error('%s')", err.Error())
	return sse.ExecuteJS(call, opts...)
}
