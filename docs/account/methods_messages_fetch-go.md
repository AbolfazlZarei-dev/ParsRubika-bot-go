# `account` — Methods messages fetch

_Source: [`methods_messages_fetch.go`](../../ParsRubika/account/methods_messages_fetch.go)_

## Symbols

### `SearchChatMessages`

```go
func (a *Account) SearchChatMessages(ctx context.Context, objectGUID, searchText, searchType string) (*events.Event, error)
```

SearchChatMessages searches messages inside a single chat. The API
uses `{type: "Text"|"Hashtag"}` (no `from_user_id`).

### `GetMessages`

```go
func (a *Account) GetMessages(ctx context.Context, objectGUID, maxID string, limit int, sort string) (*events.Event, error)
```

GetMessages fetches a window of messages from an object chat.

### `GetMessagesByID`

```go
func (a *Account) GetMessagesByID(ctx context.Context, objectGUID string, messageIDs []string) (*events.Event, error)
```

GetMessagesByID fetches specific messages by ID.

### `GetMessagesInterval`

```go
func (a *Account) GetMessagesInterval(ctx context.Context, objectGUID, middleMessageID string) (*events.Event, error)
```

GetMessagesInterval fetches a range of messages around the supplied
middle message ID. The API uses `middle_message_id`, not `mid` /
`from_ticket_id`.

### `GetMessagesUpdates`

```go
func (a *Account) GetMessagesUpdates(ctx context.Context, objectGUID string, state ...string) (*events.Event, error)
```

GetMessagesUpdates returns the message updates after `state` (the
server-side cursor; omit it to start from the newest state).

### `GetMessageShareURL`

```go
func (a *Account) GetMessageShareURL(ctx context.Context, objectGUID, messageID string) (*events.Event, error)
```

GetMessageShareURL returns the share URL of a message.

### `GetMessageReactions`

```go
func (a *Account) GetMessageReactions(ctx context.Context, objectGUID, messageID string) (*events.Event, error)
```

GetMessageReactions returns the reactions attached to a message.

### `GetChatHistory`

```go
func (a *Account) GetChatHistory(ctx context.Context, objectGUID, maxID string, limit int) (*events.Event, error)
```

GetChatHistory is a convenience wrapper around GetMessages that
fetches the most recent `limit` messages from a chat, starting from
the given `maxID` (pass "" to start from the newest). It uses
`sort=FromMax` so messages are returned newest-first.

### `GetUnreadMentions`

```go
func (a *Account) GetUnreadMentions(ctx context.Context, objectGUID, maxID, minID string) (*events.Event, error)
```

GetUnreadMentions returns unread mentions for a chat.

### `SearchGlobalMessages`

```go
func (a *Account) SearchGlobalMessages(ctx context.Context, searchText, startID string) (*events.Event, error)
```

SearchGlobalMessages searches messages globally.

### `GetPinnedMessages`

```go
func (a *Account) GetPinnedMessages(ctx context.Context, objectGUID string) (*events.Event, error)
```

GetPinnedMessages returns the pinned messages of a chat. It works
for groups, channels and user (private) chats.

### `GetScheduledHistory`

```go
func (a *Account) GetScheduledHistory(ctx context.Context, objectGUID, startID string) (*events.Event, error)
```

GetScheduledHistory returns scheduled messages.

### `GetMessageURL`

```go
func (a *Account) GetMessageURL(ctx context.Context, objectGUID, messageID string) (string, error)
```

GetMessageURL returns the share URL of a message. It is an alias for
GetMessageShareURL with a clearer name.

### `GetChatMessages`

```go
func (a *Account) GetChatMessages(ctx context.Context, objectGUID, maxID string, limit int, sort string) (*events.Event, error)
```

GetChatMessages fetches messages from a chat. It is an alias for
GetMessages.

### `GetMessagesBatch`

```go
func (a *Account) GetMessagesBatch(ctx context.Context, objectGUID string, messageIDs []string) ([]map[string]any, error)
```

GetMessagesBatch fetches messages by IDs in batches of 20.

### `GetAllMessages`

```go
func (a *Account) GetAllMessages(ctx context.Context, objectGUID string, maxPages int) ([]map[string]any, error)
```

GetAllMessages fetches all messages in a chat (up to maxPages pages).
Set maxPages to 0 for unbounded (use with caution).

### `GetFirstMessages`

```go
func (a *Account) GetFirstMessages(ctx context.Context, objectGUID string) (*events.Event, error)
```

GetFirstMessages fetches the first batch of messages in a chat by
calling getMessagesInterval with an empty middleMessageID.

### `GetLastMessages`

```go
func (a *Account) GetLastMessages(ctx context.Context, objectGUID string, limit int) (*events.Event, error)
```

GetLastMessages fetches the most recent messages of a chat.

### `GetScheduledMessages`

```go
func (a *Account) GetScheduledMessages(ctx context.Context, objectGUID string) (*events.Event, error)
```

GetScheduledMessages returns the scheduled messages of a chat.
