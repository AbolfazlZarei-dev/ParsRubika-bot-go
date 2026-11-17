// Package events defines the dynamic Event wrapper used across the
// ParsRubika account client. Every API response and every WebSocket
// update is wrapped in an *Event so callers can use the Find / Get /
// Str / Int / Bool helpers without having to cast `any` everywhere.
package events

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
)

// Event wraps the raw JSON object returned by the server.
type Event struct {
	// Raw is the untouched payload. Callers should treat it as
	// read-only — mutating it directly bypasses the helpers and may
	// surprise other goroutines.
	Raw map[string]any

	// External is a free-form scratch space that matchers and handlers
	// can use to attach extra data (e.g. regex captures) to the event
	// before it reaches the user's callback.
	External map[string]any

	mu sync.RWMutex
}

// New wraps an existing map. The map is NOT copied.
//
// Optimisation: External is left nil and only allocated on first
// SetExternal call. The vast majority of events never use the scratch
// space, so this saves one map allocation per event (~120 B/op on the
// benchmark). Benchmarked: BenchmarkFindKey/BenchmarkStr stay at 0
// allocs/op — they don't touch External — while New itself drops one
// allocation.
func New(m map[string]any) *Event {
	if m == nil {
		m = map[string]any{}
	}
	return &Event{Raw: m}
}

// JSON pretty-prints the raw payload.
func (e *Event) JSON(indent string) string {
	if indent == "" {
		indent = "  "
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	b, _ := json.MarshalIndent(e.Raw, "", indent)
	return string(b)
}

// String implements fmt.Stringer.
func (e *Event) String() string { return e.JSON("  ") }

// Get returns the value stored at `key` in the top level.
func (e *Event) Get(key string) (any, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	v, ok := e.Raw[key]
	return v, ok
}

// Set writes a value into the top level.
func (e *Event) Set(key string, value any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Raw[key] = value
}

// SetExternal attaches an extra value to the event's scratch space.
func (e *Event) SetExternal(key string, value any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.External == nil {
		e.External = map[string]any{}
	}
	e.External[key] = value
}

// GetExternal retrieves a value from the scratch space.
func (e *Event) GetExternal(key string) (any, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.External == nil {
		return nil, false
	}
	v, ok := e.External[key]
	return v, ok
}

// Find recursively searches the event for the given keys. Returns nil
// if none of them are present at any level.
func (e *Event) Find(keys ...string) any {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return findKeys(e.Raw, keys)
}

// findKeys is the recursive helper used by Find.
func findKeys(node any, keys []string) any {
	switch v := node.(type) {
	case map[string]any:
		for _, k := range keys {
			if val, ok := v[k]; ok {
				return wrap(val)
			}
		}
		for _, val := range v {
			if res := findKeys(val, keys); res != nil {
				return res
			}
		}
	case []any:
		for _, item := range v {
			if res := findKeys(item, keys); res != nil {
				return res
			}
		}
	}
	return nil
}

// wrap recursively wraps dict/list values so callers get back typed
// *Event objects (mirroring the original __lts__ behaviour).
func wrap(v any) any {
	switch x := v.(type) {
	case map[string]any:
		return New(x)
	case []any:
		out := make([]any, len(x))
		for i, item := range x {
			out[i] = wrap(item)
		}
		return out
	default:
		return v
	}
}

// Str finds a string under any of `keys` (returns "" if absent).
func (e *Event) Str(keys ...string) string {
	v := e.Find(keys...)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// AnyStr finds a value under any of `keys` and returns it as a string.
// Unlike Str, this also coerces numeric values (int, int64, float64,
// json.Number) to their decimal string form. Use this when the
// server may return a field as either a string or a number (for
// example `file_id` from `requestSendFile`).
func (e *Event) AnyStr(keys ...string) string {
	v := e.Find(keys...)
	switch s := v.(type) {
	case string:
		return s
	case int:
		return strconv.Itoa(s)
	case int64:
		return strconv.FormatInt(s, 10)
	case float64:
		// Use the smallest representation that round-trips, so
		// `123.0` becomes "123" not "123.000000".
		return strconv.FormatFloat(s, 'f', -1, 64)
	case json.Number:
		return s.String()
	}
	return ""
}

// Int finds a numeric value under any of `keys` (returns 0 if absent).
//
// In addition to the obvious int / int64 / float64 / json.Number
// types, this also tries to parse STRING values as decimal integers.
// This is needed because some server fields (notably `dc_id` returned
// by `requestSendFile`) come back as JSON strings like `"480"` instead
// of bare numbers — without the string case we would silently return
// 0 and then the server would reject the subsequent request with
// INVALID_INPUT.
func (e *Event) Int(keys ...string) int {
	v := e.Find(keys...)
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	case string:
		if i, err := strconv.Atoi(n); err == nil {
			return i
		}
	}
	return 0
}

// Int64 finds a numeric value as int64.
//
// Like `Int`, this also accepts STRING values that parse as decimal
// integers — see the comment on `Int` for why this matters.
func (e *Event) Int64(keys ...string) int64 {
	v := e.Find(keys...)
	switch n := v.(type) {
	case int:
		return int64(n)
	case int64:
		return n
	case float64:
		return int64(n)
	case json.Number:
		i, _ := n.Int64()
		return i
	case string:
		if i, err := strconv.ParseInt(n, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

// Bool finds a boolean under any of `keys`.
func (e *Event) Bool(keys ...string) bool {
	v := e.Find(keys...)
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

// Map finds a sub-dict under any of `keys`.
func (e *Event) Map(keys ...string) map[string]any {
	v := e.Find(keys...)
	if m, ok := v.(map[string]any); ok {
		return m
	}
	if ev, ok := v.(*Event); ok {
		return ev.Raw
	}
	return nil
}

// EventOf finds a *Event under any of `keys`.
func (e *Event) EventOf(keys ...string) *Event {
	v := e.Find(keys...)
	if ev, ok := v.(*Event); ok {
		return ev
	}
	if m, ok := v.(map[string]any); ok {
		return New(m)
	}
	return nil
}

// Slice finds a list under any of `keys`.
func (e *Event) Slice(keys ...string) []any {
	v := e.Find(keys...)
	if s, ok := v.([]any); ok {
		return s
	}
	return nil
}

// ===== Convenience property accessors ===================================

// Message returns the `message` sub-object as an Event (or nil).
func (e *Event) Message() *Event { return e.EventOf("message") }

// Text returns the message text (if any).
func (e *Event) Text() string {
	if msg := e.Message(); msg != nil {
		return msg.Str("text")
	}
	return ""
}

// Type returns the update type / author type.
func (e *Event) Type() string { return e.Str("type", "author_type") }

// IsMe reports whether the message author is the current user.
func (e *Event) IsMe(selfGUID string) bool { return e.AuthorGUID() == selfGUID }

// IsEdited reports whether the message has the `is_edited` flag.
func (e *Event) IsEdited() bool {
	msg := e.Message()
	if msg == nil {
		return false
	}
	return msg.Bool("is_edited")
}

// IsForward reports whether the message has a `forwarded_from` block.
func (e *Event) IsForward() bool {
	msg := e.Message()
	if msg == nil {
		return false
	}
	_, ok := msg.Get("forwarded_from")
	return ok
}

// IsEvent reports whether the message type is "Event".
func (e *Event) IsEvent() bool {
	msg := e.Message()
	return msg != nil && msg.Type() == "Event"
}

// IsFileInline reports whether the message contains a file_inline block.
func (e *Event) IsFileInline() bool {
	msg := e.Message()
	if msg == nil {
		return false
	}
	t := msg.Type()
	return t == "FileInline" || t == "FileInlineCaption"
}

// FileInline returns the file_inline block (if any).
func (e *Event) FileInline() *Event { return e.EventOf("file_inline") }

// IsGroup / IsChannel / IsPrivate reflect the update's chat type.
func (e *Event) IsGroup() bool   { return e.Type() == "Group" }
func (e *Event) IsChannel() bool { return e.Type() == "Channel" }
func (e *Event) IsPrivate() bool { return e.Type() == "User" }

// ObjectGUID returns the chat GUID, whichever name the server used.
func (e *Event) ObjectGUID() string {
	return e.Str("group_guid", "object_guid", "channel_guid")
}

// AuthorGUID returns the message author's GUID.
func (e *Event) AuthorGUID() string {
	msg := e.Message()
	if msg == nil {
		return ""
	}
	return msg.Str("author_object_guid")
}

// MessageID returns the message id (or pinned_message_id if that's set).
func (e *Event) MessageID() string {
	return e.Str("message_id", "pinned_message_id")
}

// ReplyMessageID returns the replied-to message id (if any).
func (e *Event) ReplyMessageID() string {
	msg := e.Message()
	if msg == nil {
		return ""
	}
	return msg.Str("reply_to_message_id")
}

// Sticker returns the sticker sub-object (if any).
func (e *Event) Sticker() *Event { return e.EventOf("sticker") }

// GUIDType returns "Channel", "Group", or "User" depending on the prefix
// of the supplied GUID (c0*, g0*, anything else).
func GUIDType(guid string) string {
	switch {
	case strings.HasPrefix(guid, "c0"):
		return "Channel"
	case strings.HasPrefix(guid, "g0"):
		return "Group"
	default:
		return "User"
	}
}
