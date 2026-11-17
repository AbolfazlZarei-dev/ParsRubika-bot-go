# `account` — Methods chats

_Source: [`methods_chats.go`](../../ParsRubika/account/methods_chats.go)_

## Symbols

### `GetChats`

```go
func (a *Account) GetChats(ctx context.Context, startID string) (*events.Event, error)
```

GetChats returns the user's chats.

### `GetChatsUpdates`

```go
func (a *Account) GetChatsUpdates(ctx context.Context) (*events.Event, error)
```

GetChatsUpdates fetches the latest chat updates.

The server requires a `state` parameter (Unix-seconds timestamp,
typically `now - 150`). Without it the server replies
`INVALID_INPUT: getChatsUpdates failed`, so the call always sends
`state = now - 150`.

### `SeenChats`

```go
func (a *Account) SeenChats(ctx context.Context, chatIDs []string) (*events.Event, error)
```

SeenChats marks the given chat GUIDs as seen.

### `GetAbsObjects`

```go
func (a *Account) GetAbsObjects(ctx context.Context, guids []string) (*events.Event, error)
```

GetAbsObjects returns the absolute object info for the given GUIDs.

### `SendChatActivity`

```go
func (a *Account) SendChatActivity(ctx context.Context, objectGUID string, action enums.ChatAction) (*events.Event, error)
```

SendChatActivity reports a chat activity. The API expects
`{activity: "Typing"}` (a string), not `{activity: {type: "Typing"}}`.

### `SetActionChat`

```go
func (a *Account) SetActionChat(ctx context.Context, objectGUID, action string) (*events.Event, error)
```

SetActionChat mutes / unmutes / archives / unarchives / pins / unpins
a chat. `action` is one of "Mute", "Unmute", "Archive", "UnArchive",
"Pin", "Unpin".

### `MuteChat`

```go
func (a *Account) MuteChat(ctx context.Context, objectGUID string) (*events.Event, error)
```

MuteChat mutes the supplied chat.

### `UnmuteChat`

```go
func (a *Account) UnmuteChat(ctx context.Context, objectGUID string) (*events.Event, error)
```

UnmuteChat unmutes the supplied chat.

### `ArchiveChat`

```go
func (a *Account) ArchiveChat(ctx context.Context, objectGUID string) (*events.Event, error)
```

ArchiveChat archives the supplied chat.

### `UnarchiveChat`

```go
func (a *Account) UnarchiveChat(ctx context.Context, objectGUID string) (*events.Event, error)
```

UnarchiveChat unarchives the supplied chat.

### `PinChat`

```go
func (a *Account) PinChat(ctx context.Context, objectGUID string) (*events.Event, error)
```

PinChat pins the supplied chat to the top of the chat list.

### `UnpinChat`

```go
func (a *Account) UnpinChat(ctx context.Context, objectGUID string) (*events.Event, error)
```

UnpinChat unpins the supplied chat.

### `DeleteChatHistory`

```go
func (a *Account) DeleteChatHistory(ctx context.Context, objectGUID, lastMessageID string) (*events.Event, error)
```

DeleteChatHistory clears the chat history for the supplied GUID.

### `GetLinkFromAppURL`

```go
func (a *Account) GetLinkFromAppURL(ctx context.Context, appURL string) (*events.Event, error)
```

GetLinkFromAppURL converts an `app.rubika.ir/...` URL into a deep link.

### `GetJoinLinks`

```go
func (a *Account) GetJoinLinks(ctx context.Context, objectGUID string) (*events.Event, error)
```

GetJoinLinks returns the join links of a chat.

### `CreateJoinLink`

```go
func (a *Account) CreateJoinLink(ctx context.Context, objectGUID, title string, expirationSec, usageLimit int, requestNeeded bool) (*events.Event, error)
```

CreateJoinLink creates a new join link for a chat. The API uses
`usage_limit` (not `uses_limit`) and a `request_needed` boolean.

### `GetJoinRequests`

```go
func (a *Account) GetJoinRequests(ctx context.Context, objectGUID, cursor string) (*events.Event, error)
```

GetJoinRequests returns the pending join requests.

### `ActionOnJoinRequest`

```go
func (a *Account) ActionOnJoinRequest(ctx context.Context, objectGUID, userGUID, action string) (*events.Event, error)
```

ActionOnJoinRequest accepts or declines a join request. The API uses
`{object_guid, object_type, user_guid, action}` where `object_type`
is derived from the GUID prefix ("Group" / "Channel") and `action`
is "Accept" or "Reject".

### `AcceptJoinRequest`

```go
func (a *Account) AcceptJoinRequest(ctx context.Context, objectGUID, userGUID string) (*events.Event, error)
```

AcceptJoinRequest accepts a join request.

### `RejectJoinRequest`

```go
func (a *Account) RejectJoinRequest(ctx context.Context, objectGUID, userGUID string) (*events.Event, error)
```

RejectJoinRequest rejects a join request.

### `SetChatProtectedContent`

```go
func (a *Account) SetChatProtectedContent(ctx context.Context, objectGUID string, enabled bool) (*events.Event, error)
```

SetChatProtectedContent toggles the "protected content" (restricted
content / no-forward) flag on a group or channel.

### `MuteChatFor`

```go
func (a *Account) MuteChatFor(ctx context.Context, objectGUID string, durationSeconds int) (*events.Event, error)
```

MuteChatFor mutes a chat for the given duration (in seconds). The
Rubika API requires a minimum of 1800 seconds (30 minutes).

### `DeleteNoAccessGroupChat`

```go
func (a *Account) DeleteNoAccessGroupChat(ctx context.Context, groupGUID string) (*events.Event, error)
```

DeleteNoAccessGroupChat deletes a group chat that the user cannot access.

### `SearchChats`

```go
func (a *Account) SearchChats(ctx context.Context, objectGUID, query string, limit int) (*events.Event, error)
```

SearchChats searches for messages within a specific chat. `query`
is the search term; `limit` caps the number of results (use 0 for
the server default).

### `GetChatFolder`

```go
func (a *Account) GetChatFolder(ctx context.Context, folderID string) (*events.Event, error)
```

GetChatFolder returns a single chat folder by its ID. This is the
per-folder counterpart of GetFolders.

### `GetChatsByID`

```go
func (a *Account) GetChatsByID(ctx context.Context, chatIDs []string) (*events.Event, error)
```

GetChatsByID returns info for multiple chat GUIDs at once.

### `ClearDraft`

```go
func (a *Account) ClearDraft(ctx context.Context, objectGUID string) (*events.Event, error)
```

ClearDraft clears the draft of a single chat.

### `ClearAllDrafts`

```go
func (a *Account) ClearAllDrafts(ctx context.Context) (*events.Event, error)
```

ClearAllDrafts clears all chat drafts.

### `GetAllDrafts`

```go
func (a *Account) GetAllDrafts(ctx context.Context) (*events.Event, error)
```

GetAllDrafts returns all chat drafts.

### `SetChatUseTime`

```go
func (a *Account) SetChatUseTime(ctx context.Context, objectGUID string, useTime int) (*events.Event, error)
```

SetChatUseTime sets the chat use-time hint (anti-rate-limit).

### `SetAutoDelete`

```go
func (a *Account) SetAutoDelete(ctx context.Context, objectGUIDs []string, autoDelete int) (*events.Event, error)
```

SetAutoDelete sets the auto-delete timer for chats. `autoDelete` is in
seconds (0 to disable).

### `GetChatAds`

```go
func (a *Account) GetChatAds(ctx context.Context, state int) (*events.Event, error)
```

GetChatAds returns chat advertisements.

### `ActionOnChatAds`

```go
func (a *Account) ActionOnChatAds(ctx context.Context, action, chatAdsID string) (*events.Event, error)
```

ActionOnChatAds performs an action on a chat ad (View/Click/DontShow/Report).

### `SetAskSpamAction`

```go
func (a *Account) SetAskSpamAction(ctx context.Context, objectGUID, action string) (*events.Event, error)
```

SetAskSpamAction sets the ask-spam action for a chat.

### `PinChatInFolder`

```go
func (a *Account) PinChatInFolder(ctx context.Context, objectGUID, folderID string) (*events.Event, error)
```

PinChatInFolder pins a chat inside a folder.

### `UnpinChatInFolder`

```go
func (a *Account) UnpinChatInFolder(ctx context.Context, objectGUID, folderID string) (*events.Event, error)
```

UnpinChatInFolder unpins a chat from a folder.

### `EditJoinLink`

```go
func (a *Account) EditJoinLink(ctx context.Context, objectGUID, hashLink, title string, expirationSec, usageLimit int, requestNeeded *bool) (*events.Event, error)
```

EditJoinLink edits an existing join link.

### `RevokeJoinLink`

```go
func (a *Account) RevokeJoinLink(ctx context.Context, objectGUID, hashLink string) (*events.Event, error)
```

RevokeJoinLink revokes a join link.

### `GetFolders`

```go
func (a *Account) GetFolders(ctx context.Context, lastState int) (*events.Event, error)
```

GetFolders returns the user's chat folders. `lastState` is the
server-time baseline (defaults to 0).

### `GetSuggestedFolders`

```go
func (a *Account) GetSuggestedFolders(ctx context.Context) (*events.Event, error)
```

GetSuggestedFolders returns the suggested folders.

### `DeleteFolder`

```go
func (a *Account) DeleteFolder(ctx context.Context, folderID string) (*events.Event, error)
```

DeleteFolder deletes a chat folder.

### `AddFolder`

```go
func (a *Account) AddFolder(ctx context.Context, name string, includeChats, excludeChats, includeChatTypes, excludeChatTypes []string, isAddToTop bool) (*events.Event, error)
```

AddFolder creates a new chat folder.

### `ReorderFolder`

```go
func (a *Account) ReorderFolder(ctx context.Context, folderIDs []string) (*events.Event, error)
```

ReorderFolder reorders the user's chat folders.
