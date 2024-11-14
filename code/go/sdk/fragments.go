package datastar

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type RenderFragmentOptions struct {
	EventID            string
	RetryDuration      time.Duration
	Selector           string
	MergeMode          FragmentMergeMode
	SettleDuration     time.Duration
	UseViewTransitions bool
}

type RenderFragmentOption func(*RenderFragmentOptions)

func WithSelectorf(selectorFormat string, args ...any) RenderFragmentOption {
	selector := fmt.Sprintf(selectorFormat, args...)
	return WithSelector(selector)
}

func WithSelector(selector string) RenderFragmentOption {
	return func(o *RenderFragmentOptions) {
		o.Selector = selector
	}
}

func WithMergeMode(merge FragmentMergeMode) RenderFragmentOption {
	return func(o *RenderFragmentOptions) {
		o.MergeMode = merge
	}
}

func WithSettleDuration(d time.Duration) RenderFragmentOption {
	return func(o *RenderFragmentOptions) {
		o.SettleDuration = d
	}
}

func WithUseViewTransitions(useViewTransition bool) RenderFragmentOption {
	return func(o *RenderFragmentOptions) {
		o.UseViewTransitions = useViewTransition
	}
}

type RemoveFragmentOptions struct {
	EventID            string
	RetryDuration      time.Duration
	SettleDuration     time.Duration
	UseViewTransitions *bool
}

type RemoveFragmentOption func(*RemoveFragmentOptions)

func WithRemoveEventID(id string) RemoveFragmentOption {
	return func(o *RemoveFragmentOptions) {
		o.EventID = id
	}
}

func WithRemoveRetryDuration(d time.Duration) RemoveFragmentOption {
	return func(o *RemoveFragmentOptions) {
		o.RetryDuration = d
	}
}

func WithRemoveSettleDuration(d time.Duration) RemoveFragmentOption {
	return func(o *RemoveFragmentOptions) {
		o.SettleDuration = d
	}
}

func WithRemoveUseViewTransitions(useViewTransition bool) RemoveFragmentOption {
	return func(o *RemoveFragmentOptions) {
		o.UseViewTransitions = &useViewTransition
	}
}

func (sse *ServerSentEventGenerator) RemoveFragments(selector string, opts ...RemoveFragmentOption) error {
	if selector == "" {
		panic("missing " + SelectorDatalineLiteral)
	}

	options := &RemoveFragmentOptions{
		EventID:            "",
		RetryDuration:      DefaultSSERetryDuration,
		SettleDuration:     DefaultSettleDuration,
		UseViewTransitions: nil,
	}
	for _, opt := range opts {
		opt(options)
	}

	dataRows := []string{SelectorDatalineLiteral + selector}
	if options.SettleDuration > 0 && options.SettleDuration != DefaultSettleDuration {
		settleDuration := strconv.Itoa(int(options.SettleDuration.Milliseconds()))
		dataRows = append(dataRows, SettleDurationDatalineLiteral+settleDuration)
	}
	if options.UseViewTransitions != nil {
		dataRows = append(dataRows, UseViewTransitionDatalineLiteral+strconv.FormatBool(*options.UseViewTransitions))
	}

	sendOptions := make([]SSEEventOption, 0, 2)
	if options.EventID != "" {
		sendOptions = append(sendOptions, WithSSEEventId(options.EventID))
	}
	if options.RetryDuration > 0 {
		sendOptions = append(sendOptions, WithSSERetryDuration(options.RetryDuration))
	}

	if err := sse.send(EventTypeRemove, dataRows, sendOptions...); err != nil {
		return fmt.Errorf("failed to send remove: %w", err)
	}
	return nil
}

func (sse *ServerSentEventGenerator) RenderFragment(fragment string, opts ...RenderFragmentOption) error {
	options := &RenderFragmentOptions{
		EventID:        "",
		RetryDuration:  DefaultSSERetryDuration,
		Selector:       "",
		MergeMode:      FragmentMergeModeMorph,
		SettleDuration: DefaultSettleDuration,
	}
	for _, opt := range opts {
		opt(options)
	}

	sendOptions := make([]SSEEventOption, 0, 2)
	if options.EventID != "" {
		sendOptions = append(sendOptions, WithSSEEventId(options.EventID))
	}
	if options.RetryDuration > 0 {
		sendOptions = append(sendOptions, WithSSERetryDuration(options.RetryDuration))
	}

	dataRows := make([]string, 0, 4)
	if options.Selector != "" {
		dataRows = append(dataRows, SelectorDatalineLiteral+options.Selector)
	}
	if options.MergeMode != FragmentMergeModeMorph {
		dataRows = append(dataRows, MergeDatalineLiteral+string(options.MergeMode))
	}
	if options.SettleDuration > 0 && options.SettleDuration != DefaultSettleDuration {
		settleDuration := strconv.Itoa(int(options.SettleDuration.Milliseconds()))
		dataRows = append(dataRows, SettleDurationDatalineLiteral+settleDuration)
	}
	if options.UseViewTransitions {
		dataRows = append(dataRows, UseViewTransitionDatalineLiteral+"true")
	}

	if fragment != "" {
		parts := strings.Split(fragment, "\n")
		for _, part := range parts {
			dataRows = append(dataRows, FragmentDatalineLiteral+part)
		}
	}

	if err := sse.send(
		EventTypeFragment,
		dataRows,
		sendOptions...,
	); err != nil {
		return fmt.Errorf("failed to send fragment: %w", err)
	}

	return nil
}
