package sdk

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DEFAULT_SETTLE_TIME = 300 * time.Millisecond

type FragmentMergeType string

const (
	FragmentMergeMorph   FragmentMergeType = "morph"
	FragmentMergeInner   FragmentMergeType = "inner"
	FragmentMergeOuter   FragmentMergeType = "outer"
	FragmentMergePrepend FragmentMergeType = "prepend"
	FragmentMergeAppend  FragmentMergeType = "append"
	FragmentMergeBefore  FragmentMergeType = "before"
	FragmentMergeAfter   FragmentMergeType = "after"
	FragmentMergeUpsert  FragmentMergeType = "upsert_attributes"
)

type RenderFragmentOptions struct {
	QuerySelector      string
	Merge              FragmentMergeType
	SettleDuration     time.Duration
	UseViewTransitions bool
}

type RenderFragmentOption func(*RenderFragmentOptions)

func WithQuerySelectorf(selectorFormat string, args ...any) RenderFragmentOption {
	selector := fmt.Sprintf(selectorFormat, args...)
	return WithQuerySelector(selector)
}

func WithQuerySelector(selector string) RenderFragmentOption {
	return func(o *RenderFragmentOptions) {
		o.QuerySelector = selector
	}
}

func WithMergeType(merge FragmentMergeType) RenderFragmentOption {
	return func(o *RenderFragmentOptions) {
		o.Merge = merge
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

func (sse *ServerSentEventGenerator) DeleteSelectorf(selectorFormat string, args ...any) error {
	selector := fmt.Sprintf(selectorFormat, args...)
	return sse.DeleteSelector(selector)
}

func (sse *ServerSentEventGenerator) DeleteSelector(selector string, opts ...RenderFragmentOption) error {
	if selector == "" {
		panic("missing selector")
	}

	dataRows := []string{"selector " + selector}
	if err := sse.send(EventTypeDelete, dataRows); err != nil {
		return fmt.Errorf("failed to send delete: %w", err)
	}
	return nil
}

func (sse *ServerSentEventGenerator) RenderFragment(fragment string, opts ...RenderFragmentOption) error {
	options := &RenderFragmentOptions{
		QuerySelector:  "",
		Merge:          FragmentMergeMorph,
		SettleDuration: DEFAULT_SETTLE_TIME,
	}
	for _, opt := range opts {
		opt(options)
	}

	dataRows := make([]string, 0, 4)
	if options.QuerySelector != "" {
		dataRows = append(dataRows, "selector "+options.QuerySelector)
	}
	if options.Merge != FragmentMergeMorph {
		dataRows = append(dataRows, "merge "+string(options.Merge))
	}
	if options.SettleDuration > 0 && options.SettleDuration != DEFAULT_SETTLE_TIME {
		settleTime := strconv.Itoa(int(options.SettleDuration.Milliseconds()))
		dataRows = append(dataRows, "settle "+settleTime)
	}
	if options.UseViewTransitions {
		dataRows = append(dataRows, "useViewTransition true")
	}

	if fragment != "" {
		parts := strings.Split(fragment, "\n")
		parts[0] = "fragment " + parts[0]
		dataRows = append(dataRows, parts...)
	}

	if err := sse.send(
		EventTypeFragment,
		dataRows,
		WithSSERetry(0),
	); err != nil {
		return fmt.Errorf("failed to send fragment: %w", err)
	}

	return nil
}
