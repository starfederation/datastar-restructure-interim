package sdk

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DEFAULT_SETTLE_TIME = 300 * time.Millisecond

type FragmentMergeMode string

const (
	FragmentMergeModeMorph   FragmentMergeMode = "morph"
	FragmentMergeModeInner   FragmentMergeMode = "inner"
	FragmentMergeModeOuter   FragmentMergeMode = "outer"
	FragmentMergeModePrepend FragmentMergeMode = "prepend"
	FragmentMergeModeAppend  FragmentMergeMode = "append"
	FragmentMergeModeBefore  FragmentMergeMode = "before"
	FragmentMergeModeAfter   FragmentMergeMode = "after"
	FragmentMergeModeUpsert  FragmentMergeMode = "upsertAttributes"
)

type RenderFragmentOptions struct {
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

type DeleteFragmentOptions struct {
	SettleDuration     time.Duration
	UseViewTransitions bool
}

func (sse *ServerSentEventGenerator) DeleteFragments(selector string, opts ...DeleteFragmentOptions) error {
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
		Selector:       "",
		MergeMode:      FragmentMergeModeMorph,
		SettleDuration: DEFAULT_SETTLE_TIME,
	}
	for _, opt := range opts {
		opt(options)
	}

	dataRows := make([]string, 0, 4)
	if options.Selector != "" {
		dataRows = append(dataRows, "selector "+options.Selector)
	}
	if options.MergeMode != FragmentMergeModeMorph {
		dataRows = append(dataRows, "merge "+string(options.MergeMode))
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
