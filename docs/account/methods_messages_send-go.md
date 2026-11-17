# `account` — Methods messages send

_Source: [`methods_messages_send.go`](../../ParsRubika/account/methods_messages_send.go)_

## Symbols

### `SendOptions`

```go
type SendOptions struct
```

SendOptions carries the optional arguments for SendMessage.

### `FileInlineInput`

```go
type FileInlineInput struct
```

FileInlineInput bundles everything SendMessage needs to know about
attaching a file to the outgoing message.

### `SendMessage`

```go
func (a *Account) SendMessage(ctx context.Context, objectGUID, msg string, opts ...SendOptions) (*events.Event, error)
```

SendMessage sends a message to the given object GUID.

### `SendPhoto`

```go
func (a *Account) SendPhoto(ctx context.Context, objectGUID string, photo any, caption string, opts ...SendOptions) (*events.Event, error)
```

SendPhoto sends a photo.

### `SendVideo`

```go
func (a *Account) SendVideo(ctx context.Context, objectGUID string, video any, caption string, opts ...SendOptions) (*events.Event, error)
```

SendVideo sends a video.

### `SendVoice`

```go
func (a *Account) SendVoice(ctx context.Context, objectGUID string, voice any, caption string, opts ...SendOptions) (*events.Event, error)
```

SendVoice sends a voice message.

### `SendMusic`

```go
func (a *Account) SendMusic(ctx context.Context, objectGUID string, music any, caption string, opts ...SendOptions) (*events.Event, error)
```

SendMusic sends an audio file.

### `SendGif`

```go
func (a *Account) SendGif(ctx context.Context, objectGUID string, gif any, caption string, opts ...SendOptions) (*events.Event, error)
```

SendGif sends a GIF.

### `SendDocument`

```go
func (a *Account) SendDocument(ctx context.Context, objectGUID string, file any, caption string, opts ...SendOptions) (*events.Event, error)
```

SendDocument sends a generic file.

### `SendVideoMessage`

```go
func (a *Account) SendVideoMessage(ctx context.Context, objectGUID string, video any, caption string, opts ...SendOptions) (*events.Event, error)
```

SendVideoMessage sends a round video (video-message).

### `SendSticker`

```go
func (a *Account) SendSticker(ctx context.Context, objectGUID, emojiCharacter, stickerID, stickerSetID, wHRatio string, stickerFile map[string]any, replyToMessageID string) (*events.Event, error)
```

SendSticker sends a sticker with the full sticker payload.

### `SendSimpleSticker`

```go
func (a *Account) SendSimpleSticker(ctx context.Context, objectGUID, stickerID, replyToMessageID string) (*events.Event, error)
```

SendSimpleSticker sends a sticker by sticker_id only.
