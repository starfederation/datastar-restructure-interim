package sdk

import (
	"fmt"
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

func WithQuerySelector(selectorFormat string, args ...any) RenderFragmentOption {
	selector := fmt.Sprintf(selectorFormat, args...)
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

func (sse *ServerSentEventsHandler) DeleteSelector(selector string, opts ...RenderFragmentOption) {

	if selector == "" {
		panic("missing selector")
	}

	dataRows := []string{fmt.Sprintf("selector %s", selector)}

	sse.send(EventTypeDelete, dataRows)
}

func (sse *ServerSentEventsHandler) RenderFragment(fragment string, opts ...RenderFragmentOption) error {
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
		dataRows = append(dataRows, fmt.Sprintf("selector %s", options.QuerySelector))
	}
	if options.Merge != FragmentMergeMorph {
		dataRows = append(dataRows, fmt.Sprintf("merge %s", options.Merge))
	}
	if options.SettleDuration > 0 && options.SettleDuration != DEFAULT_SETTLE_TIME {
		dataRows = append(dataRows, fmt.Sprintf("settle %d", options.SettleDuration.Milliseconds()))
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
