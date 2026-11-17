# `listeners` — Package listeners provides the five built-in listener matchers used

_Source: [`listeners.go`](../../../ParsRubika/account/listeners/listeners.go)_

## Symbols

### `Messages`

```go
type Messages struct
```

Messages is the matcher for "MessageUpdates" updates (server key
`message_updates`). When no sub-matchers are supplied it matches
every message update.

### `Chats`

```go
type Chats struct
```

Chats matches "ChatUpdates" updates (server key `chat_updates`).

### `Activities`

```go
type Activities struct
```

Activities matches "ShowActivities" updates.

### `Notifications`

```go
type Notifications struct
```

Notifications matches "ShowNotifications" updates.

### `Removals`

```go
type Removals struct
```

Removals matches "RemoveNotifications" updates.

### `ScheduledMessages`

```go
type ScheduledMessages struct
```

ScheduledMessages matches "ScheduleMessageUpdates" updates.

### `DraftMessages`

```go
type DraftMessages struct
```

DraftMessages matches "DraftMessageUpdate" updates.

### `UnconfirmedSessions`

```go
type UnconfirmedSessions struct
```

UnconfirmedSessions matches "UnconfirmedSessionUpdates" updates.

### `VoiceChatParticipants`

```go
type VoiceChatParticipants struct
```

VoiceChatParticipants matches "GroupVoiceChatParticipantUpdates" updates.

### `VoiceChatUpdates`

```go
type VoiceChatUpdates struct
```

VoiceChatUpdates matches "GroupVoiceChatUpdates" updates.

### `CallUpdates`

```go
type CallUpdates struct
```

CallUpdates matches "CallUpdates" updates.

### `CallSignals`

```go
type CallSignals struct
```

CallSignals matches "CallSignalData" updates.

### `MatcherChain`

```go
type MatcherChain interface
```

MatcherChain is the common interface implemented by every listener.
It is identical to the account.Matcher interface; it is re-declared
here to avoid an import cycle.

### `Kind`

```go
func (Messages) Kind() string { return "MessageUpdates" }
```

### `Kind`

```go
func (Chats) Kind() string { return "ChatUpdates" }
```

### `Kind`

```go
func (Activities) Kind() string { return "ShowActivities" }
```

### `Kind`

```go
func (Notifications) Kind() string { return "ShowNotifications" }
```

### `Kind`

```go
func (Removals) Kind() string { return "RemoveNotifications" }
```

### `Kind`

```go
func (ScheduledMessages) Kind() string { return "ScheduleMessageUpdates" }
```

### `Kind`

```go
func (DraftMessages) Kind() string { return "DraftMessageUpdate" }
```

### `Kind`

```go
func (UnconfirmedSessions) Kind() string { return "UnconfirmedSessionUpdates" }
```

### `Kind`

```go
func (VoiceChatParticipants) Kind() string { return "GroupVoiceChatParticipantUpdates" }
```

### `Kind`

```go
func (VoiceChatUpdates) Kind() string { return "GroupVoiceChatUpdates" }
```

### `Kind`

```go
func (CallUpdates) Kind() string { return "CallUpdates" }
```

### `Kind`

```go
func (CallSignals) Kind() string { return "CallSignalData" }
```

### `Match`

```go
func (m Messages) Match(ctx context.Context, raw map[string]any) bool
```

### `Match`

```go
func (c Chats) Match(ctx context.Context, raw map[string]any) bool
```

### `Match`

```go
func (a Activities) Match(ctx context.Context, raw map[string]any) bool
```

### `Match`

```go
func (n Notifications) Match(ctx context.Context, raw map[string]any) bool
```

### `Match`

```go
func (r Removals) Match(ctx context.Context, raw map[string]any) bool
```

### `Match`

```go
func (s ScheduledMessages) Match(ctx context.Context, raw map[string]any) bool
```

### `Match`

```go
func (d DraftMessages) Match(ctx context.Context, raw map[string]any) bool
```

### `Match`

```go
func (u UnconfirmedSessions) Match(ctx context.Context, raw map[string]any) bool
```

### `Match`

```go
func (v VoiceChatParticipants) Match(ctx context.Context, raw map[string]any) bool
```

### `Match`

```go
func (v VoiceChatUpdates) Match(ctx context.Context, raw map[string]any) bool
```

### `Match`

```go
func (c CallUpdates) Match(ctx context.Context, raw map[string]any) bool
```

### `Match`

```go
func (c CallSignals) Match(ctx context.Context, raw map[string]any) bool
```

### `NewMessages`

```go
func NewMessages(inner MatcherChain) Messages { return Messages{inner: inner} }
```

NewMessages constructs a Messages listener that wraps the supplied
matcher chain.

### `NewChats`

```go
func NewChats(inner MatcherChain) Chats { return Chats{inner: inner} }
```

NewChats constructs a Chats listener.

### `NewActivities`

```go
func NewActivities(inner MatcherChain) Activities { return Activities{inner: inner} }
```

NewActivities constructs an Activities listener.

### `NewNotifications`

```go
func NewNotifications(inner MatcherChain) Notifications { return Notifications{inner: inner} }
```

NewNotifications constructs a Notifications listener.

### `NewRemovals`

```go
func NewRemovals(inner MatcherChain) Removals { return Removals{inner: inner} }
```

NewRemovals constructs a Removals listener.

### `NewScheduledMessages`

```go
func NewScheduledMessages(inner MatcherChain) ScheduledMessages
```

NewScheduledMessages constructs a ScheduledMessages listener.

### `NewDraftMessages`

```go
func NewDraftMessages(inner MatcherChain) DraftMessages { return DraftMessages{inner: inner} }
```

NewDraftMessages constructs a DraftMessages listener.

### `NewUnconfirmedSessions`

```go
func NewUnconfirmedSessions(inner MatcherChain) UnconfirmedSessions
```

NewUnconfirmedSessions constructs an UnconfirmedSessions listener.

### `NewVoiceChatParticipants`

```go
func NewVoiceChatParticipants(inner MatcherChain) VoiceChatParticipants
```

NewVoiceChatParticipants constructs a VoiceChatParticipants listener.

### `NewVoiceChatUpdates`

```go
func NewVoiceChatUpdates(inner MatcherChain) VoiceChatUpdates { return VoiceChatUpdates{inner: inner} }
```

NewVoiceChatUpdates constructs a VoiceChatUpdates listener.

### `NewCallUpdates`

```go
func NewCallUpdates(inner MatcherChain) CallUpdates { return CallUpdates{inner: inner} }
```

NewCallUpdates constructs a CallUpdates listener.

### `NewCallSignals`

```go
func NewCallSignals(inner MatcherChain) CallSignals { return CallSignals{inner: inner} }
```

NewCallSignals constructs a CallSignals listener.
