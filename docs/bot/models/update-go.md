# `models` — Update

_Source: [`update.go`](../../../ParsRubika/bot/models/update.go)_

## Symbols

### `Update`

```go
type Update struct
```

Update is the high-level update dispatched by the bot's polling /
webhook loop.

### `MessageID`

```go
func (u *Update) MessageID() string
```

MessageID convenience accessor.

### `Text`

```go
func (u *Update) Text() string
```

Text returns the message text (if any).

### `InlineMessage`

```go
type InlineMessage struct
```

InlineMessage is the InlineMessage update type. It is dispatched by
the server when the bot is invoked in inline-mode.

### `FindKey`

```go
func (u *Update) FindKey(key string) any
```

FindKey is the convenience helper used by the rules to inspect nested
fields by name (for example "text", "message_id" or "chat_id").

### `FindKey`

```go
func (i *InlineMessage) FindKey(key string) any
```

FindKey is the inline-message counterpart of Update.FindKey.

### `Message`

```go
func (u *Update) Message() *Message
```

Message returns the message carried by the update (new or edited),
or nil for updates that carry none.

### `SenderID`

```go
func (u *Update) SenderID() string
```

SenderID returns the ID of the user that produced the update.

### `ButtonID`

```go
func (u *Update) ButtonID() string
```

ButtonID returns the ID of the inline button that triggered the update.

### `IsPrivateChat`

```go
func (u *Update) IsPrivateChat() bool { return len(u.ChatID) > 1 && u.ChatID[:2] == "u0" }
```

IsPrivateChat reports whether the update comes from a one-to-one
(user-to-bot) chat. In Rubika, user object GUIDs start with "u0";
"b0" is the bot GUID prefix, not a private-chat marker. Groups use
"g0" and channels use "c0".
