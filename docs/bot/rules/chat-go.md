# `rules` — Chat

_Source: [`chat.go`](../../../ParsRubika/bot/rules/chat.go)_

## Symbols

### `Private`

```go
type Private struct
```

Private matches when the update came from a private (user-to-bot)
chat. User object GUIDs start with "u0".

### `Check`

```go
func (Private) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `Group`

```go
type Group struct
```

Group matches when the update came from a group.

### `Check`

```go
func (Group) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `Channel`

```go
type Channel struct
```

Channel matches when the update came from a channel.

### `Check`

```go
func (Channel) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `BotSender`

```go
type BotSender struct
```

BotSender matches messages sent by a bot.

### `Check`

```go
func (BotSender) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `ChatID`

```go
type ChatID struct
```

ChatID matches by chat ID.

### `NewChatID`

```go
func NewChatID(ids ...string) *ChatID { return &ChatID{IDs: ids} }
```

NewChatID constructs a ChatID rule.

### `Check`

```go
func (c *ChatID) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `UpdateType`

```go
type UpdateType struct
```

UpdateType matches when the update type is one of `types`.

### `NewUpdateType`

```go
func NewUpdateType(types ...models.UpdateType) *UpdateType { return &UpdateType{Types: types} }
```

NewUpdateType constructs an UpdateType rule.

### `Check`

```go
func (u *UpdateType) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `Button`

```go
type Button struct
```

Button matches a button click by `button_id`. Supports regex.

### `NewButton`

```go
func NewButton(buttonID string, regex bool) *Button
```

NewButton constructs a Button rule.

### `Check`

```go
func (b *Button) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `Forward`

```go
type Forward struct
```

Forward matches forwarded messages.

### `Check`

```go
func (Forward) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `File`

```go
type File struct
```

File matches messages with a file attachment.

### `Check`

```go
func (File) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `Sticker`

```go
type Sticker struct
```

Sticker matches sticker messages.

### `Check`

```go
func (Sticker) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `Location`

```go
type Location struct
```

Location matches messages with a location attachment.

### `Check`

```go
func (Location) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `Poll`

```go
type Poll struct
```

Poll matches messages with a poll attachment.

### `Check`

```go
func (Poll) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `Contact`

```go
type Contact struct
```

Contact matches messages with a contact attachment.

### `Check`

```go
func (Contact) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `Reply`

```go
type Reply struct
```

Reply matches messages that reply to a given message ID.

### `NewReply`

```go
func NewReply(toMessageID string) *Reply { return &Reply{ToMessageID: toMessageID} }
```

NewReply constructs a Reply rule.

### `Check`

```go
func (r *Reply) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.
