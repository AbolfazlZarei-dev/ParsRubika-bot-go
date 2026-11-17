# `account` — Listeners adapter

_Source: [`listeners_adapter.go`](../../ParsRubika/account/listeners_adapter.go)_

## Symbols

### `AddListener`

```go
func AddListener(a *Account, fn Listener, l listeners.MatcherChain) int
```

AddListener registers a callback under the supplied matcher. It is
defined here (in the account package) so callers can use either
the listeners sub-package types or their own matcher.

The function signature accepts a `listeners.MatcherChain` which both
`listeners.Messages`, `listeners.Chats`, ... and every matcher from
the `matchers` sub-package implement.

### `Kind`

```go
func (lc listenerChain) Kind() string { return lc.inner.Kind() }
```

### `Match`

```go
func (lc listenerChain) Match(ctx context.Context, raw map[string]any) bool
```
