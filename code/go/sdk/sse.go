package datastar

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ServerSentEventGenerator struct {
	ctx             context.Context
	mu              *sync.Mutex
	w               http.ResponseWriter
	flusher         http.Flusher
	shouldLogPanics bool
}

func NewSSE(w http.ResponseWriter, r *http.Request) *ServerSentEventGenerator {
	flusher, ok := w.(http.Flusher)
	if !ok {
		// This is a deliberate choice as it should never occur and is an environment issue.
		// https://crawshaw.io/blog/go-and-sqlite
		// In Go, errors that are part of the standard operation of a program are returned as values.
		// Programs are expected to handle errors.
		panic("response writer does not support flushing")
	}

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
	flusher.Flush()

	sseHandler := &ServerSentEventGenerator{
		ctx:             r.Context(),
		mu:              &sync.Mutex{},
		w:               w,
		flusher:         flusher,
		shouldLogPanics: true,
	}
	return sseHandler
}

func (sse *ServerSentEventGenerator) Context() context.Context {
	return sse.ctx
}

type ServerSentEventData struct {
	Type          EventType
	EventID       string
	Data          []string
	RetryDuration time.Duration
}

type SSEEventOption func(*ServerSentEventData)

func WithSSEEventId(id string) SSEEventOption {
	return func(e *ServerSentEventData) {
		e.EventID = id
	}
}

func WithSSERetryDuration(retryDuration time.Duration) SSEEventOption {
	return func(e *ServerSentEventData) {
		e.RetryDuration = retryDuration
	}
}

var (
	eventLinePrefix = []byte("event: ")
	idLinePrefix    = []byte("id: ")
	retryLinePrefix = []byte("retry: ")
	dataLinePrefix  = []byte("data: ")
)

func writeJustError(w io.Writer, b []byte) (err error) {
	_, err = w.Write(b)
	return err
}

func (sse *ServerSentEventGenerator) send(eventType EventType, dataLines []string, opts ...SSEEventOption) error {
	sse.mu.Lock()
	defer sse.mu.Unlock()

	// create the event
	evt := ServerSentEventData{
		Type:          eventType,
		Data:          dataLines,
		RetryDuration: DefaultSSERetryDuration,
	}

	// apply options
	for _, opt := range opts {
		opt(&evt)
	}

	// write event type
	if err := errors.Join(
		writeJustError(sse.w, eventLinePrefix),
		writeJustError(sse.w, []byte(evt.Type)),
		writeJustError(sse.w, newLineBuf),
	); err != nil {
		return fmt.Errorf("failed to write event type: %w", err)
	}

	// write id if needed
	if evt.EventID != "" {
		if err := errors.Join(
			writeJustError(sse.w, idLinePrefix),
			writeJustError(sse.w, []byte(evt.EventID)),
			writeJustError(sse.w, newLineBuf),
		); err != nil {
			return fmt.Errorf("failed to write id: %w", err)
		}
	}

	// write retry if needed
	if evt.RetryDuration.Milliseconds() > 0 {
		retry := int(evt.RetryDuration.Milliseconds())
		retryStr := strconv.Itoa(retry)
		if err := errors.Join(
			writeJustError(sse.w, retryLinePrefix),
			writeJustError(sse.w, []byte(retryStr)),
			writeJustError(sse.w, newLineBuf),
		); err != nil {
			return fmt.Errorf("failed to write retry: %w", err)
		}
	}

	// write data lines
	for _, d := range evt.Data {
		if err := errors.Join(
			writeJustError(sse.w, dataLinePrefix),
			writeJustError(sse.w, []byte(d)),
			writeJustError(sse.w, newLineBuf),
		); err != nil {
			return fmt.Errorf("failed to write data: %w", err)
		}
	}

	// write double newlines to separate events
	if err := writeJustError(sse.w, doubleNewLineBuf); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	// flush the buffer to the client
	sse.flusher.Flush()

	return nil
}
