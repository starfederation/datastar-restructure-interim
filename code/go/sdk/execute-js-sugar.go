package datastar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func (sse *ServerSentEventGenerator) ConsoleLog(msg string, opts ...ExecuteJsOption) error {
	call := fmt.Sprintf("console.log(%q)", msg)
	return sse.ExecuteJs(call, opts...)
}

func (sse *ServerSentEventGenerator) ConsoleLogf(format string, args ...any) error {
	return sse.ConsoleLog(fmt.Sprintf(format, args...))
}

func (sse *ServerSentEventGenerator) ConsoleError(err error, opts ...ExecuteJsOption) error {
	call := fmt.Sprintf("console.error(%q)", err.Error())
	return sse.ExecuteJs(call, opts...)
}

func (sse *ServerSentEventGenerator) Redirectf(format string, args ...any) error {
	url := fmt.Sprintf(format, args...)
	return sse.Redirect(url)
}

func (sse *ServerSentEventGenerator) Redirect(url string, opts ...ExecuteJsOption) error {
	js := fmt.Sprintf("window.location.href = %q;", url)
	return sse.ExecuteJs(js, opts...)
}

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

func (sse *ServerSentEventGenerator) DispatchCustomEvent(eventName string, detail any, opts ...DispatchCustomEventOption) error {
	if eventName == "" {
		return fmt.Errorf("eventName is required")
	}

	detailsJSON, err := json.Marshal(detail)
	if err != nil {
		return fmt.Errorf("failed to marshal detail: %w", err)
	}

	options := DispatchCustomEventOptions{
		EventID:       "",
		RetryDuration: DefaultSseRetryDuration,
		Selector:      DefaultCustomEventSelector,
		Bubbles:       DefaultCustomEventBubbles,
		Cancelable:    DefaultCustomEventCancelable,
		Composed:      DefaultCustomEventComposed,
	}

	for _, opt := range opts {
		opt(&options)
	}

	elementsJS := `[document]`
	if options.Selector != "" && options.Selector != DefaultCustomEventSelector {
		elementsJS = fmt.Sprintf(`document.querySelectorAll(%q)`, options.Selector)
	}

	js := fmt.Sprintf(`
const elements = %s

const event = new CustomEvent(%q, {
	bubbles: %t,
	cancelable: %t,
	composed: %t,
	detail: %s,
});

elements.forEach((element) => {
	element.dispatchEvent(event);
});
	`,
		elementsJS,
		eventName,
		options.Bubbles,
		options.Cancelable,
		options.Composed,
		string(detailsJSON),
	)

	executeOptions := make([]ExecuteJsOption, 0)
	if options.EventID != "" {
		executeOptions = append(executeOptions, WithExecuteJsEventID(options.EventID))
	}
	if options.RetryDuration != 0 {
		executeOptions = append(executeOptions, WithExecuteJsRetryDuration(options.RetryDuration))
	}

	return sse.ExecuteJs(js, executeOptions...)

}

func (sse *ServerSentEventGenerator) ReplaceURL(u url.URL, opts ...ExecuteJsOption) error {
	js := fmt.Sprintf(`window.history.replaceState({}, "", %q)`, u.String())
	return sse.ExecuteJs(js, opts...)
}

func (sse *ServerSentEventGenerator) ReplaceURLQuerystring(r *http.Request, values url.Values, opts ...ExecuteJsOption) error {
	u := *r.URL
	u.RawQuery = values.Encode()
	return sse.ReplaceURL(u, opts...)
}
