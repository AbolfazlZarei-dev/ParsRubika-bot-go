# `account` — Methods ownership

_Source: [`methods_ownership.go`](../../ParsRubika/account/methods_ownership.go)_

## Symbols

### `RequestChangeOwner`

```go
func (a *Account) RequestChangeOwner(ctx context.Context, objectGUID, newOwnerGUID string) (*events.Event, error)
```

RequestChangeOwner requests ownership transfer to a new user.

### `AcceptOwnershipRequest`

```go
func (a *Account) AcceptOwnershipRequest(ctx context.Context, objectGUID string) (*events.Event, error)
```

AcceptOwnershipRequest accepts a pending ownership transfer.

### `RejectOwnershipRequest`

```go
func (a *Account) RejectOwnershipRequest(ctx context.Context, objectGUID string) (*events.Event, error)
```

RejectOwnershipRequest rejects a pending ownership transfer.

### `CancelChangeOwner`

```go
func (a *Account) CancelChangeOwner(ctx context.Context, objectGUID string) (*events.Event, error)
```

CancelChangeOwner cancels a pending ownership transfer.

### `GetPendingOwner`

```go
func (a *Account) GetPendingOwner(ctx context.Context, objectGUID string) (*events.Event, error)
```

GetPendingOwner returns the pending owner of a chat.
