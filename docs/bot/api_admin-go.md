# `bot` — Api admin

_Source: [`api_admin.go`](../../ParsRubika/bot/api_admin.go)_

## Symbols

### `SendQuiz`

```go
func (b *Bot) SendQuiz(ctx context.Context, chatID, question string, options []string, correctOptionIdx int, explanation string, replyTo string, disableNotif bool, keypad *models.Keypad, keypadType models.ChatKeypadType) (*models.MessageID, error)
```

SendQuiz is a convenience wrapper around SendPoll that sets type=Quiz
and requires a correct option index.

### `ClosePoll`

```go
func (b *Bot) ClosePoll(ctx context.Context, chatID, messageID string) error
```

ClosePoll closes/stops a poll so no further votes are accepted. The
Rubika bot API exposes this as a dedicated endpoint in newer versions;
when not available, use the account-side setPollAction.

### `GetChatAdmins`

```go
func (b *Bot) GetChatAdmins(ctx context.Context, chatID string) ([]map[string]any, error)
```

GetChatAdmins is an alias for GetChatAdministrators (common naming).

### `GetChatMembers`

```go
func (b *Bot) GetChatMembers(ctx context.Context, chatID, startID string) (map[string]any, error)
```

GetChatMembers returns a page of chat members. startID is the
pagination cursor (empty for the first page).

### `GetChatInfo`

```go
func (b *Bot) GetChatInfo(ctx context.Context, chatID string) (*models.Chat, error)
```

GetChatInfo is an alias for GetChat.

### `SetChatTitle`

```go
func (b *Bot) SetChatTitle(ctx context.Context, chatID, title string) error
```

SetChatTitle updates a group/channel title.

### `SetChatDescription`

```go
func (b *Bot) SetChatDescription(ctx context.Context, chatID, description string) error
```

SetChatDescription updates a group/channel description.

### `SetChatPhoto`

```go
func (b *Bot) SetChatPhoto(ctx context.Context, chatID, fileID string) error
```

SetChatPhoto sets the chat avatar from an uploaded file_id.

### `RemoveChatPhoto`

```go
func (b *Bot) RemoveChatPhoto(ctx context.Context, chatID string) error
```

RemoveChatPhoto removes the chat avatar.

### `AddChatMembers`

```go
func (b *Bot) AddChatMembers(ctx context.Context, chatID string, userIDs []string) error
```

AddChatMembers adds multiple users to a group/channel.

### `BanChatMembers`

```go
func (b *Bot) BanChatMembers(ctx context.Context, chatID string, userIDs []string) error
```

BanChatMembers bans multiple members sequentially. Returns the first
error encountered (if any). Successfully banned members before the
error are not rolled back.

### `UnbanChatMembers`

```go
func (b *Bot) UnbanChatMembers(ctx context.Context, chatID string, userIDs []string) error
```

UnbanChatMembers unbans multiple members sequentially.

### `RestrictChatMember`

```go
func (b *Bot) RestrictChatMember(ctx context.Context, chatID, userID string, until int64) error
```

RestrictChatMember restricts a member in a chat until the given time
(Unix seconds; 0 or negative means forever).

### `PromoteChatMember`

```go
func (b *Bot) PromoteChatMember(ctx context.Context, chatID, userID string, rights []string) error
```

PromoteChatMember elevates a member to admin with the supplied rights.
rights is a list of permission strings (e.g. "SendMessages", "PinMessages").

### `DemoteChatMember`

```go
func (b *Bot) DemoteChatMember(ctx context.Context, chatID, userID string) error
```

DemoteChatMember removes admin rights from a member.

### `ExportChatInviteLink`

```go
func (b *Bot) ExportChatInviteLink(ctx context.Context, chatID string) (string, error)
```

ExportChatInviteLink generates (or re-exports) the primary invite link.

### `RevokeChatInviteLink`

```go
func (b *Bot) RevokeChatInviteLink(ctx context.Context, chatID, link string) error
```

RevokeChatInviteLink revokes an existing invite link.

### `CreateGroup`

```go
func (b *Bot) CreateGroup(ctx context.Context, title string, memberUserIDs []string) (map[string]any, error)
```

CreateGroup creates a new group with the given title and initial members.

### `CreateChannel`

```go
func (b *Bot) CreateChannel(ctx context.Context, title, description string) (map[string]any, error)
```

CreateChannel creates a new channel.

### `LeaveChat`

```go
func (b *Bot) LeaveChat(ctx context.Context, chatID string) error
```

LeaveChat leaves a group or channel.

### `RemoveChatKeypad`

```go
func (b *Bot) RemoveChatKeypad(ctx context.Context, chatID string) error
```

RemoveChatKeypad removes the reply keypad for a chat (alias for
EditChatKeypad with type=Remove).

### `RegisterAllEndpoints`

```go
func (b *Bot) RegisterAllEndpoints(ctx context.Context, url string) error
```

RegisterAllEndpoints registers all 5 webhook endpoint types in one
call.

### `MemberIsAdmin`

```go
func (b *Bot) MemberIsAdmin(ctx context.Context, chatID, userID string) (bool, error)
```

MemberIsAdmin returns true if the user is an admin or the chat owner.

### `MemberIsOwner`

```go
func (b *Bot) MemberIsOwner(ctx context.Context, chatID, userID string) (bool, error)
```

MemberIsOwner returns true if the user is the chat owner/creator.

### `MemberIsJoined`

```go
func (b *Bot) MemberIsJoined(ctx context.Context, chatID, userID string) (bool, error)
```

MemberIsJoined returns true if the user is currently a member of the
chat (i.e. getChatMember returns a non-empty result).

### `CheckJoin`

```go
func (b *Bot) CheckJoin(ctx context.Context, chatID, userID string) (bool, error)
```

CheckJoin verifies that a user has joined the given channel/group. It
fetches the member list and scans for userID. This is a best-effort
check — very large chats may not return all members in one page.

### `GetAvatarMe`

```go
func (b *Bot) GetAvatarMe(ctx context.Context, saveAs string, asBytes bool, chunkSize int, progress func(downloaded, total int64)) (any, error)
```

GetAvatarMe downloads the bot's own avatar. Returns the saved path or
raw bytes depending on saveAs/asBytes (same semantics as DownloadFile).

### `CancelFunc`

```go
type CancelFunc func()
```

CancelFunc aborts a pending background action.

### `AutoDeleteMessage`

```go
func (b *Bot) AutoDeleteMessage(ctx context.Context, chatID, messageID string, delay time.Duration) CancelFunc
```

AutoDeleteMessage schedules a message for deletion after the given
delay. The returned CancelFunc aborts the pending delete if called
before the delay elapses.

### `DeleteAfter`

```go
func (b *Bot) DeleteAfter(ctx context.Context, chatID, messageID string, delay time.Duration) CancelFunc
```

DeleteAfter is an alias for AutoDeleteMessage with a clearer name.

### `EditAfter`

```go
func (b *Bot) EditAfter(ctx context.Context, chatID, messageID, newText string, delay time.Duration) CancelFunc
```

EditAfter schedules a message edit after the given delay.

### `DisableInlineKeyboard`

```go
func (b *Bot) DisableInlineKeyboard(ctx context.Context, chatID, messageID string, delay time.Duration) CancelFunc
```

DisableInlineKeyboard removes the inline keypad of a message after an
optional delay (0 = immediately).

### `SendScheduledMessage`

```go
func (b *Bot) SendScheduledMessage(ctx context.Context, chatID, text string, delay time.Duration, opts ...SendOptions) CancelFunc
```

SendScheduledMessage sleeps for delay then sends the message. Returns
the eventual *MessageID via the callback (non-blocking).

### `BulkResult`

```go
type BulkResult struct
```

BulkResult is the outcome of a single bulk operation target.

### `SendBulkMessage`

```go
func (b *Bot) SendBulkMessage(ctx context.Context, chatIDs []string, text string, concurrency int, delayBetween time.Duration, opts ...SendOptions) []BulkResult
```

SendBulkMessage sends the same text to multiple chats concurrently.
concurrency limits the number of in-flight sends. delayBetween is the
pause between launches (0 = no pause). Returns per-chat results.

### `DeleteBulkMessage`

```go
func (b *Bot) DeleteBulkMessage(ctx context.Context, chatID string, messageIDs []string, concurrency int) []BulkResult
```

DeleteBulkMessage deletes multiple messages from the same chat
concurrently. messageIDs must all belong to chatID.

### `EditBulkMessage`

```go
func (b *Bot) EditBulkMessage(ctx context.Context, chatID string, messages map[string]string, concurrency int) []BulkResult
```

EditBulkMessage edits multiple messages in the same chat. messages maps
messageID → new text.

### `SendButtonJoin`

```go
func (b *Bot) SendButtonJoin(ctx context.Context, chatID, title, username, text string, replyTo string) (*models.MessageID, error)
```

SendButtonJoin sends a message with a "Join Channel" inline button.

### `SendButtonURL`

```go
func (b *Bot) SendButtonURL(ctx context.Context, chatID, title, url, text string, replyTo string) (*models.MessageID, error)
```

SendButtonURL sends a message with a URL inline button.
