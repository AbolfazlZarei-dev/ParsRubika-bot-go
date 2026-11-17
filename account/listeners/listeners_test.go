package listeners

import (
	"context"
	"testing"
)

type trueMatcher struct{ kind string }

func (m trueMatcher) Kind() string                                   { return m.kind }
func (m trueMatcher) Match(_ context.Context, _ map[string]any) bool { return true }

type falseMatcher struct{ kind string }

func (m falseMatcher) Kind() string                                   { return m.kind }
func (m falseMatcher) Match(_ context.Context, _ map[string]any) bool { return false }

func TestMessagesListener(t *testing.T) {
	m := NewMessages(trueMatcher{kind: "MessageUpdates"})
	if m.Kind() != "MessageUpdates" {
		t.Errorf("Kind = %q, want MessageUpdates", m.Kind())
	}
	if !m.Match(context.Background(), map[string]any{}) {
		t.Error("Messages listener should match when inner matcher returns true")
	}
	mFalse := NewMessages(falseMatcher{kind: "MessageUpdates"})
	if mFalse.Match(context.Background(), map[string]any{}) {
		t.Error("Messages listener should not match when inner matcher returns false")
	}
	// A nil inner matcher matches everything.
	mNil := Messages{}
	if !mNil.Match(context.Background(), map[string]any{}) {
		t.Error("Messages with nil inner should match everything")
	}
}

func TestChatsListener(t *testing.T) {
	m := NewChats(trueMatcher{kind: "ChatUpdates"})
	if m.Kind() != "ChatUpdates" {
		t.Errorf("Kind = %q, want ChatUpdates", m.Kind())
	}
}

func TestActivitiesListener(t *testing.T) {
	m := NewActivities(trueMatcher{kind: "ShowActivities"})
	if m.Kind() != "ShowActivities" {
		t.Errorf("Kind = %q, want ShowActivities", m.Kind())
	}
}

func TestNotificationsListener(t *testing.T) {
	m := NewNotifications(trueMatcher{kind: "ShowNotifications"})
	if m.Kind() != "ShowNotifications" {
		t.Errorf("Kind = %q, want ShowNotifications", m.Kind())
	}
}

func TestRemovalsListener(t *testing.T) {
	m := NewRemovals(trueMatcher{kind: "RemoveNotifications"})
	if m.Kind() != "RemoveNotifications" {
		t.Errorf("Kind = %q, want RemoveNotifications", m.Kind())
	}
}
