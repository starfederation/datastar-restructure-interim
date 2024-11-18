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

type MergeSignalsOptions struct {
	EventID       string
	RetryDuration time.Duration
	OnlyIfMissing bool
}

type MergeSignalsOption func(*MergeSignalsOptions)

func WithMergeSignalsEventID(id string) MergeSignalsOption {
	return func(o *MergeSignalsOptions) {
		o.EventID = id
	}
}

func WithMergeSignalsRetryDuration(retryDuration time.Duration) MergeSignalsOption {
	return func(o *MergeSignalsOptions) {
		o.RetryDuration = retryDuration
	}
}

func WithOnlyIfMissing(onlyIfMissing bool) MergeSignalsOption {
	return func(o *MergeSignalsOptions) {
		o.OnlyIfMissing = onlyIfMissing
	}
}

func (sse *ServerSentEventGenerator) MergeSignals(storeContents []byte, opts ...MergeSignalsOption) error {
	options := &MergeSignalsOptions{
		EventID:       "",
		RetryDuration: DefaultSseRetryDuration,
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
	if options.RetryDuration != DefaultSseRetryDuration {
		sendOptions = append(sendOptions, WithSSERetryDuration(options.RetryDuration))
	}

	if err := sse.Send(
		EventTypeMergeSignals,
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
		EventTypeRemoveSignals,
		[]string{PathsDatalineLiteral + strings.Join(paths, " ")},
	); err != nil {
		return fmt.Errorf("failed to send delete from store: %w", err)
	}
	return nil
}
