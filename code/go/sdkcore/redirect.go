package sdk

import "fmt"

func (sse *ServerSentEventsHandler) Redirect(urlFormat string, args ...any) {
	url := fmt.Sprintf(urlFormat, args...)
	sse.send(
		EventTypeRedirect,
		[]string{
			fmt.Sprintf("url: %s", url),
		},
		WithSSERetry(0),
	)
}
