# `account` — Methods voicechat

_Source: [`methods_voicechat.go`](../../ParsRubika/account/methods_voicechat.go)_

## Symbols

### `CreateChannelVoiceChat`

```go
func (a *Account) CreateChannelVoiceChat(ctx context.Context, channelGUID, title string) (*events.Event, error)
```

CreateChannelVoiceChat creates a voice chat inside a channel.

### `DiscardChannelVoiceChat`

```go
func (a *Account) DiscardChannelVoiceChat(ctx context.Context, channelGUID, voiceChatID string) (*events.Event, error)
```

DiscardChannelVoiceChat closes a channel voice chat.

### `SetChannelVoiceChatSetting`

```go
func (a *Account) SetChannelVoiceChatSetting(ctx context.Context, channelGUID, voiceChatID string, settings map[string]any) (*events.Event, error)
```

SetChannelVoiceChatSetting updates voice-chat settings. Only the
non-empty fields are sent; `updated_parameters` is built automatically.

### `GetGroupVoiceChat`

```go
func (a *Account) GetGroupVoiceChat(ctx context.Context, groupGUID, voiceChatID string) (*events.Event, error)
```

GetGroupVoiceChat returns info about a group voice chat.

### `GetGroupVoiceChatParticipants`

```go
func (a *Account) GetGroupVoiceChatParticipants(ctx context.Context, groupGUID, voiceChatID, startID string) (*events.Event, error)
```

GetGroupVoiceChatParticipants returns the participants of a voice chat.

### `GetChatVoiceChat`

```go
func (a *Account) GetChatVoiceChat(ctx context.Context, objectGUID string) (*events.Event, error)
```

GetChatVoiceChat returns the voice chat ID for a chat. Only groups
and channels support voice chats; user/bot/service GUIDs return an
explicit error rather than being mis-dispatched to GetGroupInfo.

### `CreateGroupVoiceChat`

```go
func (a *Account) CreateGroupVoiceChat(ctx context.Context, groupGUID, title string) (*events.Event, error)
```

CreateGroupVoiceChat creates a voice chat inside a group. The
server expects the uniqueness token `rnd` and the protocol
`version` alongside the GUID; `title` is optional.

### `LeaveGroupVoiceChat`

```go
func (a *Account) LeaveGroupVoiceChat(ctx context.Context, groupGUID, voiceChatID string) (*events.Event, error)
```

LeaveGroupVoiceChat leaves a group voice chat.

### `GetGroupVoiceChatUpdates`

```go
func (a *Account) GetGroupVoiceChatUpdates(ctx context.Context, groupGUID, voiceChatID string) (*events.Event, error)
```

GetGroupVoiceChatUpdates fetches the latest voice-chat updates.

### `SetGroupVoiceChatSetting`

```go
func (a *Account) SetGroupVoiceChatSetting(ctx context.Context, groupGUID, voiceChatID string, settings map[string]any) (*events.Event, error)
```

SetGroupVoiceChatSetting updates voice-chat settings. Only the
non-empty fields in `settings` are sent; the `updated_parameters`
array is built automatically.

### `StartVoiceChat`

```go
func (a *Account) StartVoiceChat(ctx context.Context, groupGUID, title string) (*events.Event, error)
```

StartVoiceChat creates a voice chat in a group (alias for
CreateGroupVoiceChat with a more intuitive name).

### `RequestSpeakVoiceChat`

```go
func (a *Account) RequestSpeakVoiceChat(ctx context.Context, objectGUID, voiceChatID, participantGUID string) (*events.Event, error)
```

RequestSpeakVoiceChat requests to speak in a voice chat.

### `CancelRequestSpeakVoiceChat`

```go
func (a *Account) CancelRequestSpeakVoiceChat(ctx context.Context, objectGUID, voiceChatID, participantGUID string) (*events.Event, error)
```

CancelRequestSpeakVoiceChat cancels a pending speak request.

### `JoinVoiceChat`

```go
func (a *Account) JoinVoiceChat(ctx context.Context, objectGUID, voiceChatID, sdpOfferData string) (*events.Event, error)
```

JoinVoiceChat joins a voice chat. The API method is dispatched to
`joinGroupVoiceChat` or `joinChannelVoiceChat` based on the GUID
prefix. `sdpOfferData` is the WebRTC SDP offer (may be empty for a
simple presence-only join).

### `DiscardVoiceChat`

```go
func (a *Account) DiscardVoiceChat(ctx context.Context, objectGUID, voiceChatID string) (*events.Event, error)
```

DiscardVoiceChat closes a voice chat (dispatches to Group/Channel).

### `SetVoiceChatState`

```go
func (a *Account) SetVoiceChatState(ctx context.Context, objectGUID, voiceChatID, action, participantGUID string) (*events.Event, error)
```

SetVoiceChatState updates the voice-chat state. `action` is "Mute"
or "Unmute". The API method is dispatched to
`setGroupVoiceChatState` / `setChannelVoiceChatState`.

### `SendGroupVoiceChatActivity`

```go
func (a *Account) SendGroupVoiceChatActivity(ctx context.Context, objectGUID, voiceChatID, action, participantGUID string) (*events.Event, error)
```

SendGroupVoiceChatActivity sends an activity notification in a voice
chat. `participantGUID` defaults to the current user's GUID.
