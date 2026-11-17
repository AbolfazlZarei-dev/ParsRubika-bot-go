# `handy` — Package handy is a collection of small, ergonomic helpers that

_Source: [`handy.go`](../../../ParsRubika/extras/handy/handy.go)_

## Symbols

### `Echo`

```go
func Echo(b *bot.Bot)
```

Echo registers a single handler that replies with the same text it
received. Useful for testing the bot.

### `PingReply`

```go
func PingReply(b *bot.Bot)
```

PingReply registers a single handler that replies "pong!" to any
message that equals "ping" (case-insensitive).

### `CommandHelp`

```go
func CommandHelp(b *bot.Bot, commands map[string]string)
```

CommandHelp registers a `/help` command that lists every command
the bot knows about. The list is built from the supplied
`(command, description)` pairs.

### `MentionChat`

```go
func MentionChat(b *bot.Bot)
```

MentionChat replies to every update with the chat's GUID. Useful
for debugging / onboarding.

### `MentionUserID`

```go
func MentionUserID(b *bot.Bot)
```

MentionUserID replies with the user's GUID on the message "userid".

### `CopyMessage`

```go
func CopyMessage(b *bot.Bot)
```

CopyMessage registers a handler that copies any received message
back to the same chat. Useful for echoing stickers, files, etc.
