package datastar

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNoPathsProvided = errors.New("no paths provided")
)

type MergeStoreOptions struct {
	EventID       string
	RetryDuration time.Duration
	OnlyIfMissing bool
}

type MergeStoreOption func(*MergeStoreOptions)

func WithMergeStoreEventID(id string) MergeStoreOption {
	return func(o *MergeStoreOptions) {
		o.EventID = id
	}
}

func WithMergeStoreRetryDuration(retryDuration time.Duration) MergeStoreOption {
	return func(o *MergeStoreOptions) {
		o.RetryDuration = retryDuration
	}
}

func WithOnlyIfMissing(onlyIfMissing bool) MergeStoreOption {
	return func(o *MergeStoreOptions) {
		o.OnlyIfMissing = onlyIfMissing
	}
}

func (sse *ServerSentEventGenerator) MergeStore(storeContents []byte, opts ...MergeStoreOption) error {
	options := &MergeStoreOptions{
		EventID:       "",
		RetryDuration: DefaultSSERetryDuration,
		OnlyIfMissing: false,
	}
	for _, opt := range opts {
		opt(options)
	}

	dataRows := make([]string, 0, 32)
	if options.OnlyIfMissing {
		dataRows = append(dataRows, OnlyIfMissingDatalineLiteral+strconv.FormatBool(options.OnlyIfMissing))
	}
	lines := bytes.Split(storeContents, newLineBuf)
	for _, line := range lines {
		dataRows = append(dataRows, StoreDatalineLiteral+string(line))
	}

	sendOptions := make([]SSEEventOption, 0, 2)
	if options.EventID != "" {
		sendOptions = append(sendOptions, WithSSEEventId(options.EventID))
	}
	if options.RetryDuration != DefaultSSERetryDuration {
		sendOptions = append(sendOptions, WithSSERetryDuration(options.RetryDuration))
	}

	if err := sse.Send(
		EventTypeMergeStore,
		dataRows,
		sendOptions...,
	); err != nil {
		return fmt.Errorf("failed to send merge store: %w", err)
	}
	return nil
}

func (sse *ServerSentEventGenerator) DeleteFromStore(paths ...string) error {
	if len(paths) == 0 {
		return ErrNoPathsProvided
	}

	if err := sse.Send(
		EventTypeRemoveFromStore,
		[]string{PathsDatalineLiteral + strings.Join(paths, " ")},
	); err != nil {
		return fmt.Errorf("failed to send delete from store: %w", err)
	}
	return nil
}
