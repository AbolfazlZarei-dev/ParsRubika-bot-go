# `account` — Methods helpers

_Source: [`methods_helpers.go`](../../ParsRubika/account/methods_helpers.go)_

## Symbols

### `String`

```go
func (a *Account) String() string
```

String returns a human-readable description of the account.

### `BanMembers`

```go
func (a *Account) BanMembers(ctx context.Context, objectGUID string, memberGUIDs []string) error
```

BanMembers bans multiple members from a group or channel. The GUID
prefix (g0/c0) selects the correct underlying API.

### `UnbanMembers`

```go
func (a *Account) UnbanMembers(ctx context.Context, objectGUID string, memberGUIDs []string) error
```

UnbanMembers unbans multiple members from a group or channel.

### `UnbanAllMembers`

```go
func (a *Account) UnbanAllMembers(ctx context.Context, objectGUID string, excludeGUIDs []string) error
```

UnbanAllMembers fetches the banned-members list and unbans everyone.
excludeGUIDs (if non-nil) lists members that should remain banned.
The banned-members API returns members under "in_chat_members" with
each member's GUID in the "user_guid" field (not "object_guid").

### `ClearBlackList`

```go
func (a *Account) ClearBlackList(ctx context.Context, objectGUID string, excludeGUIDs []string) error
```

ClearBlackList is an alias for UnbanAllMembers.

### `AcceptAllJoinRequests`

```go
func (a *Account) AcceptAllJoinRequests(ctx context.Context, objectGUID string, excludeGUIDs []string) error
```

AcceptAllJoinRequests accepts every pending join request for a group
or channel. excludeGUIDs (if non-nil) lists users to skip.

### `GetMessageInfo`

```go
func (a *Account) GetMessageInfo(ctx context.Context, objectGUID, messageID string) (*events.Event, error)
```

GetMessageInfo fetches a single message by ID (delegates to
GetMessagesByID with a single-element slice).

### `GetLastMessage`

```go
func (a *Account) GetLastMessage(ctx context.Context, objectGUID string) (*events.Event, error)
```

GetLastMessage returns the last message in a chat. Rubika chat-info
responses nest the last message under the "last_message" key (an
object containing "message_id"), not as a top-level "last_message_id"
string. We extract the message_id from that nested object.

### `GetFirstMessage`

```go
func (a *Account) GetFirstMessage(ctx context.Context, objectGUID string) (*events.Event, error)
```

GetFirstMessage returns the first message in a chat by fetching the
messages interval and picking the lowest message_id. Message IDs are
compared numerically (not lexicographically) to handle IDs with
different digit counts correctly.

### `SearchMessages`

```go
func (a *Account) SearchMessages(ctx context.Context, objectGUID, query string, maxPages int) ([]map[string]any, error)
```

SearchMessages searches messages in a chat by text (client-side filter
over getMessages). It fetches up to maxPages pages of 50 messages each and returns those
whose text contains the query (case-insensitive). The cursor advances
to the OLDEST message in each page (msgs[len-1]) so each page fetches
a distinct window of messages.

### `DeleteMyMessages`

```go
func (a *Account) DeleteMyMessages(ctx context.Context, objectGUID string) error
```

DeleteMyMessages deletes all of the current user's messages in a chat
by fetching and deleting them in batches. The author is identified by
the "author_object_guid" field (not "author_guid"). The cursor advances
to the OLDEST message in each page so the next call fetches an earlier
window — using the newest would re-fetch nearly the same page and
loop ~forever.

### `DeleteAllMessages`

```go
func (a *Account) DeleteAllMessages(ctx context.Context, objectGUID string) error
```

DeleteAllMessages deletes every message in a chat (regardless of
author). Use with caution. The cursor advances to the OLDEST message
in each page (see DeleteMyMessages for rationale).

### `GetGUID`

```go
func (a *Account) GetGUID(ctx context.Context, input string) (string, error)
```

GetGUID resolves a username, join-link, or GUID to a GUID. If the
input is already a GUID (starts with u0/g0/c0/b0/s0) it is returned
as-is. Join-links (rubika.ir/joing/... or /joinc/...) are resolved
via groupPreviewByJoinLink / channelPreviewByJoinLink.

### `GetMyGUID`

```go
func (a *Account) GetMyGUID() string { return a.guid }
```

GetMyGUID returns the current user's GUID.

### `CheckJoin`

```go
func (a *Account) CheckJoin(ctx context.Context, objectGUID, userGUID string) (bool, error)
```

CheckJoin verifies that a user has joined a group/channel by fetching
the member list and scanning for the user's GUID.

### `BulkSendResult`

```go
type BulkSendResult struct
```

BulkSendResult is the outcome of a single bulk send target.

### `SendBulkMessage`

```go
func (a *Account) SendBulkMessage(ctx context.Context, objectGUIDs []string, text string, concurrency int, delayBetween time.Duration, opts ...SendOptions) []BulkSendResult
```

SendBulkMessage sends the same text to multiple chats concurrently.
concurrency limits in-flight sends.
