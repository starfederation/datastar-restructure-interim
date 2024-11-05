package sdk

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/delaneyj/gostar/elements"
	sdkcore "github.com/starfederation/datastar/go/sdkcore"
	"github.com/valyala/bytebufferpool"
)

func WithMergeMorph() sdkcore.RenderFragmentOption {
	return sdkcore.WithMergeType(sdkcore.FragmentMergeMorph)
}

func WithMergePrepend() sdkcore.RenderFragmentOption {
	return sdkcore.WithMergeType(sdkcore.FragmentMergePrepend)
}

func WithMergeAppend() sdkcore.RenderFragmentOption {
	return sdkcore.WithMergeType(sdkcore.FragmentMergeAppend)
}

func WithMergeBefore() sdkcore.RenderFragmentOption {
	return sdkcore.WithMergeType(sdkcore.FragmentMergeBefore)
}

func WithMergeAfter() sdkcore.RenderFragmentOption {
	return sdkcore.WithMergeType(sdkcore.FragmentMergeAfter)
}

func WithMergeUpsertAttributes() sdkcore.RenderFragmentOption {
	return sdkcore.WithMergeType(sdkcore.FragmentMergeUpsert)
}

func WithQuerySelectorID(id string) sdkcore.RenderFragmentOption {
	return sdkcore.WithQuerySelector("#%s", id)
}

func WithViewTransitions() sdkcore.RenderFragmentOption {
	return func(o *sdkcore.RenderFragmentOptions) {
		o.UseViewTransitions = false
	}
}

func WithoutViewTransitions() sdkcore.RenderFragmentOption {
	return func(o *sdkcore.RenderFragmentOptions) {
		o.UseViewTransitions = true
	}
}

func (sse *ServerSentEventsHandler) RenderFragmentTempl(c templ.Component, opts ...sdkcore.RenderFragmentOption) error {
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

func RenderFragmentGostar(sse *ServerSentEventsHandler, child elements.ElementRenderer, opts ...sdkcore.RenderFragmentOption) error {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	if err := child.Render(buf); err != nil {
		return fmt.Errorf("failed to render: %w", err)
	}
	return sse.RenderFragment(buf.String(), opts...)
}
