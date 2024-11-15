package datastar

import (
	"encoding/json"
	"fmt"
)

func (sse *ServerSentEventGenerator) MarshalAndMergeStore(store any, opts ...MergeStoreOption) error {
	b, err := json.Marshal(store)
	if err != nil {
		panic(err)
	}
	if err := sse.MergeStore(b, opts...); err != nil {
		return fmt.Errorf("failed to merge store: %w", err)
	}

	return nil
}

func (sse *ServerSentEventGenerator) MarshalAndMergeStoreIfMissing(store any, opts ...MergeStoreOption) error {
	if err := sse.MarshalAndMergeStore(
		store,
		append(opts, WithOnlyIfMissing(true))...,
	); err != nil {
		return fmt.Errorf("failed to merge store if missing: %w", err)
	}
	return nil
}

func (sse *ServerSentEventGenerator) MergeStoreIfMissingRaw(storeJSON string) error {
	if err := sse.MergeStore([]byte(storeJSON), WithOnlyIfMissing(true)); err != nil {
		return fmt.Errorf("failed to merge store if missing: %w", err)
	}
	return nil
}
