package sdk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const DEFAULT_SSE_SEND_RETRY = 1 * time.Second

type ServerSentEventsHandler struct {
	ctx             context.Context
	mu              *sync.Mutex
	w               http.ResponseWriter
	flusher         http.Flusher
	shouldLogPanics bool
	nextID          int
}

func NewSSE(w http.ResponseWriter, r *http.Request) (*ServerSentEventsHandler, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, ErrFlushingNotSupported
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
	flusher.Flush()

	sseHandler := &ServerSentEventsHandler{
		ctx:             r.Context(),
		mu:              &sync.Mutex{},
		w:               w,
		flusher:         flusher,
		shouldLogPanics: true,
		nextID:          0,
	}
	return sseHandler, nil
}

func (sse *ServerSentEventsHandler) Context() context.Context {
	return sse.ctx
}

type SSEEvent struct {
	Type            EventType
	ShouldIncludeId bool
	ID              string
	Data            []string
	Retry           time.Duration
}

type SSEEventOption func(*SSEEvent)

func WithSSEId(id string) SSEEventOption {
	return func(e *SSEEvent) {
		e.ID = id
		e.ShouldIncludeId = true
	}
}

func WithSSERetry(retry time.Duration) SSEEventOption {
	return func(e *SSEEvent) {
		e.Retry = retry
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

func (sse *ServerSentEventsHandler) send(eventType EventType, dataLines []string, opts ...SSEEventOption) error {
	sse.mu.Lock()
	defer sse.mu.Unlock()

	evt := SSEEvent{
		Type:  eventType,
		Data:  dataLines,
		Retry: DEFAULT_SSE_SEND_RETRY,
	}

	// apply options
	for _, opt := range opts {
		opt(&evt)
	}

	// write event type
	if err := errors.Join(
		writeJustError(sse.w, eventLinePrefix),
		writeJustError(sse.w, evt.Type),
		writeJustError(sse.w, newLineBuf),
	); err != nil {
		return fmt.Errorf("failed to write event type: %w", err)
	}

	// write id if needed
	if evt.ShouldIncludeId {
		if evt.ID == "" {
			evt.ID = fmt.Sprintf("%d", sse.nextID)
			sse.nextID++
		}

		if err := errors.Join(
			writeJustError(sse.w, idLinePrefix),
			writeJustError(sse.w, []byte(evt.ID)),
			writeJustError(sse.w, newLineBuf),
		); err != nil {
			return fmt.Errorf("failed to write id: %w", err)
		}
	}

	// write retry if needed
	if evt.Retry.Milliseconds() > 0 {
		if err := errors.Join(
			writeJustError(sse.w, retryLinePrefix),
			writeJustError(sse.w, []byte(fmt.Sprintf("%d", evt.Retry.Milliseconds()))),
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

	sse.flusher.Flush()

	return nil
}
