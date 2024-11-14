package datastar

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/delaneyj/gostar/elements"
	"github.com/valyala/bytebufferpool"
)

var ValidFragmentMergeTypes = []FragmentMergeMode{
	FragmentMergeModeMorph,
	FragmentMergeModeInner,
	FragmentMergeModeOuter,
	FragmentMergeModePrepend,
	FragmentMergeModeAppend,
	FragmentMergeModeBefore,
	FragmentMergeModeAfter,
	FragmentMergeModeUpsertAttributes,
}

func FragmentMergeTypeFromString(s string) (FragmentMergeMode, error) {
	for _, t := range ValidFragmentMergeTypes {
		if string(t) == s {
			return t, nil
		}
	}
	return "", fmt.Errorf("invalid fragment merge type: %s", s)
}

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
	return WithMergeMode(FragmentMergeModeUpsertAttributes)
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

func (sse *ServerSentEventGenerator) RenderFragmentf(format string, args ...any) error {
	return sse.RenderFragment(fmt.Sprintf(format, args...))
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

func GetSSE(urlFormat string, args ...any) string {
	return fmt.Sprintf(`$get('%s')`, fmt.Sprintf(urlFormat, args...))
}
func PostSSE(urlFormat string, args ...any) string {
	return fmt.Sprintf(`$post('%s')`, fmt.Sprintf(urlFormat, args...))
}
func PutSSE(urlFormat string, args ...any) string {
	return fmt.Sprintf(`$put('%s')`, fmt.Sprintf(urlFormat, args...))
}
func PatchSSE(urlFormat string, args ...any) string {
	return fmt.Sprintf(`$patch('%s')`, fmt.Sprintf(urlFormat, args...))
}
func DeleteSSE(urlFormat string, args ...any) string {
	return fmt.Sprintf(`$delete('%s')`, fmt.Sprintf(urlFormat, args...))
}
