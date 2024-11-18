package datastar

import (
	"strconv"
	"strings"
	"time"
)

type ExecuteJsOptions struct {
	EventID       string
	RetryDuration time.Duration
	Type          string
	AutoRemove    *bool
}

type ExecuteJsOption func(*ExecuteJsOptions)

func WithExecuteJsEventID(id string) ExecuteJsOption {
	return func(o *ExecuteJsOptions) {
		o.EventID = id
	}
}

func WithExecuteJsRetryDuration(retryDuration time.Duration) ExecuteJsOption {
	return func(o *ExecuteJsOptions) {
		o.RetryDuration = retryDuration
	}
}

func WithExecuteJsType(t string) ExecuteJsOption {
	return func(o *ExecuteJsOptions) {
		o.Type = t
	}
}

func WithExecuteJsAutoRemove(autoremove bool) ExecuteJsOption {
	return func(o *ExecuteJsOptions) {
		o.AutoRemove = &autoremove
	}
}

func (sse *ServerSentEventGenerator) ExecuteJs(scriptContents string, opts ...ExecuteJsOption) error {
	options := &ExecuteJsOptions{
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

	if options.AutoRemove != nil && *options.AutoRemove != DefaultExecuteJsAutoRemove {
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
