# `bot` — Api chats

_Source: [`api_chats.go`](../../ParsRubika/bot/api_chats.go)_

## Symbols

### `GetMe`

```go
func (b *Bot) GetMe(ctx context.Context) (*models.Bot, error)
```

GetMe returns the bot identity.

### `SetCommands`

```go
func (b *Bot) SetCommands(ctx context.Context, cmds []models.BotCommand) error
```

SetCommands registers the bot's command list.

### `UpdateBotEndpoints`

```go
func (b *Bot) UpdateBotEndpoints(ctx context.Context, url, endpointType string) (map[string]any, error)
```

UpdateBotEndpoints registers or unregisters the webhook endpoints.

### `GetChat`

```go
func (b *Bot) GetChat(ctx context.Context, chatID string) (*models.Chat, error)
```

GetChat returns the chat object for `chatID`.

### `ForwardMessage`

```go
func (b *Bot) ForwardMessage(ctx context.Context, fromChatID, messageID, toChatID string, disableNotif bool) (*models.MessageID, error)
```

ForwardMessage forwards a message to another chat.

### `GetChatAdministrators`

```go
func (b *Bot) GetChatAdministrators(ctx context.Context, chatID string) ([]map[string]any, error)
```

GetChatAdministrators returns the admin list of a chat.

### `GetChatMemberCount`

```go
func (b *Bot) GetChatMemberCount(ctx context.Context, chatID string) (int, error)
```

GetChatMemberCount returns the member count of a chat.

### `BanChatMember`

```go
func (b *Bot) BanChatMember(ctx context.Context, chatID, userID string) error
```

BanChatMember bans a member.

### `UnbanChatMember`

```go
func (b *Bot) UnbanChatMember(ctx context.Context, chatID, userID string) error
```

UnbanChatMember revokes a ban.

### `GetChatMember`

```go
func (b *Bot) GetChatMember(ctx context.Context, chatID, userID string) (map[string]any, error)
```

GetChatMember fetches a single chat member.

### `PinChatMessage`

```go
func (b *Bot) PinChatMessage(ctx context.Context, chatID, messageID string, disableNotif bool) error
```

PinChatMessage pins a message.

### `UnpinChatMessage`

```go
func (b *Bot) UnpinChatMessage(ctx context.Context, chatID, messageID string) error
```

UnpinChatMessage unpins a message.

### `UnpinAllChatMessages`

```go
func (b *Bot) UnpinAllChatMessages(ctx context.Context, chatID string) error
```

UnpinAllChatMessages unpins all messages in a chat. The Rubika
API only accepts a chat_id for this operation.

### `EditChatKeypad`

```go
func (b *Bot) EditChatKeypad(ctx context.Context, chatID string, keypadType models.ChatKeypadType, keypad *models.Keypad) error
```

EditChatKeypad edits the chat keypad.

### `RequestSendFileURL`

```go
func (b *Bot) RequestSendFileURL(ctx context.Context, fileType models.FileType) (string, error)
```

RequestSendFileURL requests an upload URL for the supplied file type.

### `GetFile`

```go
func (b *Bot) GetFile(ctx context.Context, fileID string) (string, error)
```

GetFile returns the download URL of a previously uploaded file.

### `DownloadFile`

```go
func (b *Bot) DownloadFile(ctx context.Context, fileID, saveAs string, asBytes bool, chunkSize int, progress func(downloaded, total int64)) (any, error)
```

DownloadFile downloads a file by file_id. When `saveAs` is empty the
raw bytes are returned. The optional `progress` callback is invoked
after each chunk.
