# `scheduler` — Package scheduler wraps the bot's scheduler to provide a simple

_Source: [`send.go`](../../../ParsRubika/extras/scheduler/send.go)_

## Symbols

### `Sender`

```go
type Sender interface
```

Sender is the minimal interface a scheduler client needs (both
`*bot.Bot` and `*account.Account` satisfy it).

### `SendLater`

```go
func SendLater(parent context.Context, sender Sender, chatID, text string, delay time.Duration) (cancel func())
```

SendLater schedules `text` to be sent to `chatID` after `delay`.
The supplied `sender` is used to send the message; pass the bot
client or the account client.

The returned `cancel` function can be invoked to cancel the
scheduled send (e.g. when the user types `/cancel` before the
reminder fires).

### `SendRecurring`

```go
func SendRecurring(parent context.Context, sender Sender, chatID, text string, interval time.Duration) (cancel func())
```

SendRecurring schedules `text` to be sent to `chatID` every
`interval`. The returned `cancel` function stops the recurring
send.
