# `account` — Methods channels

_Source: [`methods_channels.go`](../../ParsRubika/account/methods_channels.go)_

## Symbols

### `GetChannelInfo`

```go
func (a *Account) GetChannelInfo(ctx context.Context, channelGUID string) (*events.Event, error)
```

GetChannelInfo fetches information about a channel.

### `AddChannel`

```go
func (a *Account) AddChannel(ctx context.Context, title, description, username string, channelType ...string) (*events.Event, error)
```

AddChannel creates a new channel owned by the current user. The
optional channelType ("Public" in maxrubika's default) is sent only
when provided, so existing calls keep their exact old payload.

### `RemoveChannel`

```go
func (a *Account) RemoveChannel(ctx context.Context, channelGUID string) (*events.Event, error)
```

RemoveChannel removes (deletes) a channel.

### `EditChannelInfo`

```go
func (a *Account) EditChannelInfo(ctx context.Context, channelGUID, title, description, username string) (*events.Event, error)
```

EditChannelInfo edits a channel's metadata. Pass an empty string for
any field you do not want to change. The `updated_parameters` list
is built automatically.

### `CheckChannelUsername`

```go
func (a *Account) CheckChannelUsername(ctx context.Context, username string) (*events.Event, error)
```

CheckChannelUsername checks whether a username is free for a channel.

### `UpdateChannelUsername`

```go
func (a *Account) UpdateChannelUsername(ctx context.Context, channelGUID, username string) (*events.Event, error)
```

UpdateChannelUsername changes the channel's username.

### `AddChannelMembers`

```go
func (a *Account) AddChannelMembers(ctx context.Context, channelGUID string, userGUIDs []string) (*events.Event, error)
```

AddChannelMembers adds members to a channel. The API expects
`member_guids` (not `members`).

### `GetChannelAllMembers`

```go
func (a *Account) GetChannelAllMembers(ctx context.Context, channelGUID, startID, searchText string, limit int) (*events.Event, error)
```

GetChannelAllMembers returns the channel member list with an
optional search-text filter and pagination start-id.

### `GetChannelAdminMembers`

```go
func (a *Account) GetChannelAdminMembers(ctx context.Context, channelGUID, startID string) (*events.Event, error)
```

GetChannelAdminMembers returns the admin member list (paged).

### `GetChannelAdminAccessList`

```go
func (a *Account) GetChannelAdminAccessList(ctx context.Context, channelGUID, memberGUID string) (*events.Event, error)
```

GetChannelAdminAccessList returns the access rights of an admin.
`memberGUID` is required.

### `GetChannelLink`

```go
func (a *Account) GetChannelLink(ctx context.Context, channelGUID string) (*events.Event, error)
```

GetChannelLink returns the invite link of a channel.

### `SetChannelLink`

```go
func (a *Account) SetChannelLink(ctx context.Context, channelGUID string) (*events.Event, error)
```

SetChannelLink generates a new invite link.

### `BanChannelMember`

```go
func (a *Account) BanChannelMember(ctx context.Context, channelGUID, userGUID string) (*events.Event, error)
```

BanChannelMember bans a member from a channel (action = "Set").

### `UnbanChannelMember`

```go
func (a *Account) UnbanChannelMember(ctx context.Context, channelGUID, userGUID string) (*events.Event, error)
```

UnbanChannelMember unbans a channel member (action = "Unset").

### `ChannelPreviewByJoinLink`

```go
func (a *Account) ChannelPreviewByJoinLink(ctx context.Context, link string) (*events.Event, error)
```

ChannelPreviewByJoinLink fetches the channel preview before joining.
The link is normalised to its trailing slug.

### `JoinChannelByLink`

```go
func (a *Account) JoinChannelByLink(ctx context.Context, link string) (*events.Event, error)
```

JoinChannelByLink joins a channel via its invite link. The API only
needs the `hash_link` — the channel GUID is implicit.

### `JoinChannelAction`

```go
func (a *Account) JoinChannelAction(ctx context.Context, channelGUID, action string) (*events.Event, error)
```

JoinChannelAction accepts/declines a join request. `action` is
"Join" or "Leave".

### `SeenChannelMessages`

```go
func (a *Account) SeenChannelMessages(ctx context.Context, channelGUID, minID, maxID string) (*events.Event, error)
```

SeenChannelMessages marks a range of channel messages as seen.
`minID` and `maxID` define the inclusive range.

### `GetChannelSeenCount`

```go
func (a *Account) GetChannelSeenCount(ctx context.Context, channelGUID, minID, maxID string) (*events.Event, error)
```

GetChannelSeenCount returns the seen count for a range of messages.

### `GetChannelStatistics`

```go
func (a *Account) GetChannelStatistics(ctx context.Context, channelGUID string) (*events.Event, error)
```

GetChannelStatistics returns channel-specific statistics. The
Rubika API exposes a dedicated `channelStatistics` method for this
(distinct from the account-wide `getAllStatistics`), so we call it
directly to avoid returning global stats when the caller asked for
channel-only data.

### `SetChannelAdmin`

```go
func (a *Account) SetChannelAdmin(ctx context.Context, channelGUID, memberGUID, action string, accessList []string) (*events.Event, error)
```

SetChannelAdmin promotes or demotes a channel admin.

### `PromoteChannelAdmin`

```go
func (a *Account) PromoteChannelAdmin(ctx context.Context, channelGUID, memberGUID string, accessList []string) (*events.Event, error)
```

PromoteChannelAdmin promotes a channel admin.

### `DemoteChannelAdmin`

```go
func (a *Account) DemoteChannelAdmin(ctx context.Context, channelGUID, memberGUID string) (*events.Event, error)
```

DemoteChannelAdmin demotes a channel admin.

### `SetChannelType`

```go
func (a *Account) SetChannelType(ctx context.Context, channelGUID, channelType string) (*events.Event, error)
```

SetChannelType sets the channel type (Public or Private).

### `EditChannelTitle`

```go
func (a *Account) EditChannelTitle(ctx context.Context, channelGUID, title string) (*events.Event, error)
```

EditChannelTitle is a convenience wrapper for editing just the title.
