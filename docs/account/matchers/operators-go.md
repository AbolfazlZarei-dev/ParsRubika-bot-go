# `matchers` — Operators

_Source: [`operators.go`](../../../ParsRubika/account/matchers/operators.go)_

## Symbols

### `OperatorKind`

```go
type OperatorKind int
```

OperatorKind enumerates the relational operators supported by
Operator-based matchers.

### `OpEqual`

```go
const OpEqual OperatorKind = iota // ==
```

### `OpNotEqual`

```go
const OpNotEqual // !=
```

### `OpLess`

```go
const OpLess // <
```

### `OpLessEqual`

```go
const OpLessEqual // <=
```

### `OpGreater`

```go
const OpGreater // >
```

### `OpGreaterEqual`

```go
const OpGreaterEqual // >=
```

### `Operator`

```go
type Operator struct
```

Operator is a small helper that pairs a value with an operator and
can be used by any numeric/string matcher. It is NOT a matcher on
its own; matchers embed it and call `Operator.Match(actual)`.

### `Match`

```go
func (o Operator) Match(actual any) bool
```

Match returns true when the operator's comparison holds for `actual`.

### `OpMatcher`

```go
type OpMatcher struct
```

OpMatcher is a generic matcher that pairs an accessor with an
Operator. It is the building block for the typed matchers below
(IsPinned, IsMute, MessageID, etc.).

### `Kind`

```go
func (m *OpMatcher) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *OpMatcher) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `NewMessageID`

```go
func NewMessageID(op OperatorKind, value any) *OpMatcher
```

NewMessageID creates a matcher that compares the update's message_id
against `value` using `op`.

### `NewReplyMessageID`

```go
func NewReplyMessageID(op OperatorKind, value any) *OpMatcher
```

NewReplyMessageID creates a matcher that compares the update's
reply_to_message_id against `value`.

### `NewTime`

```go
func NewTime(op OperatorKind, value any) *OpMatcher
```

NewTime creates a matcher that compares the update's time against
`value` using `op`.

### `IsPinned`

```go
type IsPinned struct
```

IsPinned matches updates whose message has the `pinned` flag set.

### `NewIsPinned`

```go
func NewIsPinned(kind string) *IsPinned
```

NewIsPinned constructs an IsPinned matcher.

### `Kind`

```go
func (m *IsPinned) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsPinned) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsMute`

```go
type IsMute struct
```

IsMute matches updates whose chat has the `mute` flag set.

### `NewIsMute`

```go
func NewIsMute(kind string) *IsMute
```

NewIsMute constructs an IsMute matcher.

### `Kind`

```go
func (m *IsMute) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsMute) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsInContact`

```go
type IsInContact struct
```

IsInContact matches updates whose author is in the user's address
book. The client cannot tell from the raw update alone, so this
matcher only fires when the server explicitly sets the
`is_in_contact` flag in the message.

### `NewIsInContact`

```go
func NewIsInContact(kind string) *IsInContact
```

NewIsInContact constructs an IsInContact matcher.

### `Kind`

```go
func (m *IsInContact) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsInContact) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsVideo`

```go
type IsVideo struct
```

IsVideo matches messages with a file_inline of type "Video".

### `NewIsVideo`

```go
func NewIsVideo(kind string) *IsVideo
```

NewIsVideo constructs an IsVideo matcher.

### `Kind`

```go
func (m *IsVideo) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsVideo) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsVoice`

```go
type IsVoice struct
```

IsVoice matches messages with a file_inline of type "Voice".

### `NewIsVoice`

```go
func NewIsVoice(kind string) *IsVoice
```

NewIsVoice constructs an IsVoice matcher.

### `Kind`

```go
func (m *IsVoice) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsVoice) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsMusic`

```go
type IsMusic struct
```

IsMusic matches messages with a file_inline of type "Music".

### `NewIsMusic`

```go
func NewIsMusic(kind string) *IsMusic
```

NewIsMusic constructs an IsMusic matcher.

### `Kind`

```go
func (m *IsMusic) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsMusic) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsGif`

```go
type IsGif struct
```

IsGif matches messages with a file_inline of type "Gif".

### `NewIsGif`

```go
func NewIsGif(kind string) *IsGif
```

NewIsGif constructs an IsGif matcher.

### `Kind`

```go
func (m *IsGif) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsGif) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsContact`

```go
type IsContact struct
```

IsContact matches messages with a contact attachment.

### `NewIsContact`

```go
func NewIsContact(kind string) *IsContact
```

NewIsContact constructs an IsContact matcher.

### `Kind`

```go
func (m *IsContact) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsContact) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsLocation`

```go
type IsLocation struct
```

IsLocation matches messages with a location attachment.

### `NewIsLocation`

```go
func NewIsLocation(kind string) *IsLocation
```

NewIsLocation constructs an IsLocation matcher.

### `Kind`

```go
func (m *IsLocation) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsLocation) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsPoll`

```go
type IsPoll struct
```

IsPoll matches messages with a poll attachment.

### `NewIsPoll`

```go
func NewIsPoll(kind string) *IsPoll
```

NewIsPoll constructs an IsPoll matcher.

### `Kind`

```go
func (m *IsPoll) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsPoll) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsEvent`

```go
type IsEvent struct
```

IsEvent matches service/event messages (join/leave/pin/title etc.).

### `NewIsEvent`

```go
func NewIsEvent(kind string) *IsEvent
```

NewIsEvent constructs an IsEvent matcher.

### `Kind`

```go
func (m *IsEvent) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsEvent) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `EventType`

```go
type EventType struct
```

EventType matches event messages of a specific type (e.g.
"AddedGroupMembers", "PinnedMessageUpdated").

### `NewEventType`

```go
func NewEventType(eventType, kind string) *EventType
```

NewEventType constructs an EventType matcher.

### `Kind`

```go
func (m *EventType) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *EventType) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsForwardedFromUser`

```go
type IsForwardedFromUser struct
```

IsForwardedFromUser matches messages forwarded from a user.

### `NewIsForwardedFromUser`

```go
func NewIsForwardedFromUser(kind string) *IsForwardedFromUser
```

NewIsForwardedFromUser constructs an IsForwardedFromUser matcher.

### `Kind`

```go
func (m *IsForwardedFromUser) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsForwardedFromUser) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsForwardedFromChannel`

```go
type IsForwardedFromChannel struct
```

IsForwardedFromChannel matches messages forwarded from a channel.

### `NewIsForwardedFromChannel`

```go
func NewIsForwardedFromChannel(kind string) *IsForwardedFromChannel
```

NewIsForwardedFromChannel constructs an IsForwardedFromChannel matcher.

### `Kind`

```go
func (m *IsForwardedFromChannel) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsForwardedFromChannel) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `HasReaction`

```go
type HasReaction struct
```

HasReaction matches messages that have at least one reaction.

### `NewHasReaction`

```go
func NewHasReaction(kind string) *HasReaction
```

NewHasReaction constructs a HasReaction matcher.

### `Kind`

```go
func (m *HasReaction) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *HasReaction) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.
