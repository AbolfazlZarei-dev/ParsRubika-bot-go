# `bot` — Copy

_Source: [`copy.go`](../../ParsRubika/bot/copy.go)_

## Symbols

### `ErrNoMessage`

```go
var ErrNoMessage = errors.New("bot: update carries no message")
```

ErrNoMessage is returned by helpers that need a message the update
does not carry.

### `ErrNothingToCopy`

```go
var ErrNothingToCopy = errors.New("bot: message has no copyable content")
```

ErrNothingToCopy is returned when a message has no copyable content.

### `CopyMessage`

```go
func (b *Bot) CopyMessage(ctx context.Context, msg *models.Message, toChatID string, opts ...SendOptions) (*models.MessageID, error)
```

CopyMessage re-sends the content of msg to toChatID without any
"forwarded from" attribution. Text, files (by file_id), stickers,
locations, contacts and polls are supported.
