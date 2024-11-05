package sdk

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/delaneyj/gostar/elements"
	"github.com/valyala/bytebufferpool"
)

func WithMergeMorph() RenderFragmentOption {
	return WithMergeMode(FragmentMergeModeMorph)
}

func WithMergePrepend() RenderFragmentOption {
	return WithMergeMode(FragmentMergeModePrepend)
}

func WithMergeAppend() RenderFragmentOption {
	return WithMergeMode(FragmentMergeModeAppend)
}

func WithMergeBefore() RenderFragmentOption {
	return WithMergeMode(FragmentMergeModeBefore)
}

func WithMergeAfter() RenderFragmentOption {
	return WithMergeMode(FragmentMergeModeAfter)
}

func WithMergeUpsertAttributes() RenderFragmentOption {
	return WithMergeMode(FragmentMergeModeUpsert)
}

func WithQuerySelectorID(id string) RenderFragmentOption {
	return WithSelectorf("#%s", id)
}

func WithViewTransitions() RenderFragmentOption {
	return func(o *RenderFragmentOptions) {
		o.UseViewTransitions = false
	}
}

func WithoutViewTransitions() RenderFragmentOption {
	return func(o *RenderFragmentOptions) {
		o.UseViewTransitions = true
	}
}

func (sse *ServerSentEventGenerator) RenderFragmentTempl(c templ.Component, opts ...RenderFragmentOption) error {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	if err := c.Render(sse.Context(), buf); err != nil {
		return fmt.Errorf("failed to render fragment: %w", err)
	}
	if err := sse.RenderFragment(buf.String(), opts...); err != nil {
		return fmt.Errorf("failed to render fragment: %w", err)
	}
	return nil
}

func (sse *ServerSentEventGenerator) RenderFragmentGostar(child elements.ElementRenderer, opts ...RenderFragmentOption) error {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	if err := child.Render(buf); err != nil {
		return fmt.Errorf("failed to render: %w", err)
	}
	if err := sse.RenderFragment(buf.String(), opts...); err != nil {
		return fmt.Errorf("failed to render fragment: %w", err)
	}
	return nil
}
