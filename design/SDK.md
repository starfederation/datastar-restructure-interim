# Architecture Decision Record: Datastar SDK

## Summary

### Issue

Datastar has had a few helper tools in the past for different languages.  The SDK effort is to unify around the tooling needed for Hypermedai On Whatever your Like (HOWL) based UIs.  Although Datastar the library can use any plugins the default bundle includes robust Server Sent Event (SSE) base approach.  Most current languages and backend don't have great tooling around the style of delivering content to the frontend.

### Decision

Provide a SDK in a language agnostic way, to that end

1. Keep SDK as minimal as possible
2. Allow per language/framework extended features to live in a SDK ***sugar*** version

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


# Library

> [!WARNING] All naming conventions are shown using `Go` as the standard, thing may change per language norms but please keep as close as possible.

## ServerSentEventsHandler

1. ***There must*** be a `ServerSentEventsHandler` namespace.  In Go this is implemented as a struct, but could be a class or even namespace in languages such as C.
2. ### Construction / Initialization
   1. ***There must*** be a way to create a new instance of this object based on the incoming `HTTP` Request and Response objects.
   2. The `ServerSentEventsHandler` ***must*** default to a flusher interface that has the following response headers set by default
      1. `Cache-Control = nocache`
      2. `Connection = keep-alive`
      3. `Content-Type = text/event-stream`
   3. Then the created response ***should*** `flush` immediately to avoid timeouts while 0-♾️ events are created
   4. `ServerSentEventsHandler` ***should*** include an incrementing number to be used as an id for events when one is not provided
3.