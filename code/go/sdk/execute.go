package datastar

import (
	"strconv"
	"strings"
	"time"
)

type ExecuteJSOptions struct {
	EventID          string
	RetryDuration    time.Duration
	AutoRemoveScript *bool
}

type ExecuteJSOption func(*ExecuteJSOptions)

func WithExecuteJSEventID(id string) ExecuteJSOption {
	return func(o *ExecuteJSOptions) {
		o.EventID = id
	}
}

func WithExecuteJSRetryDuration(retryDuration time.Duration) ExecuteJSOption {
	return func(o *ExecuteJSOptions) {
		o.RetryDuration = retryDuration
	}
}

func WithExecuteJSAutoRemoveScript(autoremove bool) ExecuteJSOption {
	return func(o *ExecuteJSOptions) {
		o.AutoRemoveScript = &autoremove
	}
}

func (sse *ServerSentEventGenerator) ExecuteJS(scriptContents string, opts ...ExecuteJSOption) error {
	options := &ExecuteJSOptions{
		EventID:          "",
		RetryDuration:    DefaultSseRetryDuration,
		AutoRemoveScript: nil,
	}
	for _, opt := range opts {
		opt(options)
	}

	sendOpts := make([]SSEEventOption, 0, 2)
	if options.EventID != "" {
		sendOpts = append(sendOpts, WithSSEEventId(options.EventID))
	}

	if options.RetryDuration != DefaultSseRetryDuration {
		sendOpts = append(sendOpts, WithSSERetryDuration(options.RetryDuration))
	}

	dataLines := make([]string, 0, 2)
	if options.AutoRemoveScript != nil && *options.AutoRemoveScript != DefaultAutoRemoveScript {
		dataLines = append(dataLines, AutoRemoveScriptDatalineLiteral+strconv.FormatBool(*options.AutoRemoveScript))
	}

	scriptLines := strings.Split(scriptContents, "\n")
	for _, line := range scriptLines {
		dataLines = append(dataLines, ScriptDatalineLiteral+line)
	}

	return sse.Send(
		EventTypeExecuteJs,
		dataLines,
		sendOpts...,
	)
}
