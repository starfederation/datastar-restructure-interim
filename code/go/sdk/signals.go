package datastar

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrNoPathsProvided = errors.New("no paths provided")
)

type PatchStoreOptions struct {
	EventID       string
	RetryDuration time.Duration
	OnlyIfMissing bool
}

type PatchStoreOption func(*PatchStoreOptions)

func WithPatchStoreEventID(id string) PatchStoreOption {
	return func(o *PatchStoreOptions) {
		o.EventID = id
	}
}

func WithPatchStoreRetryDuration(retryDuration time.Duration) PatchStoreOption {
	return func(o *PatchStoreOptions) {
		o.RetryDuration = retryDuration
	}
}

func WithOnlyIfMissing(onlyIfMissing bool) PatchStoreOption {
	return func(o *PatchStoreOptions) {
		o.OnlyIfMissing = onlyIfMissing
	}
}

func (sse *ServerSentEventGenerator) PatchStore(storeContents []byte, opts ...PatchStoreOption) error {
	options := &PatchStoreOptions{
		EventID:       "",
		RetryDuration: DefaultSSERetryDuration,
		OnlyIfMissing: false,
	}
	for _, opt := range opts {
		opt(options)
	}

	dataRows := make([]string, 0, 32)
	if options.OnlyIfMissing {
		dataRows = append(dataRows, fmt.Sprintf("onlyIfMissing %t", options.OnlyIfMissing))
	}
	lines := bytes.Split(storeContents, newLineBuf)
	for _, line := range lines {
		dataRows = append(dataRows, "store "+string(line))
	}

	sendOptions := make([]SSEEventOption, 0, 2)
	if options.EventID != "" {
		sendOptions = append(sendOptions, WithSSEEventId(options.EventID))
	}
	if options.RetryDuration != DefaultSSERetryDuration {
		sendOptions = append(sendOptions, WithSSERetryDuration(options.RetryDuration))
	}

	if err := sse.send(
		EventTypeSignal,
		dataRows,
		sendOptions...,
	); err != nil {
		return fmt.Errorf("failed to send patch store: %w", err)
	}
	return nil
}

func (sse *ServerSentEventGenerator) DeleteFromStore(paths ...string) error {
	if len(paths) == 0 {
		return ErrNoPathsProvided
	}

	if err := sse.send(
		EventTypeRemove,
		[]string{"paths " + strings.Join(paths, " ")},
	); err != nil {
		return fmt.Errorf("failed to send delete from store: %w", err)
	}
	return nil
}
