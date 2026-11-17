# `bot` — Api convenience

_Source: [`api_convenience.go`](../../ParsRubika/bot/api_convenience.go)_

## Symbols

### `SendImage`

```go
func (b *Bot) SendImage(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error)
```

SendImage sends an image file. Either `filePath` (local file to
upload) or `fileID` (previously uploaded file ID) must be set.
`text` is the optional caption.

### `SendVideo`

```go
func (b *Bot) SendVideo(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error)
```

SendVideo sends a video file. Either `filePath` or `fileID` must be
set. `text` is the optional caption.

### `SendGif`

```go
func (b *Bot) SendGif(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error)
```

SendGif sends a GIF file. Either `filePath` or `fileID` must be set.
`text` is the optional caption.

### `SendAudio`

```go
func (b *Bot) SendAudio(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error)
```

SendAudio sends an audio file. Either `filePath` or `fileID` must be
set. `text` is the optional caption.

### `SendMusic`

```go
func (b *Bot) SendMusic(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error)
```

SendMusic is an alias for SendAudio (kept for backwards
compatibility). It sends an audio file as a music message.

### `SendVoice`

```go
func (b *Bot) SendVoice(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error)
```

SendVoice sends a voice message. The file name is forced to `.ogg`
so the client renders it as a voice player.

### `SendMessageSimple`

```go
func (b *Bot) SendMessageSimple(ctx context.Context, chatID, text string) (*models.MessageID, error)
```

SendMessageSimple is a one-liner for SendMessage with just text.

### `SendDocument`

```go
func (b *Bot) SendDocument(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error)
```

SendDocument sends a generic document/file. It is an alias for SendFile
with FileTypeFile.

### `SendVideoMessage`

```go
func (b *Bot) SendVideoMessage(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error)
```

SendVideoMessage sends a round video message (VideoMessage type).
On the bot API side this is sent as a Video file; the round rendering
is determined by the client.

### `SendContactReply`

```go
func (b *Bot) SendContactReply(ctx context.Context, upd *models.Update, firstName, lastName, phone string) (*models.MessageID, error)
```

SendContactReply sends a contact as a reply to an update.

### `SendLocationReply`

```go
func (b *Bot) SendLocationReply(ctx context.Context, upd *models.Update, latitude, longitude float64) (*models.MessageID, error)
```

SendLocationReply sends a location as a reply to an update.

### `SendPollReply`

```go
func (b *Bot) SendPollReply(ctx context.Context, upd *models.Update, question string, options []string, pollType string) (*models.MessageID, error)
```

SendPollReply sends a poll as a reply to an update.

### `SendStickerReply`

```go
func (b *Bot) SendStickerReply(ctx context.Context, upd *models.Update, stickerID string) (*models.MessageID, error)
```

SendStickerReply sends a sticker as a reply to an update.

### `Reply`

```go
func Reply(ctx context.Context, b *Bot, u *models.Update, text string, opts ...SendOptions) (*models.MessageID, error)
```

Reply is a convenience helper that replies to a given update. It is
a free function (because we cannot add methods to a non-local type
from another package) — use it as `bot.Reply(ctx, b, upd, "text")`.

### `ReplyFile`

```go
func ReplyFile(ctx context.Context, b *Bot, u *models.Update, filePath, text string, fileType models.FileType, opts ...SendOptions) (*models.MessageID, error)
```

ReplyFile is a convenience helper that replies with a file.

### `Delete`

```go
func Delete(ctx context.Context, b *Bot, u *models.Update) error
```

Delete is a convenience helper that deletes the update's message.

### `EditText`

```go
func EditText(ctx context.Context, b *Bot, u *models.Update, newText string, opts ...EditOptions) error
```

EditText is a convenience helper that edits the update's message.

### `EditKeypad`

```go
func EditKeypad(ctx context.Context, b *Bot, u *models.Update, keypad *models.Keypad) error
```

EditKeypad is a convenience helper that edits the update's keypad.

### `IsGroupUpdate`

```go
func IsGroupUpdate(u *models.Update) bool
```

IsGroupUpdate reports whether the update came from a group chat.

### `IsChannelUpdate`

```go
func IsChannelUpdate(u *models.Update) bool
```

IsChannelUpdate reports whether the update came from a channel.

### `IsPrivateUpdate`

```go
func IsPrivateUpdate(u *models.Update) bool
```

IsPrivateUpdate reports whether the update came from a private
(user-to-bot) chat. In Rubika, user object GUIDs start with "u0";
"b0" is the bot GUID prefix, not a private-chat marker.
