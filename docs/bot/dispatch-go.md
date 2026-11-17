# `bot` — Dispatch

_Source: [`dispatch.go`](../../ParsRubika/bot/dispatch.go)_

## Symbols

### `Dispatch`

```go
func (b *Bot) Dispatch(ctx context.Context, u *models.Update) error
```

Dispatch feeds an update to the bot synchronously, exactly as the
polling loop would. It is handy for tests and custom transports.
