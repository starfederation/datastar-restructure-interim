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

## ServerSentEventGenerator

***There must*** be a `ServerSentEventGenerator` namespace.  In Go this is implemented as a struct, but could be a class or even namespace in languages such as C.

### Construction / Initialization
   1. ***There must*** be a way to create a new instance of this object based on the incoming `HTTP` Request and Response objects.
   2. The `ServerSentEventGenerator` ***must*** default to a flusher interface that has the following response headers set by default
      1. `Cache-Control = nocache`
      2. `Connection = keep-alive`
      3. `Content-Type = text/event-stream`
   3. Then the created response ***should*** `flush` immediately to avoid timeouts while 0-♾️ events are created
   4. `ServerSentEventGenerator` ***should*** include an incrementing number to be used as an id for events when one is not provided
   5. Multiple calls using `ServerSentEventGenerator` should be single threaded to guaruantee order.  The Go implementation use a mutex to facilitate this behavior but might not be need in a some environments

### private `send`

all top level `ServerSentEventGenerator` ***should*** use a unified sending function.
####  Args
   *  eventType [EventType](#EventType)
   *  dataLines []string
   *  options
      * `id` (string) Each event ***may*** include an `id`.  This can be used by the backend to replaye events.  If one is not provided the server ***must*** include an monotonically incrementing id
      * `retry` (duration) Each event ***may*** include a `retry` value.  If one is not provided the SDK ***must*** default to `1 second`.

#### Logic
When called the function ***must*** write to the response buffer the following in specified order
1.   ***Must*** write `event: EVENT_TYPE\n` where `EVENT_TYPE` is [EventType](#EventType)
2.   ***Must*** write `id: ID\n` where `ID` is either a user defined string or a monotonically increased integer starting at 0
3.   For each string in the provided `dataLines`, you ***must*** write `data: DATA\n` where `DATA` is the provided string.
4.  ***Must*** write a `\n\n` to complete the event per the SSE spec.
5.  Afterward the writer ***should*** immediately flush.  This can be confounded by other middlewares such as compression layers

## EventType
An enum of Datastr supported events.  Will be a string over the wire
Currently valid values are

* datastar-fragment
* datastar-signal
* datastar-delete
* datastar-redirect
* datastar-console