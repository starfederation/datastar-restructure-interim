# Architecture Decision Record: Datastar SDK

## Summary

### Issue

Datastar has had a few helper tools in the past for different languages.  The SDK effort is to unify around the tooling needed for Hypermedai On Whatever your Like (HOWL) based UIs.  Although Datastar the library can use any plugins the default bundle includes robust Server Sent Event (SSE) base approach.  Most current languages and backend don't have great tooling around the style of delivering content to the frontend.

### Decision

Provide a SDK in a language agnostic way, to that end

1. Keep ***core*** SDK as minimal as possible
2. Allow per language/framework extended features to live in a ***full*** SDK

### Status

- [x] Create a document (this) to allow any one to make a spec compliant SDK for any language or framework
- [ ] Provide a [reference implementation](../code/go/sdkcore) in Go
- [ ] Provide SDKs for
  - [ ] JS/TS
  - [ ] PHP
  - [ ] .NET
  - [ ] Python
  - [ ] Java
  - [ ] Haskell?

## Details

### Assumptions

The core mechanics of Datastar's SSE support is

1. Data gets send to browser as SSE events
2. Data comes in via JSON from browser under a `datastar` namespace

### Constraints

foo| bar |baz
--|--|--
foo | bar | bar