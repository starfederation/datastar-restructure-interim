# Release Notes for Datastar

## 0.20.0 - Unreleased

> [!WARNING]
> This update contains breaking changes to attributes and SSE events.

### Added

- Added SDKs for .NET, Go, PHP and TypeScript.
- Added the `data-persist` attribute.
- Added the `data-replace-url` attribute.

### Changed

- Changed the `$$` prefix to `🚀` (or no prefix) for all plugin actions.
- Changed the `data-header` syntax to accept an object of key-value pairs.
- Renamed the `datastar-delete` event to `datastar-remove`.
- Renamed the `upsert_attributes` merge type to `upsertAttributes` in the fragment event.
- Renamed the `settle` option to `settleDuration` in the fragment event.
- Renamed the `vt` option to `useViewTransition` in the fragment event.

### Removed

- Removed the `local` and `session` modifiers from `data-store`. Use the `data-persist` attribute instead.
