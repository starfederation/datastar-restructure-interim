package datastar

import (
	"fmt"

	"github.com/goccy/go-json"
)

func (sse *ServerSentEventGenerator) DispatchCustomEventMarshal(eventName string, detail any, opts ...DispatchCustomEventOption) error {
	b, err := json.Marshal(detail)
	if err != nil {
		return fmt.Errorf("failed to marshal detail: %w", err)
	}

	return sse.DispatchCustomEvent(eventName, string(b), opts...)
}
