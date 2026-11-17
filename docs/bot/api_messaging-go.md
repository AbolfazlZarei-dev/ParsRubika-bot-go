# `bot` — Api messaging

_Source: [`api_messaging.go`](../../ParsRubika/bot/api_messaging.go)_

## Symbols

### `SendOptions`

```go
type SendOptions struct
```

SendOptions bundles the optional arguments for the various send* calls.

### `SendMessage`

```go
func (b *Bot) SendMessage(ctx context.Context, chatID, msg string, opts ...SendOptions) (*models.MessageID, error)
```

SendMessage sends a text message.

### `SendSticker`

```go
func (b *Bot) SendSticker(ctx context.Context, chatID, stickerID string, opts ...SendOptions) (*models.MessageID, error)
```

SendSticker sends a sticker by sticker_id.

### `SendFile`

```go
func (b *Bot) SendFile(ctx context.Context, chatID, filePath, fileID, msg, fileName string, fileType models.FileType, opts ...SendOptions) (*models.MessageID, error)
```

SendFile sends a file. Either `filePath` or `fileID` must be set.

### `SendPoll`

```go
func (b *Bot) SendPoll(ctx context.Context, chatID, question string, options []string, pollType string, allowMulti, anonymous bool, correctIdx *int, explanation, replyTo string, disableNotif bool, keypad *models.Keypad, keypadType models.ChatKeypadType) (*models.MessageID, error)
```

SendPoll sends a poll.

### `SendLocation`

```go
func (b *Bot) SendLocation(ctx context.Context, chatID string, latitude, longitude float64, keypad *models.Keypad, disableNotif bool, replyTo string, keypadType models.ChatKeypadType) (*models.MessageID, error)
```

SendLocation sends a location attachment.

### `SendContact`

```go
func (b *Bot) SendContact(ctx context.Context, chatID, firstName, lastName, phone string, keypad *models.Keypad, disableNotif bool, replyTo string, keypadType models.ChatKeypadType) (*models.MessageID, error)
```

SendContact sends a contact attachment.

### `EditOptions`

```go
type EditOptions struct
```

EditOptions carries optional arguments for edit operations.

### `EditMessageText`

```go
func (b *Bot) EditMessageText(ctx context.Context, chatID, messageID, msg string, opts ...EditOptions) error
```

EditMessageText edits a message's text.

### `EditMessageKeypad`

```go
func (b *Bot) EditMessageKeypad(ctx context.Context, chatID, messageID string, keypad *models.Keypad) error
```

EditMessageKeypad replaces the inline keypad of an existing message.

### `DeleteMessage`

```go
func (b *Bot) DeleteMessage(ctx context.Context, chatID, messageID string) error
```

DeleteMessage removes a message.
