# `account` — Dispatch

_Source: [`dispatch.go`](../../ParsRubika/account/dispatch.go)_

## Symbols

### `AddListener`

```go
func (a *Account) AddListener(fn Listener, m Matcher) int
```

AddListener registers a callback under the supplied matcher. It
returns the integer ID of the listener; pass it to RemoveListener to
deregister.

### `RemoveListener`

```go
func (a *Account) RemoveListener(id int) bool
```

RemoveListener removes a previously-registered listener by its ID.
