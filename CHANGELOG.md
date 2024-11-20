# Release Notes for Datastar

## 0.20.0 - Unreleased

> [!WARNING]
> This update contains breaking changes to attributes, actions and SSE events.

### Added

- Added SDKs for Go, PHP, .NET  and TypeScript.
- Added the `data-persist` attribute.
- Added the `data-replace-url` attribute.
- Added the `datastar-execute-script` SSE event.
- Added `replaceUrl` examples from both signals and the backend.
- Added `data-sse-indicator` attribute, which upsert a signal to show when SSE is active.

### Changed

- Changed the `$$` prefix to `$` for action plugins.
- Changed the `data-header` syntax to accept an object of key-value pairs.
- Renamed the `datastar-fragment` SSE event to `datastar-merge-fragments`.
- Renamed the `datastar-signal` SSE event to `datastar-merge-signals`.
- Renamed the `fragment` dataline literal for SSE events to `fragments`.
- Renamed the `store` dataline literal for SSE events to `signals`.
- Renamed the `upsert_attributes` merge mode to `upsertAttributes` in the fragment event.
- Renamed the `settle` option to `settleDuration` in the fragment event and changed the default value to `300`.
- Renamed the `vt` option to `useViewTransition` in the fragment event and changed the default value to `false`.
- Turned second argument of sse actions from `onlyRemoteSignals` to an optional object with `header` and `onlyRemoteSignals` keys currently defaulting to `{}` and `true` respectively.
- Now use error codes that make roughly with HTTP status codes.  Reduced binary size significantly.
- `data-model` will attempt to upsert missing signals

### Removed

- Removed the `~ref` syntax. Use the new `$ref()` action instead.
- Removed the `local` and `session` modifiers from `data-store`. Use the new `data-persist` attribute instead.
- Removed the `datastar-delete` SSE event. Use the new `datastar-remove-fragments` and `datastar-remove-signals` SSE events instead.
- Removed `sendDatastarEvent` from ctx.  We have to rethink how to expose events for a better try at the inspector.
- Removed `data-show` now covered by `data-class`
- Removed `$$isFetching` action and `data-fetch-indicator` attribute.  Use data-sse-indicator instead.
- Removed `data-header` attribute.  Use `header` args in sse actions instead.
- The concept of `_dsPlugins` is removed, more consisten architecture made is unnecessary.
- Removed the `datastar-redirect` and `datastar-console` SSE events.  Use the new `datastar-execute-script` SSE event instead.
