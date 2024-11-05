package sdk

import (
	"encoding/json"
	"fmt"
)

func (sse *ServerSentEventGenerator) PatchStoreMarshal(store any, opts ...PatchStoreOption) error {
	b, err := json.Marshal(store)
	if err != nil {
		panic(err)
	}
	if err := sse.PatchStore(b, opts...); err != nil {
		return fmt.Errorf("failed to patch store: %w", err)
	}

	return nil
}

func (sse *ServerSentEventGenerator) PatchStoreMarshalIfMissing(store any, opts ...PatchStoreOption) error {
	if err := sse.PatchStoreMarshal(
		store,
		append(opts, WithOnlyIfMissing(true))...,
	); err != nil {
		return fmt.Errorf("failed to patch store if missing: %w", err)
	}
	return nil
}

func (sse *ServerSentEventGenerator) PatchStoreIfMissingRaw(storeJSON string) error {
	if err := sse.PatchStore([]byte(storeJSON), WithOnlyIfMissing(true)); err != nil {
		return fmt.Errorf("failed to patch store if missing: %w", err)
	}
	return nil
}
