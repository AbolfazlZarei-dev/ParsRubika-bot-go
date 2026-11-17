# `voicekit` — Package voicekit provides high-level helpers for the Rubika

_Source: [`voicekit.go`](../../../ParsRubika/extras/voicekit/voicekit.go)_

## Symbols

### `State`

```go
type State string
```

State is the state of a voice-chat participant.

### `StateMute`

```go
const StateMute State = "Mute"
```

### `StateUnmute`

```go
const StateUnmute State = "Unmute"
```

### `Activity`

```go
type Activity string
```

Activity is the activity reported in a voice chat.

### `ActivitySpeaking`

```go
const ActivitySpeaking Activity = "Speaking"
```

### `Join`

```go
func Join(ctx context.Context, acc *account.Account, objectGUID, voiceChatID, sdpOffer string) (*events.Event, error)
```

Join joins a voice chat.

### `Leave`

```go
func Leave(ctx context.Context, acc *account.Account, objectGUID, voiceChatID string) (*events.Event, error)
```

Leave leaves a voice chat.

### `MuteParticipant`

```go
func MuteParticipant(ctx context.Context, acc *account.Account, objectGUID, voiceChatID, participantGUID string) (*events.Event, error)
```

MuteParticipant mutes a participant in a voice chat.

### `UnmuteParticipant`

```go
func UnmuteParticipant(ctx context.Context, acc *account.Account, objectGUID, voiceChatID, participantGUID string) (*events.Event, error)
```

UnmuteParticipant unmutes a participant in a voice chat.

### `ReportSpeaking`

```go
func ReportSpeaking(ctx context.Context, acc *account.Account, objectGUID, voiceChatID, participantGUID string) (*events.Event, error)
```

ReportSpeaking sends a speaking activity notification.

### `Create`

```go
func Create(ctx context.Context, acc *account.Account, groupGUID, title string) (*events.Event, error)
```

Create creates a voice chat in a group.

### `Close`

```go
func Close(ctx context.Context, acc *account.Account, objectGUID, voiceChatID string) (*events.Event, error)
```

Close closes a voice chat.

### `GetVoiceChatID`

```go
func GetVoiceChatID(chatInfo *events.Event) string
```

GetVoiceChatID extracts the voice_chat_id from a chat-info
event. This is a convenience helper since the ID is nested
deeply in the response.

### `FormatError`

```go
func FormatError(err error) string
```

FormatError returns a human-readable description of a voice-chat
error.
