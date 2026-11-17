# `bot` — Polling

_Source: [`polling.go`](../../ParsRubika/bot/polling.go)_

## Symbols

### `GetUpdates`

```go
func (b *Bot) GetUpdates(ctx context.Context, limit int, offsetID string) ([]*models.Update, error)
```

GetUpdates fetches a single batch of updates without touching the
stored offset or the duplicate filter.

### `Updater`

```go
func (b *Bot) Updater(ctx context.Context, limit int, pollTimeout time.Duration) ([]*models.Update, error)
```

Updater performs one long-poll round: it uses the stored offset,
advances it, filters duplicates and returns the fresh updates.

When the bot starts without a stored offset and SkipBacklog is on,
the first successful batch is discarded so old messages are not
replayed.

### `ResetOffset`

```go
func (b *Bot) ResetOffset()
```

ResetOffset forgets the stored polling offset.

### `OffsetFile`

```go
func (b *Bot) OffsetFile() string
```

OffsetFile returns the path of the on-disk offset file (empty when
persistence is disabled).
