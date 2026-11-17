# `account` — Methods messages actions

_Source: [`methods_messages_actions.go`](../../ParsRubika/account/methods_messages_actions.go)_

## Symbols

### `TranscribeVoice`

```go
func (a *Account) TranscribeVoice(ctx context.Context, objectGUID, messageID string) (*events.Event, error)
```

TranscribeVoice sends a voice for transcription.

### `GetTranscription`

```go
func (a *Account) GetTranscription(ctx context.Context, messageID, transcriptionID string) (*events.Event, error)
```

GetTranscription fetches the result of a previous transcription.
The API expects `{message_id, transcription_id}`.

### `MarkAsRead`

```go
func (a *Account) MarkAsRead(ctx context.Context, objectGUID string) error
```

MarkAsRead marks a chat as read by fetching the last message id and
calling seenChats with it. The seenChats API expects a map of
chat_guid → last_seen_message_id; passing 0 marks nothing.

### `MarkAsSeen`

```go
func (a *Account) MarkAsSeen(ctx context.Context, objectGUID string) error
```

MarkAsSeen is an alias for MarkAsRead.

### `ClickInlineButton`

```go
func (a *Account) ClickInlineButton(ctx context.Context, objectGUID, messageID, buttonID string) (*events.Event, error)
```

ClickInlineButton clicks an inline button by sending a sendMessageAPICall
with the button's ID.

### `EditMessage`

```go
func (a *Account) EditMessage(ctx context.Context, objectGUID, messageID, msg string, parseMode enums.ParseMode) (*events.Event, error)
```

EditMessage edits the text of an existing message.

### `DeleteMessages`

```go
func (a *Account) DeleteMessages(ctx context.Context, objectGUID string, messageIDs []string) (*events.Event, error)
```

DeleteMessages deletes messages. `deleteType` is "Global" (default,
deletes for everyone) or "Local" (deletes only for the current user).

### `DeleteMessagesLocal`

```go
func (a *Account) DeleteMessagesLocal(ctx context.Context, objectGUID string, messageIDs []string) (*events.Event, error)
```

DeleteMessagesLocal deletes messages only for the current user
(Local delete — other participants can still see them).

### `ForwardMessages`

```go
func (a *Account) ForwardMessages(ctx context.Context, fromObjectGUID string, messageIDs []string, toObjectGUID string, hideAuthor bool, isMute bool) (*events.Event, error)
```

ForwardMessages forwards messages between chats. When `hideAuthor`
is true, the original sender is hidden.

### `SetPinMessage`

```go
func (a *Account) SetPinMessage(ctx context.Context, objectGUID, messageID string) (*events.Event, error)
```

SetPinMessage pins a message (action = "Pin").

### `SetUnpinMessage`

```go
func (a *Account) SetUnpinMessage(ctx context.Context, objectGUID, messageID string) (*events.Event, error)
```

SetUnpinMessage unpins a message (action = "Unpin").

### `UnpinAllMessages`

```go
func (a *Account) UnpinAllMessages(ctx context.Context, objectGUID string) (*events.Event, error)
```

UnpinAllMessages unpins all messages in a chat.

### `AutoDeleteMessage`

```go
func (a *Account) AutoDeleteMessage(ctx context.Context, objectGUID, messageID string, delay time.Duration) error
```

AutoDeleteMessage schedules a message for deletion after `delay`.

### `AddReaction`

```go
func (a *Account) AddReaction(ctx context.Context, objectGUID, messageID string, reactionID int) (*events.Event, error)
```

AddReaction adds a reaction to a message. The API uses
`actionOnMessageReaction` with `action: "Add"` and an int `reaction_id`.

### `RemoveReaction`

```go
func (a *Account) RemoveReaction(ctx context.Context, objectGUID, messageID string, reactionID int) (*events.Event, error)
```

RemoveReaction removes a reaction from a message.

### `ClickMessageURL`

```go
func (a *Account) ClickMessageURL(ctx context.Context, objectGUID, messageID, linkURL string) (*events.Event, error)
```

ClickMessageURL clicks a link inside a message.

### `ClickKeyPadButton`

```go
func (a *Account) ClickKeyPadButton(ctx context.Context, objectGUID, messageID, buttonID string) (*events.Event, error)
```

ClickKeyPadButton clicks a chat keypad button by its ID.

### `SendMessageAPICall`

```go
func (a *Account) SendMessageAPICall(ctx context.Context, objectGUID, messageID, buttonID string) (*events.Event, error)
```

SendMessageAPICall clicks an inline button on the user-side by
sending a message with the button's aux_data.

### `FeedbackTranscription`

```go
func (a *Account) FeedbackTranscription(ctx context.Context, messageID, transcriptionID, feedbackType, feedbackText string) (*events.Event, error)
```

FeedbackTranscription sends feedback on a voice transcription.

### `SendNowScheduledMessage`

```go
func (a *Account) SendNowScheduledMessage(ctx context.Context, objectGUID, messageID string) (*events.Event, error)
```

SendNowScheduledMessage sends a scheduled message immediately.
