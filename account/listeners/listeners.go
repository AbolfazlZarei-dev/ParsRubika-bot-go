// Package listeners provides the five built-in listener matchers used
// by the account client: Messages, Chats, Activities, Notifications,
// and Removals. Each one is just a thin wrapper that fixes the
// matcher's Kind() to the CamelCase name the server actually sends
// (the `capitalize()` of the snake_case key in the decrypted
// `data_enc` payload — see account/handlews.go).
package listeners

import (
	"context"
)

// Messages is the matcher for "MessageUpdates" updates (server key
// `message_updates`). When no sub-matchers are supplied it matches
// every message update.
type Messages struct{ inner MatcherChain }

// Chats matches "ChatUpdates" updates (server key `chat_updates`).
type Chats struct{ inner MatcherChain }

// Activities matches "ShowActivities" updates.
type Activities struct{ inner MatcherChain }

// Notifications matches "ShowNotifications" updates.
type Notifications struct{ inner MatcherChain }

// Removals matches "RemoveNotifications" updates.
type Removals struct{ inner MatcherChain }

// ScheduledMessages matches "ScheduleMessageUpdates" updates.
type ScheduledMessages struct{ inner MatcherChain }

// DraftMessages matches "DraftMessageUpdate" updates.
type DraftMessages struct{ inner MatcherChain }

// UnconfirmedSessions matches "UnconfirmedSessionUpdates" updates.
type UnconfirmedSessions struct{ inner MatcherChain }

// VoiceChatParticipants matches "GroupVoiceChatParticipantUpdates" updates.
type VoiceChatParticipants struct{ inner MatcherChain }

// VoiceChatUpdates matches "GroupVoiceChatUpdates" updates.
type VoiceChatUpdates struct{ inner MatcherChain }

// CallUpdates matches "CallUpdates" updates.
type CallUpdates struct{ inner MatcherChain }

// CallSignals matches "CallSignalData" updates.
type CallSignals struct{ inner MatcherChain }

// MatcherChain is the common interface implemented by every listener.
// It is identical to the account.Matcher interface; it is re-declared
// here to avoid an import cycle.
type MatcherChain interface {
	Kind() string
	Match(ctx context.Context, raw map[string]any) bool
}

// ===== Kind implementations ============================================
//
// These names MUST match `capitalize(<server_key>)`:
//
// capitalize("message_updates") = "MessageUpdates"
// capitalize("chat_updates") = "ChatUpdates"
// capitalize("show_activities") = "ShowActivities"
// capitalize("show_notifications") = "ShowNotifications"
// capitalize("remove_notifications") = "RemoveNotifications"
//
// The list below is the complete set of update names the server sends.

func (Messages) Kind() string              { return "MessageUpdates" }
func (Chats) Kind() string                 { return "ChatUpdates" }
func (Activities) Kind() string            { return "ShowActivities" }
func (Notifications) Kind() string         { return "ShowNotifications" }
func (Removals) Kind() string              { return "RemoveNotifications" }
func (ScheduledMessages) Kind() string     { return "ScheduleMessageUpdates" }
func (DraftMessages) Kind() string         { return "DraftMessageUpdate" }
func (UnconfirmedSessions) Kind() string   { return "UnconfirmedSessionUpdates" }
func (VoiceChatParticipants) Kind() string { return "GroupVoiceChatParticipantUpdates" }
func (VoiceChatUpdates) Kind() string      { return "GroupVoiceChatUpdates" }
func (CallUpdates) Kind() string           { return "CallUpdates" }
func (CallSignals) Kind() string           { return "CallSignalData" }

// ===== Match implementations ==========================================

func (m Messages) Match(ctx context.Context, raw map[string]any) bool {
	if m.inner == nil {
		return true
	}
	return m.inner.Match(ctx, raw)
}

func (c Chats) Match(ctx context.Context, raw map[string]any) bool {
	if c.inner == nil {
		return true
	}
	return c.inner.Match(ctx, raw)
}

func (a Activities) Match(ctx context.Context, raw map[string]any) bool {
	if a.inner == nil {
		return true
	}
	return a.inner.Match(ctx, raw)
}

func (n Notifications) Match(ctx context.Context, raw map[string]any) bool {
	if n.inner == nil {
		return true
	}
	return n.inner.Match(ctx, raw)
}

func (r Removals) Match(ctx context.Context, raw map[string]any) bool {
	if r.inner == nil {
		return true
	}
	return r.inner.Match(ctx, raw)
}

func (s ScheduledMessages) Match(ctx context.Context, raw map[string]any) bool {
	if s.inner == nil {
		return true
	}
	return s.inner.Match(ctx, raw)
}

func (d DraftMessages) Match(ctx context.Context, raw map[string]any) bool {
	if d.inner == nil {
		return true
	}
	return d.inner.Match(ctx, raw)
}

func (u UnconfirmedSessions) Match(ctx context.Context, raw map[string]any) bool {
	if u.inner == nil {
		return true
	}
	return u.inner.Match(ctx, raw)
}

func (v VoiceChatParticipants) Match(ctx context.Context, raw map[string]any) bool {
	if v.inner == nil {
		return true
	}
	return v.inner.Match(ctx, raw)
}

func (v VoiceChatUpdates) Match(ctx context.Context, raw map[string]any) bool {
	if v.inner == nil {
		return true
	}
	return v.inner.Match(ctx, raw)
}

func (c CallUpdates) Match(ctx context.Context, raw map[string]any) bool {
	if c.inner == nil {
		return true
	}
	return c.inner.Match(ctx, raw)
}

func (c CallSignals) Match(ctx context.Context, raw map[string]any) bool {
	if c.inner == nil {
		return true
	}
	return c.inner.Match(ctx, raw)
}

// ===== Constructors =====================================================

// NewMessages constructs a Messages listener that wraps the supplied
// matcher chain.
func NewMessages(inner MatcherChain) Messages { return Messages{inner: inner} }

// NewChats constructs a Chats listener.
func NewChats(inner MatcherChain) Chats { return Chats{inner: inner} }

// NewActivities constructs an Activities listener.
func NewActivities(inner MatcherChain) Activities { return Activities{inner: inner} }

// NewNotifications constructs a Notifications listener.
func NewNotifications(inner MatcherChain) Notifications { return Notifications{inner: inner} }

// NewRemovals constructs a Removals listener.
func NewRemovals(inner MatcherChain) Removals { return Removals{inner: inner} }

// NewScheduledMessages constructs a ScheduledMessages listener.
func NewScheduledMessages(inner MatcherChain) ScheduledMessages {
	return ScheduledMessages{inner: inner}
}

// NewDraftMessages constructs a DraftMessages listener.
func NewDraftMessages(inner MatcherChain) DraftMessages { return DraftMessages{inner: inner} }

// NewUnconfirmedSessions constructs an UnconfirmedSessions listener.
func NewUnconfirmedSessions(inner MatcherChain) UnconfirmedSessions {
	return UnconfirmedSessions{inner: inner}
}

// NewVoiceChatParticipants constructs a VoiceChatParticipants listener.
func NewVoiceChatParticipants(inner MatcherChain) VoiceChatParticipants {
	return VoiceChatParticipants{inner: inner}
}

// NewVoiceChatUpdates constructs a VoiceChatUpdates listener.
func NewVoiceChatUpdates(inner MatcherChain) VoiceChatUpdates { return VoiceChatUpdates{inner: inner} }

// NewCallUpdates constructs a CallUpdates listener.
func NewCallUpdates(inner MatcherChain) CallUpdates { return CallUpdates{inner: inner} }

// NewCallSignals constructs a CallSignals listener.
func NewCallSignals(inner MatcherChain) CallSignals { return CallSignals{inner: inner} }
