package sdk

import (
	"encoding/json"
	"fmt"

	sdkcore "github.com/starfederation/datastar/go/sdkcore"
)

func (sse *ServerSentEventsHandler) PatchStoreMarshal(store any, opts ...sdkcore.PatchStoreOption) error {
	b, err := json.Marshal(store)
	if err != nil {
		panic(err)
	}
	if err := sse.PatchStore(b, opts...); err != nil {
		return fmt.Errorf("failed to patch store: %w", err)
	}

	return nil
}

func (sse *ServerSentEventsHandler) PatchStoreMarshalIfMissing(store any, opts ...sdkcore.PatchStoreOption) error {
	if err := sse.PatchStoreMarshal(
		store,
		append(opts, sdkcore.WithOnlyIfMissing(true))...,
	); err != nil {
		return fmt.Errorf("failed to patch store if missing: %w", err)
	}
	return nil
}

func (sse *ServerSentEventsHandler) PatchStoreIfMissingRaw(storeJSON string) {
	if err := sse.PatchStore([]byte(storeJSON), sdkcore.WithOnlyIfMissing(true)); err != nil {
		panic(err)
	}
}
