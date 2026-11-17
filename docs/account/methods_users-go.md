# `account` — Methods users

_Source: [`methods_users.go`](../../ParsRubika/account/methods_users.go)_

## Symbols

### `GetObjectByUsername`

```go
func (a *Account) GetObjectByUsername(ctx context.Context, username string) (*events.Event, error)
```

GetObjectByUsername resolves a username to an object GUID. The
leading `@` is stripped automatically.

### `UserIsAdmin`

```go
func (a *Account) UserIsAdmin(ctx context.Context, objectGUID, userGUID string) (bool, error)
```

UserIsAdmin reports whether `userGUID` is an admin of `objectGUID`.
The API does not have a dedicated `userIsAdmin` method, so we paginate
`getGroupAdminMembers` / `getChannelAdminMembers` and check membership
client-side.

### `GetUserInfo`

```go
func (a *Account) GetUserInfo(ctx context.Context, userGUID string) (*events.Event, error)
```

GetUserInfo returns information about a user. When `userGUID` is
empty the server returns information about the current user.

### `GetMe`

```go
func (a *Account) GetMe(ctx context.Context) (*events.Event, error)
```

GetMe is a convenience wrapper for GetUserInfo that fetches the
currently authenticated user.

### `SetBlockUser`

```go
func (a *Account) SetBlockUser(ctx context.Context, userGUID, action string) (*events.Event, error)
```

SetBlockUser toggles the blocked state of `userGUID`. `action` is
"Block" or "Unblock".

### `BlockUser`

```go
func (a *Account) BlockUser(ctx context.Context, userGUID string) (*events.Event, error)
```

BlockUser blocks the supplied user.

### `UnblockUser`

```go
func (a *Account) UnblockUser(ctx context.Context, userGUID string) (*events.Event, error)
```

UnblockUser unblocks the supplied user.

### `DeleteUserChat`

```go
func (a *Account) DeleteUserChat(ctx context.Context, userGUID, lastMessageID string) (*events.Event, error)
```

DeleteUserChat deletes the chat history with the supplied user GUID.

### `CheckUserUsername`

```go
func (a *Account) CheckUserUsername(ctx context.Context, username string) (*events.Event, error)
```

CheckUserUsername reports whether the supplied username exists.

### `GetCommonGroups`

```go
func (a *Account) GetCommonGroups(ctx context.Context, userGUID string) (*events.Event, error)
```

GetCommonGroups returns the chats shared with the supplied user.

### `GetLastOnline`

```go
func (a *Account) GetLastOnline(ctx context.Context, userGUID string) (*events.Event, error)
```

GetLastOnline returns the online status of the supplied user. This
is a composite helper that fetches the user info and reads the
`online_time` field. It guards against nil intermediate events so a
missing "user" key does not panic.

### `RequestVoiceCall`

```go
func (a *Account) RequestVoiceCall(ctx context.Context, userGUID string) (*events.Event, error)
```

RequestVoiceCall initiates a voice call with the supplied user.

### `RequestVideoCall`

```go
func (a *Account) RequestVideoCall(ctx context.Context, userGUID string) (*events.Event, error)
```

RequestVideoCall initiates a video call with the supplied user.

### `DiscardCall`

```go
func (a *Account) DiscardCall(ctx context.Context, callID string, duration int, reason string) (*events.Event, error)
```

DiscardCall ends an active call. `reason` is "Missed" or "Disconnect".

### `GetBotInfo`

```go
func (a *Account) GetBotInfo(ctx context.Context, botGUID string) (*events.Event, error)
```

GetBotInfo returns info about a bot (GUID must start with "b0").

### `StartBot`

```go
func (a *Account) StartBot(ctx context.Context, botGUID string) (*events.Event, error)
```

StartBot sends a "/start" message to a bot.

### `StopBot`

```go
func (a *Account) StopBot(ctx context.Context, botGUID string) (*events.Event, error)
```

StopBot stops a bot session.

### `DeleteBotChat`

```go
func (a *Account) DeleteBotChat(ctx context.Context, botGUID, lastMessageID string) (*events.Event, error)
```

DeleteBotChat deletes the chat history with a bot.

### `GetInlineBots`

```go
func (a *Account) GetInlineBots(ctx context.Context, botGUIDs []string) (*events.Event, error)
```

GetInlineBots returns info about multiple inline bots.

### `GetServiceInfo`

```go
func (a *Account) GetServiceInfo(ctx context.Context, serviceGUID string) (*events.Event, error)
```

GetServiceInfo returns info about a service chat (GUID must start
with "s0").

### `DeleteServiceChat`

```go
func (a *Account) DeleteServiceChat(ctx context.Context, serviceGUID, lastMessageID string) (*events.Event, error)
```

DeleteServiceChat deletes the chat history with a service chat.
