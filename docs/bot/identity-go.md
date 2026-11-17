# `bot` — Identity

_Source: [`identity.go`](../../ParsRubika/bot/identity.go)_

## Symbols

### `Me`

```go
func (b *Bot) Me(ctx context.Context) (*models.Bot, error)
```

Me returns the bot's own profile. The answer is cached for ten
minutes, so it is cheap to call from handlers.

### `Name`

```go
func (b *Bot) Name(ctx context.Context) (string, error)
```

Name returns the display name of the bot.

### `Username`

```go
func (b *Bot) Username(ctx context.Context) (string, error)
```

Username returns the bot's @username (without the "@").
