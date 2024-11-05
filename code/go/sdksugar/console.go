package sdk

import sdkcore "github.com/starfederation/datastar/go/sdk"

func (sse *ServerSentEventsHandler) ConsoleError(format string, args ...interface{}) {
	sse.Console(sdkcore.ConsoleLogModeError, format, args...)
}

func (sse *ServerSentEventsHandler) ConsoleWarn(format string, args ...interface{}) {
	sse.Console(sdkcore.ConsoleLogModeWarn, format, args...)
}

func (sse *ServerSentEventsHandler) ConsoleInfo(format string, args ...interface{}) {
	sse.Console(sdkcore.ConsoleLogModeInfo, format, args...)
}

func (sse *ServerSentEventsHandler) ConsoleDebug(format string, args ...interface{}) {
	sse.Console(sdkcore.ConsoleLogModeDebug, format, args...)
}

func (sse *ServerSentEventsHandler) ConsoleGroup(label string) {
	sse.Console(sdkcore.ConsoleLogModeGroup, label)
}

func (sse *ServerSentEventsHandler) ConsoleGroupEnd() {
	sse.Console(sdkcore.ConsoleLogModeGroupEnd, "")
}
