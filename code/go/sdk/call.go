package datastar

func (sse *ServerSentEventGenerator) Call(js string, opts ...SSEEventOption) error {
	return sse.Send(
		EventTypeCall,
		[]string{js},
		opts...,
	)
}
