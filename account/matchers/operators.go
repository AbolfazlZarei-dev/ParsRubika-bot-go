package matchers

import (
	"context"
	"strings"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// OperatorKind enumerates the relational operators supported by
// Operator-based matchers.
type OperatorKind int

const (
	OpEqual        OperatorKind = iota // ==
	OpNotEqual                         // !=
	OpLess                             // <
	OpLessEqual                        // <=
	OpGreater                          // >
	OpGreaterEqual                     // >=
)

// Operator is a small helper that pairs a value with an operator and
// can be used by any numeric/string matcher. It is NOT a matcher on
// its own; matchers embed it and call `Operator.Match(actual)`.
type Operator struct {
	Kind  OperatorKind
	Value any
}

// Match returns true when the operator's comparison holds for `actual`.
func (o Operator) Match(actual any) bool {
	switch o.Kind {
	case OpEqual:
		return cmpAny(actual, o.Value) == 0
	case OpNotEqual:
		return cmpAny(actual, o.Value) != 0
	case OpLess:
		return cmpAny(actual, o.Value) < 0
	case OpLessEqual:
		return cmpAny(actual, o.Value) <= 0
	case OpGreater:
		return cmpAny(actual, o.Value) > 0
	case OpGreaterEqual:
		return cmpAny(actual, o.Value) >= 0
	}
	return false
}

// cmpAny compares two values. Numbers are compared numerically,
// strings lexicographically. Mixed types are compared by their string
// representation.
func cmpAny(a, b any) int {
	switch av := a.(type) {
	case int:
		if bv, ok := b.(int); ok {
			return compareInt(av, bv)
		}
		if bv, ok := b.(int64); ok {
			return compareInt64(int64(av), bv)
		}
		if bv, ok := b.(float64); ok {
			return compareFloat(float64(av), bv)
		}
	case int64:
		if bv, ok := b.(int64); ok {
			return compareInt64(av, bv)
		}
		if bv, ok := b.(int); ok {
			return compareInt64(av, int64(bv))
		}
		if bv, ok := b.(float64); ok {
			return compareFloat(float64(av), bv)
		}
	case float64:
		if bv, ok := b.(float64); ok {
			return compareFloat(av, bv)
		}
	case string:
		if bv, ok := b.(string); ok {
			return strings.Compare(av, bv)
		}
	}
	return strings.Compare(toStr(a), toStr(b))
}

func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func compareInt64(a, b int64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func compareFloat(a, b float64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func toStr(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case int:
		return intToStr(int64(x))
	case int64:
		return intToStr(x)
	case float64:
		return floatToStr(x)
	case bool:
		if x {
			return "true"
		}
		return "false"
	}
	return ""
}

func intToStr(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func floatToStr(f float64) string {
	// Simple %g-style formatting without importing fmt.
	if f == float64(int64(f)) {
		return intToStr(int64(f))
	}
	// Fallback: rough 6-decimal precision.
	neg := f < 0
	if neg {
		f = -f
	}
	intPart := int64(f)
	frac := f - float64(intPart)
	// Multiply by 1e6, round to int.
	fracInt := int64(frac*1e6 + 0.5)
	// Build "intPart.fracInt" with leading zeros on frac.
	s := intToStr(intPart) + "." + padLeft(intToStr(fracInt), 6)
	if neg {
		return "-" + s
	}
	return s
}

func padLeft(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return strings.Repeat("0", n-len(s)) + s
}

// OpMatcher is a generic matcher that pairs an accessor with an
// Operator. It is the building block for the typed matchers below
// (IsPinned, IsMute, MessageID, etc.).
type OpMatcher struct {
	kind     string
	accessor func(*events.Event) any
	op       Operator
}

// Kind implements Matcher.
func (m *OpMatcher) Kind() string { return m.kind }

// Match implements Matcher.
func (m *OpMatcher) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	if m.accessor == nil {
		return false
	}
	actual := m.accessor(ev)
	if m.op.Kind == OpEqual && m.op.Value == nil {
		return actual != nil
	}
	return m.op.Match(actual)
}

// NewMessageID creates a matcher that compares the update's message_id
// against `value` using `op`.
func NewMessageID(op OperatorKind, value any) *OpMatcher {
	return &OpMatcher{
		kind:     "MessageUpdates",
		accessor: func(e *events.Event) any { return e.MessageID() },
		op:       Operator{Kind: op, Value: value},
	}
}

// NewReplyMessageID creates a matcher that compares the update's
// reply_to_message_id against `value`.
func NewReplyMessageID(op OperatorKind, value any) *OpMatcher {
	return &OpMatcher{
		kind:     "MessageUpdates",
		accessor: func(e *events.Event) any { return e.ReplyMessageID() },
		op:       Operator{Kind: op, Value: value},
	}
}

// NewTime creates a matcher that compares the update's time against
// `value` using `op`.
func NewTime(op OperatorKind, value any) *OpMatcher {
	return &OpMatcher{
		kind:     "MessageUpdates",
		accessor: func(e *events.Event) any { return e.Str("time") },
		op:       Operator{Kind: op, Value: value},
	}
}

// ===== Boolean predicates ==============================================

// IsPinned matches updates whose message has the `pinned` flag set.
type IsPinned struct{ kind string }

// NewIsPinned constructs an IsPinned matcher.
func NewIsPinned(kind string) *IsPinned {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsPinned{kind: kind}
}

// Kind implements Matcher.
func (m *IsPinned) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsPinned) Match(_ context.Context, raw map[string]any) bool {
	return asEvent(raw).Find("pinned") != nil
}

// IsMute matches updates whose chat has the `mute` flag set.
type IsMute struct{ kind string }

// NewIsMute constructs an IsMute matcher.
func NewIsMute(kind string) *IsMute {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsMute{kind: kind}
}

// Kind implements Matcher.
func (m *IsMute) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsMute) Match(_ context.Context, raw map[string]any) bool {
	return asEvent(raw).Find("mute") != nil
}

// IsInContact matches updates whose author is in the user's address
// book. The client cannot tell from the raw update alone, so this
// matcher only fires when the server explicitly sets the
// `is_in_contact` flag in the message.
type IsInContact struct{ kind string }

// NewIsInContact constructs an IsInContact matcher.
func NewIsInContact(kind string) *IsInContact {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsInContact{kind: kind}
}

// Kind implements Matcher.
func (m *IsInContact) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsInContact) Match(_ context.Context, raw map[string]any) bool {
	return asEvent(raw).Find("is_in_contact") != nil
}

// ===== Media type matchers ============================================

// IsVideo matches messages with a file_inline of type "Video".
type IsVideo struct{ kind string }

// NewIsVideo constructs an IsVideo matcher.
func NewIsVideo(kind string) *IsVideo {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsVideo{kind: kind}
}

// Kind implements Matcher.
func (m *IsVideo) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsVideo) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	return ev.IsFileInline() && ev.FileInline() != nil && ev.FileInline().Type() == "Video"
}

// IsVoice matches messages with a file_inline of type "Voice".
type IsVoice struct{ kind string }

// NewIsVoice constructs an IsVoice matcher.
func NewIsVoice(kind string) *IsVoice {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsVoice{kind: kind}
}

// Kind implements Matcher.
func (m *IsVoice) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsVoice) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	return ev.IsFileInline() && ev.FileInline() != nil && ev.FileInline().Type() == "Voice"
}

// IsMusic matches messages with a file_inline of type "Music".
type IsMusic struct{ kind string }

// NewIsMusic constructs an IsMusic matcher.
func NewIsMusic(kind string) *IsMusic {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsMusic{kind: kind}
}

// Kind implements Matcher.
func (m *IsMusic) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsMusic) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	return ev.IsFileInline() && ev.FileInline() != nil && ev.FileInline().Type() == "Music"
}

// IsGif matches messages with a file_inline of type "Gif".
type IsGif struct{ kind string }

// NewIsGif constructs an IsGif matcher.
func NewIsGif(kind string) *IsGif {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsGif{kind: kind}
}

// Kind implements Matcher.
func (m *IsGif) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsGif) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	return ev.IsFileInline() && ev.FileInline() != nil && ev.FileInline().Type() == "Gif"
}

// IsContact matches messages with a contact attachment.
type IsContact struct{ kind string }

// NewIsContact constructs an IsContact matcher.
func NewIsContact(kind string) *IsContact {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsContact{kind: kind}
}

// Kind implements Matcher.
func (m *IsContact) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsContact) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	return ev.IsFileInline() && ev.FileInline() != nil && ev.FileInline().Type() == "Contact"
}

// IsLocation matches messages with a location attachment.
type IsLocation struct{ kind string }

// NewIsLocation constructs an IsLocation matcher.
func NewIsLocation(kind string) *IsLocation {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsLocation{kind: kind}
}

// Kind implements Matcher.
func (m *IsLocation) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsLocation) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	return ev.IsFileInline() && ev.FileInline() != nil && ev.FileInline().Type() == "Location"
}

// IsPoll matches messages with a poll attachment.
type IsPoll struct{ kind string }

// NewIsPoll constructs an IsPoll matcher.
func NewIsPoll(kind string) *IsPoll {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsPoll{kind: kind}
}

// Kind implements Matcher.
func (m *IsPoll) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsPoll) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	return ev.IsFileInline() && ev.FileInline() != nil && ev.FileInline().Type() == "Poll"
}

// IsEvent matches service/event messages (join/leave/pin/title etc.).
type IsEvent struct{ kind string }

// NewIsEvent constructs an IsEvent matcher.
func NewIsEvent(kind string) *IsEvent {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsEvent{kind: kind}
}

// Kind implements Matcher.
func (m *IsEvent) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsEvent) Match(_ context.Context, raw map[string]any) bool {
	return asEvent(raw).IsEvent()
}

// EventType matches event messages of a specific type (e.g.
// "AddedGroupMembers", "PinnedMessageUpdated").
type EventType struct {
	eventType string
	kind      string
}

// NewEventType constructs an EventType matcher.
func NewEventType(eventType, kind string) *EventType {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &EventType{eventType: eventType, kind: kind}
}

// Kind implements Matcher.
func (m *EventType) Kind() string { return m.kind }

// Match implements Matcher.
func (m *EventType) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	if !ev.IsEvent() {
		return false
	}
	msg := ev.Message()
	if msg == nil {
		return false
	}
	ed := msg.EventOf("event_data")
	if ed == nil {
		return false
	}
	return ed.Str("type") == m.eventType
}

// IsForwardedFromUser matches messages forwarded from a user.
type IsForwardedFromUser struct{ kind string }

// NewIsForwardedFromUser constructs an IsForwardedFromUser matcher.
func NewIsForwardedFromUser(kind string) *IsForwardedFromUser {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsForwardedFromUser{kind: kind}
}

// Kind implements Matcher.
func (m *IsForwardedFromUser) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsForwardedFromUser) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	if !ev.IsForward() {
		return false
	}
	msg := ev.Message()
	if msg == nil {
		return false
	}
	ff := msg.EventOf("forwarded_from")
	if ff == nil {
		return false
	}
	return ff.Str("type_from") == "User"
}

// IsForwardedFromChannel matches messages forwarded from a channel.
type IsForwardedFromChannel struct{ kind string }

// NewIsForwardedFromChannel constructs an IsForwardedFromChannel matcher.
func NewIsForwardedFromChannel(kind string) *IsForwardedFromChannel {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsForwardedFromChannel{kind: kind}
}

// Kind implements Matcher.
func (m *IsForwardedFromChannel) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsForwardedFromChannel) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	if !ev.IsForward() {
		return false
	}
	msg := ev.Message()
	if msg == nil {
		return false
	}
	ff := msg.EventOf("forwarded_from")
	if ff == nil {
		return false
	}
	return ff.Str("type_from") == "Channel"
}

// HasReaction matches messages that have at least one reaction.
type HasReaction struct{ kind string }

// NewHasReaction constructs a HasReaction matcher.
func NewHasReaction(kind string) *HasReaction {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &HasReaction{kind: kind}
}

// Kind implements Matcher.
func (m *HasReaction) Kind() string { return m.kind }

// Match implements Matcher.
func (m *HasReaction) Match(_ context.Context, raw map[string]any) bool {
	return asEvent(raw).Find("reactions") != nil
}
