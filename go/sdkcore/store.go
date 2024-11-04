package sdk

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNoPathsProvided = errors.New("no paths provided")
)

func (sse *ServerSentEventsHandler) DeleteFromStore(paths ...string) error {
	if len(paths) == 0 {
		return ErrNoPathsProvided
	}

	if err := sse.send(
		EventTypeDelete,
		[]string{
			fmt.Sprintf("paths %s", strings.Join(paths, " ")),
		},
	); err != nil {
		return fmt.Errorf("failed to send delete from store: %w", err)
	}
	return nil
}

type PatchStoreOptions struct {
	OnlyIfMissing bool
}

type PatchStoreOption func(*PatchStoreOptions)

func WithOnlyIfMissing(onlyIfMissing bool) PatchStoreOption {
	return func(o *PatchStoreOptions) {
		o.OnlyIfMissing = onlyIfMissing
	}
}

func (sse *ServerSentEventsHandler) PatchStore(storeContents []byte, opts ...PatchStoreOption) error {
	options := &PatchStoreOptions{
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
		dataRows = append(dataRows, fmt.Sprintf("store %s", line))
	}

	if err := sse.send(
		EventTypeSignal,
		dataRows,
	); err != nil {
		return fmt.Errorf("failed to send patch store: %w", err)
	}
	return nil
}
