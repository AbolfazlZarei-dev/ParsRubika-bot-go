# `account` — Methods aliases

_Source: [`methods_aliases.go`](../../ParsRubika/account/methods_aliases.go)_

## Symbols

### `GetChatInfo`

```go
func (a *Account) GetChatInfo(ctx context.Context, objectGUID string) (*events.Event, error)
```

GetChatInfo is a dispatcher that fetches chat info by GUID type.
It is an alias for GetInfo with a clearer name.

### `GetInfoByLink`

```go
func (a *Account) GetInfoByLink(ctx context.Context, link string) (*events.Event, error)
```

GetInfoByLink resolves a join link to a chat preview. It dispatches
to channelPreviewByJoinLink or groupPreviewByJoinLink based on the
URL path.

### `GetInfoByUsername`

```go
func (a *Account) GetInfoByUsername(ctx context.Context, username string) (*events.Event, error)
```

GetInfoByUsername resolves a @username to an object.

### `JoinChat`

```go
func (a *Account) JoinChat(ctx context.Context, input string) (*events.Event, error)
```

JoinChat joins a group or channel by link or GUID. It is a dispatcher
alias for Join.

### `LeaveChatByGUID`

```go
func (a *Account) LeaveChatByGUID(ctx context.Context, objectGUID string) (*events.Event, error)
```

LeaveChatByGUID leaves a chat by GUID. It is a dispatcher alias for
LeaveChat.

### `JoinChannelByGUID`

```go
func (a *Account) JoinChannelByGUID(ctx context.Context, channelGUID string) (*events.Event, error)
```

JoinChannelByGUID joins a channel by GUID (action=Join).

### `LeaveChannel`

```go
func (a *Account) LeaveChannel(ctx context.Context, channelGUID string) (*events.Event, error)
```

LeaveChannel leaves a channel (action=Leave).

### `ArchiveChannel`

```go
func (a *Account) ArchiveChannel(ctx context.Context, channelGUID string) (*events.Event, error)
```

ArchiveChannel archives a channel (action=Archive).

### `GetChannelAdmins`

```go
func (a *Account) GetChannelAdmins(ctx context.Context, channelGUID, startID string) (*events.Event, error)
```

GetChannelAdmins is an alias for GetChannelAdminMembers.

### `GetChannelMembers`

```go
func (a *Account) GetChannelMembers(ctx context.Context, channelGUID, startID, searchText string, limit int) (*events.Event, error)
```

GetChannelMembers is an alias for GetChannelAllMembers.

### `GetGroupAdmins`

```go
func (a *Account) GetGroupAdmins(ctx context.Context, groupGUID, startID string) (*events.Event, error)
```

GetGroupAdmins is an alias for GetGroupAdminMembers.

### `GetGroupMembers`

```go
func (a *Account) GetGroupMembers(ctx context.Context, groupGUID, startID, searchText string) (*events.Event, error)
```

GetGroupMembers is an alias for GetGroupAllMembers.

### `GetBannedMembers`

```go
func (a *Account) GetBannedMembers(ctx context.Context, objectGUID, startID string) (*events.Event, error)
```

GetBannedMembers fetches banned members from a group or channel.

### `SetAdmin`

```go
func (a *Account) SetAdmin(ctx context.Context, objectGUID, memberGUID string, accessList []string) (*events.Event, error)
```

SetAdmin promotes a user to admin in a group or channel (dispatcher).

### `UnsetAdmin`

```go
func (a *Account) UnsetAdmin(ctx context.Context, objectGUID, memberGUID string) (*events.Event, error)
```

UnsetAdmin demotes an admin in a group or channel (dispatcher).

### `GetAdminAccessList`

```go
func (a *Account) GetAdminAccessList(ctx context.Context, objectGUID, memberGUID string) (*events.Event, error)
```

GetAdminAccessList fetches the admin access list for a member.

### `ChangeOwner`

```go
func (a *Account) ChangeOwner(ctx context.Context, objectGUID, newOwnerGUID string) (*events.Event, error)
```

ChangeOwner requests ownership transfer of a group or channel.
This is an alias for RequestChangeOwner.

### `CancelChangeOwnerOp`

```go
func (a *Account) CancelChangeOwnerOp(ctx context.Context, objectGUID string) (*events.Event, error)
```

CancelChangeOwner cancels a pending ownership transfer.
This is an alias for the existing method.

### `SearchGlobalChats`

```go
func (a *Account) SearchGlobalChats(ctx context.Context, query string) (*events.Event, error)
```

SearchGlobalChats searches all chats globally. Alias for
SearchGlobalObjects.

### `SearchGlobalUsers`

```go
func (a *Account) SearchGlobalUsers(ctx context.Context, query string) (*events.Event, error)
```

SearchGlobalUsers searches globally for users.

### `SearchGlobalChannels`

```go
func (a *Account) SearchGlobalChannels(ctx context.Context, query string) (*events.Event, error)
```

SearchGlobalChannels searches globally for channels.

### `SearchGlobalGroups`

```go
func (a *Account) SearchGlobalGroups(ctx context.Context, query string) (*events.Event, error)
```

SearchGlobalGroups searches globally for groups.

### `StartBotByGUID`

```go
func (a *Account) StartBotByGUID(ctx context.Context, botGUID string) (*events.Event, error)
```

StartBotByGUID starts a bot by sending /start. Alias for StartBot.

### `StopBotByGUID`

```go
func (a *Account) StopBotByGUID(ctx context.Context, botGUID string) (*events.Event, error)
```

StopBotByGUID stops a bot. Alias for StopBot.

### `GetBotInfoByGUID`

```go
func (a *Account) GetBotInfoByGUID(ctx context.Context, botGUID string) (*events.Event, error)
```

GetBotInfoByGUID fetches bot info. Alias for GetBotInfo.

### `GetInlineBotsList`

```go
func (a *Account) GetInlineBotsList(ctx context.Context, botGUIDs []string) (*events.Event, error)
```

GetInlineBotsList fetches info for multiple inline bots.

### `GetServiceInfoByGUID`

```go
func (a *Account) GetServiceInfoByGUID(ctx context.Context, serviceGUID string) (*events.Event, error)
```

GetServiceInfoByGUID fetches service chat info.

### `DeleteServiceChatByGUID`

```go
func (a *Account) DeleteServiceChatByGUID(ctx context.Context, serviceGUID, lastMessageID string) (*events.Event, error)
```

DeleteServiceChatByGUID deletes a service chat.

### `DeleteBotChatByGUID`

```go
func (a *Account) DeleteBotChatByGUID(ctx context.Context, botGUID, lastMessageID string) (*events.Event, error)
```

DeleteBotChatByGUID deletes a bot chat.

### `GetInfo`

```go
func (a *Account) GetInfo(ctx context.Context, objectGUID string) (*events.Event, error)
```

GetInfo returns generic info for the supplied object — it dispatches
to getUserInfo / getGroupInfo / getChannelInfo / getBotInfo /
getServiceInfo based on the GUID prefix, or `getObjectByUsername`
when the input is a username. We check prefixes directly (rather than
relying on events.GUIDType) because GUIDType returns "User" for
anything that is not a group or channel — including bots (b0) and
service chats (s0), which would be mis-dispatched to getUserInfo.

### `Join`

```go
func (a *Account) Join(ctx context.Context, input string) (*events.Event, error)
```

Join joins an object. If the input looks like a link, it dispatches
to `joinGroup` / `joinChannelByLink`. Otherwise (GUID), it dispatches
to `joinChannelAction(channelGUID, "Join")`.

### `LeaveChat`

```go
func (a *Account) LeaveChat(ctx context.Context, objectGUID string) (*events.Event, error)
```

LeaveChat leaves a chat — it dispatches to `joinChannelAction(channel,
"Remove")` (channel) or `leaveGroup` (group).

### `BanMember`

```go
func (a *Account) BanMember(ctx context.Context, objectGUID, memberGUID string) (*events.Event, error)
```

BanMember bans a member from a chat — it dispatches based on GUID
prefix to `banGroupMember` / `banChannelMember`.

### `UnbanMember`

```go
func (a *Account) UnbanMember(ctx context.Context, objectGUID, memberGUID string) (*events.Event, error)
```

UnbanMember unbans a member from a chat.

### `GetMembers`

```go
func (a *Account) GetMembers(ctx context.Context, objectGUID string) ([]map[string]any, error)
```

GetMembers returns the full member list of a chat (group/channel).
