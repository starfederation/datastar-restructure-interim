package datastar

import "fmt"

func (sse *ServerSentEventGenerator) Redirectf(format string, args ...any) error {
	url := fmt.Sprintf(format, args...)
	return sse.Redirect(url)
}

func (sse *ServerSentEventGenerator) Redirect(url string, opts ...SSEEventOption) error {
	if err := sse.send(
		EventTypeRedirect,
		[]string{"url " + url},
		opts...,
	); err != nil {
		return fmt.Errorf("failed to send redirect: %w", err)
	}
	return nil
}
