# `account` — Methods gifs

_Source: [`methods_gifs.go`](../../ParsRubika/account/methods_gifs.go)_

## Symbols

### `GetMyGifSet`

```go
func (a *Account) GetMyGifSet(ctx context.Context) (*events.Event, error)
```

GetMyGifSet returns the user's saved gif set.

### `AddToMyGifSet`

```go
func (a *Account) AddToMyGifSet(ctx context.Context, objectGUID, messageID string) (*events.Event, error)
```

AddToMyGifSet adds a gif to the user's gif set.

### `RemoveFromMyGifSet`

```go
func (a *Account) RemoveFromMyGifSet(ctx context.Context, fileID string) (*events.Event, error)
```

RemoveFromMyGifSet removes a gif from the user's gif set. The API
uses `file_id`, not `object_guid + message_id`.

### `DeleteAllGifs`

```go
func (a *Account) DeleteAllGifs(ctx context.Context) error
```

DeleteAllGifs removes all gifs from the user's gif set. The file_id
may come back as either a string or a JSON number depending on the
server, so we use the event's AnyStr accessor (via a temporary Event
wrapper) to coerce it.
