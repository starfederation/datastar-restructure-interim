package datastar

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DispatchCustomEventOptions struct {
	EventID       string
	RetryDuration time.Duration
	Selector      string
	Bubbles       bool
	Cancelable    bool
	Composed      bool
}

type DispatchCustomEventOption func(*DispatchCustomEventOptions)

func WithDispatchCustomEventEventID(id string) DispatchCustomEventOption {
	return func(o *DispatchCustomEventOptions) {
		o.EventID = id
	}
}

func WithDispatchCustomEventRetryDuration(retryDuration time.Duration) DispatchCustomEventOption {
	return func(o *DispatchCustomEventOptions) {
		o.RetryDuration = retryDuration
	}
}

func WithDispatchCustomEventSelector(selector string) DispatchCustomEventOption {
	return func(o *DispatchCustomEventOptions) {
		o.Selector = selector
	}
}

func WithDispatchCustomEventBubbles(bubbles bool) DispatchCustomEventOption {
	return func(o *DispatchCustomEventOptions) {
		o.Bubbles = bubbles
	}
}

func WithDispatchCustomEventCancelable(cancelable bool) DispatchCustomEventOption {
	return func(o *DispatchCustomEventOptions) {
		o.Cancelable = cancelable
	}
}

func WithDispatchCustomEventComposed(composed bool) DispatchCustomEventOption {
	return func(o *DispatchCustomEventOptions) {
		o.Composed = composed
	}
}

func (sse *ServerSentEventGenerator) DispatchCustomEvent(eventName string, detailsJSON string, opts ...DispatchCustomEventOption) error {
	if eventName == "" {
		return fmt.Errorf("eventName is required")
	}

	options := &DispatchCustomEventOptions{
		EventID:       "",
		RetryDuration: DefaultSseRetryDuration,
		Selector:      DefaultCustomEventSelector,
		Bubbles:       DefaultCustomEventBubbles,
		Cancelable:    DefaultCustomEventCancelable,
		Composed:      DefaultCustomEventComposed,
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

	dataLines := make([]string, 0, 5)

	dataLines = append(dataLines, EventNameDatalineLiteral+eventName)

	if options.Selector != DefaultCustomEventSelector {
		dataLines = append(dataLines, SelectorDatalineLiteral+options.Selector)
	}
	if options.Bubbles != DefaultCustomEventBubbles {
		dataLines = append(dataLines, BubblesDatalineLiteral+strconv.FormatBool(options.Bubbles))
	}
	if options.Cancelable != DefaultCustomEventCancelable {
		dataLines = append(dataLines, CancelableDatalineLiteral+strconv.FormatBool(options.Cancelable))
	}
	if options.Composed != DefaultCustomEventComposed {
		dataLines = append(dataLines, ComposedDatalineLiteral+strconv.FormatBool(options.Composed))
	}

	parts := strings.Split(detailsJSON, "\n")
	for _, part := range parts {
		dataLines = append(dataLines, DetailJsonDatalineLiteral+part)
	}

	return sse.Send(
		EventTypeDispatchCustomEvent,
		dataLines,
		sendOpts...,
	)
}
