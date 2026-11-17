# `bot` — Aliases

_Source: [`aliases.go`](../../ParsRubika/bot/aliases.go)_

## Symbols

### `NewText`

```go
func NewText(pattern string, regex bool) *rules.Text { return rules.NewText(pattern, regex) }
```

NewText re-exports rules.NewText.

### `NewCommand`

```go
func NewCommand(cmds, prefixes []string, caseInsensitive, allowNoPrefix bool) *rules.Command
```

NewCommand re-exports rules.NewCommand.

### `NewChatID`

```go
func NewChatID(ids ...string) *rules.ChatID { return rules.NewChatID(ids...) }
```

NewChatID re-exports rules.NewChatID.

### `NewUpdateType`

```go
func NewUpdateType(types ...models.UpdateType) *rules.UpdateType
```

NewUpdateType re-exports rules.NewUpdateType.

### `NewButton`

```go
func NewButton(buttonID string, regex bool) *rules.Button { return rules.NewButton(buttonID, regex) }
```

NewButton re-exports rules.NewButton.

### `NewReply`

```go
func NewReply(toMessageID string) *rules.Reply { return rules.NewReply(toMessageID) }
```

NewReply re-exports rules.NewReply.

### `ReplyUpdate`

```go
func ReplyUpdate(ctx context.Context, b *Bot, u *models.Update, msg string, opts ...SendOptions) (*models.MessageID, error)
```

ReplyUpdate is the convenience function that replies to an update.
Use it as `bot.ReplyUpdate(ctx, b, upd, "Hello!")`.

### `ReplyFileUpdate`

```go
func ReplyFileUpdate(ctx context.Context, b *Bot, u *models.Update, filePath, msg string, fileType models.FileType, opts ...SendOptions) (*models.MessageID, error)
```

ReplyFileUpdate is the convenience function that replies with a file.

### `DeleteUpdate`

```go
func DeleteUpdate(ctx context.Context, b *Bot, u *models.Update) error
```

DeleteUpdate deletes the update's message.

### `EditTextUpdate`

```go
func EditTextUpdate(ctx context.Context, b *Bot, u *models.Update, newText string, opts ...EditOptions) error
```

EditTextUpdate edits the update's message text.

### `EditKeypadUpdate`

```go
func EditKeypadUpdate(ctx context.Context, b *Bot, u *models.Update, keypad *models.Keypad) error
```

EditKeypadUpdate edits the update's inline keypad.
