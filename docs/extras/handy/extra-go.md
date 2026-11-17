# `handy` — Extra

_Source: [`extra.go`](../../../ParsRubika/extras/handy/extra.go)_

## Symbols

### `StartCommand`

```go
func StartCommand(b *bot.Bot, welcomeMessage string)
```

StartCommand registers a `/start` handler that sends a welcome message.
The welcome message can include the user's name via {name} placeholder.

### `IDCommand`

```go
func IDCommand(b *bot.Bot)
```

IDCommand registers an `/id` command that returns the chat ID and
sender ID.

### `TimeCommand`

```go
func TimeCommand(b *bot.Bot)
```

TimeCommand registers a `/time` command that returns the current server
time.

### `PingCommand`

```go
func PingCommand(b *bot.Bot)
```

PingCommand registers a `/ping` command that measures response time
and reports it back to the user.

### `UserInfoCommand`

```go
func UserInfoCommand(b *bot.Bot)
```

UserInfoCommand registers a `/whoami` command that returns the user's
info.

### `ReplyToMedia`

```go
func ReplyToMedia(b *bot.Bot)
```

ReplyToMedia registers a handler that replies to any media message
with the file ID.

### `AdminOnly`

```go
func AdminOnly(b *bot.Bot, rejectMessage string, handler func(ctx context.Context, b *bot.Bot, upd *bot.Update) error) func(ctx context.Context, b *bot.Bot, upd *bot.Update) error
```

AdminOnly wraps a handler so it only runs for chat admins. Non-admins
receive the rejectMessage (if non-empty).

### `PrivateOnly`

```go
func PrivateOnly(handler func(ctx context.Context, b *bot.Bot, upd *bot.Update) error) func(ctx context.Context, b *bot.Bot, upd *bot.Update) error
```

PrivateOnly wraps a handler so it only runs in private chats.

### `GroupOnly`

```go
func GroupOnly(handler func(ctx context.Context, b *bot.Bot, upd *bot.Update) error) func(ctx context.Context, b *bot.Bot, upd *bot.Update) error
```

GroupOnly wraps a handler so it only runs in groups.

### `ChannelOnly`

```go
func ChannelOnly(handler func(ctx context.Context, b *bot.Bot, upd *bot.Update) error) func(ctx context.Context, b *bot.Bot, upd *bot.Update) error
```

ChannelOnly wraps a handler so it only runs in channels.

### `SendTyping`

```go
func SendTyping(ctx context.Context, b *bot.Bot, chatID string) error
```

SendTyping sends a "typing" chat action. This is a no-op for the bot
API (which does not support chat activities), but is kept so code can be shared
with the account client.

### `ForwardToChat`

```go
func ForwardToChat(ctx context.Context, b *bot.Bot, fromChatID, messageID, toChatID string) error
```

ForwardToChat forwards a message to another chat.

### `SendKeypad`

```go
func SendKeypad(ctx context.Context, b *bot.Bot, chatID, text string, buttons [][]string) error
```

SendKeypad sends a message with a simple inline keypad built from a
list of (id, text) button pairs arranged in rows.

### `QuickReply`

```go
func QuickReply(ctx context.Context, b *bot.Bot, upd *bot.Update, text string) (string, error)
```

QuickReply sends a text reply to an update and returns the message ID.

### `DeleteAfter`

```go
func DeleteAfter(ctx context.Context, b *bot.Bot, chatID, text string, delay time.Duration) (bot.CancelFunc, error)
```

DeleteAfter sends a message and schedules its deletion after delay.
