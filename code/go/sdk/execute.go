package datastar

import (
	"strconv"
	"strings"
	"time"
)

type ExecuteJSOptions struct {
	EventID       string
	RetryDuration time.Duration
	Type          string
	AutoRemove    *bool
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

func WithExecuteJSType(t string) ExecuteJSOption {
	return func(o *ExecuteJSOptions) {
		o.Type = t
	}
}

func WithExecuteJSAutoRemove(autoremove bool) ExecuteJSOption {
	return func(o *ExecuteJSOptions) {
		o.AutoRemove = &autoremove
	}
}

func (sse *ServerSentEventGenerator) ExecuteJS(scriptContents string, opts ...ExecuteJSOption) error {
	options := &ExecuteJSOptions{
		RetryDuration: DefaultSseRetryDuration,
		Type:          DefaultExecuteJsType,
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

	dataLines := make([]string, 0, 64)

	if options.AutoRemove != nil && *options.AutoRemove != DefaultExecuteJsautoRemove {
		dataLines = append(dataLines, AutoRemoveDatalineLiteral+strconv.FormatBool(*options.AutoRemove))
	}

	if options.Type != DefaultExecuteJsType {
		dataLines = append(dataLines, TypeDatalineLiteral, options.Type)
	}

	scriptLines := strings.Split(scriptContents, NewLine)
	for _, line := range scriptLines {
		dataLines = append(dataLines, ScriptDatalineLiteral+line)
	}

	return sse.Send(
		EventTypeExecuteJs,
		dataLines,
		sendOpts...,
	)
}
