# `account` — Methods groups

_Source: [`methods_groups.go`](../../ParsRubika/account/methods_groups.go)_

## Symbols

### `EditSlowMode`

```go
func (a *Account) EditSlowMode(ctx context.Context, groupGUID string, slowMode int) (*events.Event, error)
```

EditSlowMode sets the slow-mode interval for a group. `slowMode`
is in seconds (0 to disable).

### `EditGroupTimer`

```go
func (a *Account) EditGroupTimer(ctx context.Context, groupGUID string, timer int) (*events.Event, error)
```

EditGroupTimer is an alias for EditSlowMode.

### `LockGroup`

```go
func (a *Account) LockGroup(ctx context.Context, groupGUID string) (*events.Event, error)
```

LockGroup removes SendMessages from the default access list.

### `UnlockGroup`

```go
func (a *Account) UnlockGroup(ctx context.Context, groupGUID string) (*events.Event, error)
```

UnlockGroup adds SendMessages to the default access list.

### `GetGroupMessageReadParticipants`

```go
func (a *Account) GetGroupMessageReadParticipants(ctx context.Context, groupGUID, messageID string) (*events.Event, error)
```

GetGroupMessageReadParticipants returns who has read a specific message.

### `EditGroupInfoFull`

```go
func (a *Account) EditGroupInfoFull(ctx context.Context, groupGUID, title, description string, extraValues map[string]any, updatedParameters []string) (*events.Event, error)
```

EditGroupInfoFull is the full-featured version of EditGroupInfo that
accepts extra parameters (slow_mode, event_messages, reactions, etc.).

### `GetGroupInfo`

```go
func (a *Account) GetGroupInfo(ctx context.Context, groupGUID string) (*events.Event, error)
```

GetGroupInfo fetches information about a group.

### `AddGroup`

```go
func (a *Account) AddGroup(ctx context.Context, title, description string, memberGUIDs []string) (*events.Event, error)
```

AddGroup creates a new group.

### `RemoveGroup`

```go
func (a *Account) RemoveGroup(ctx context.Context, groupGUID string) (*events.Event, error)
```

RemoveGroup deletes a group owned by the current user.

### `EditGroupInfo`

```go
func (a *Account) EditGroupInfo(ctx context.Context, groupGUID, title, description string) (*events.Event, error)
```

EditGroupInfo edits a group's metadata. Pass an empty string for any
field you do not want to change. The `updatedParameters` list is built
automatically from the non-empty fields.

### `EditGroupTitle`

```go
func (a *Account) EditGroupTitle(ctx context.Context, groupGUID, title string) (*events.Event, error)
```

EditGroupTitle is a convenience wrapper that only updates the title.

### `EditGroupDescription`

```go
func (a *Account) EditGroupDescription(ctx context.Context, groupGUID, description string) (*events.Event, error)
```

EditGroupDescription is a convenience wrapper for the description only.

### `JoinGroup`

```go
func (a *Account) JoinGroup(ctx context.Context, link string) (*events.Event, error)
```

JoinGroup joins a group via its invite link. The link is normalised
to just the trailing slug before being sent as `hash_link`.

### `LeaveGroup`

```go
func (a *Account) LeaveGroup(ctx context.Context, groupGUID string) (*events.Event, error)
```

LeaveGroup leaves a group.

### `AddGroupMembers`

```go
func (a *Account) AddGroupMembers(ctx context.Context, groupGUID string, memberGUIDs []string) (*events.Event, error)
```

AddGroupMembers adds members to a group.

### `BanGroupMember`

```go
func (a *Account) BanGroupMember(ctx context.Context, groupGUID, userGUID string) (*events.Event, error)
```

BanGroupMember bans a group member (action = "Set").

### `UnbanGroupMember`

```go
func (a *Account) UnbanGroupMember(ctx context.Context, groupGUID, userGUID string) (*events.Event, error)
```

UnbanGroupMember unbans a group member (action = "Unset").

### `GetGroupAllMembers`

```go
func (a *Account) GetGroupAllMembers(ctx context.Context, groupGUID, startID, searchText string) (*events.Event, error)
```

GetGroupAllMembers returns the group member list with an optional
search-text filter and a pagination start-id.

### `GetGroupAdminMembers`

```go
func (a *Account) GetGroupAdminMembers(ctx context.Context, groupGUID, startID string) (*events.Event, error)
```

GetGroupAdminMembers returns the admin member list (paged).

### `GetGroupAdminAccessList`

```go
func (a *Account) GetGroupAdminAccessList(ctx context.Context, groupGUID, memberGUID string) (*events.Event, error)
```

GetGroupAdminAccessList returns the access rights of an admin.

### `SetGroupAdmin`

```go
func (a *Account) SetGroupAdmin(ctx context.Context, groupGUID, userGUID, action string, accessList []string) (*events.Event, error)
```

SetGroupAdmin promotes/demotes a group admin. `action` is "SetAdmin"
or "UnsetAdmin". `accessList` is the list of permissions (e.g.
["BanMember", "PinMessages"]).

### `PromoteGroupAdmin`

```go
func (a *Account) PromoteGroupAdmin(ctx context.Context, groupGUID, userGUID string, accessList []string) (*events.Event, error)
```

PromoteGroupAdmin is a convenience wrapper that promotes a member.

### `DemoteGroupAdmin`

```go
func (a *Account) DemoteGroupAdmin(ctx context.Context, groupGUID, userGUID string) (*events.Event, error)
```

DemoteGroupAdmin is a convenience wrapper that demotes a member.

### `GetGroupDefaultAccess`

```go
func (a *Account) GetGroupDefaultAccess(ctx context.Context, groupGUID string) (*events.Event, error)
```

GetGroupDefaultAccess returns the default access rights of a group.

### `SetGroupDefaultAccess`

```go
func (a *Account) SetGroupDefaultAccess(ctx context.Context, groupGUID string, accessList []string) (*events.Event, error)
```

SetGroupDefaultAccess updates the default access rights of a group.

### `GetGroupLink`

```go
func (a *Account) GetGroupLink(ctx context.Context, groupGUID string) (*events.Event, error)
```

GetGroupLink returns the group's invite link.

### `SetGroupLink`

```go
func (a *Account) SetGroupLink(ctx context.Context, groupGUID string) (*events.Event, error)
```

SetGroupLink regenerates the group's invite link.

### `GroupPreviewByJoinLink`

```go
func (a *Account) GroupPreviewByJoinLink(ctx context.Context, link string) (*events.Event, error)
```

GroupPreviewByJoinLink fetches the preview before joining via a link.
The link is normalised to its trailing slug before being sent as
`hash_link`.

### `GetGroupOnlineCount`

```go
func (a *Account) GetGroupOnlineCount(ctx context.Context, groupGUID string) (*events.Event, error)
```

GetGroupOnlineCount returns the number of online members.

### `GetGroupMentionList`

```go
func (a *Account) GetGroupMentionList(ctx context.Context, groupGUID, searchMention string) (*events.Event, error)
```

GetGroupMentionList returns the mention list of a group, optionally
filtered by `searchMention`.

### `GetBannedGroupMembers`

```go
func (a *Account) GetBannedGroupMembers(ctx context.Context, groupGUID, startID string) (*events.Event, error)
```

GetBannedGroupMembers returns the list of banned members (paged).
