package matchers

import (
	"context"
	"testing"
)

func TestTextMatcher(t *testing.T) {
	ctx := context.Background()
	m := NewText("hello", false, "MessageUpdates")
	if m.Kind() != "MessageUpdates" {
		t.Errorf("Kind = %q, want MessageUpdates", m.Kind())
	}
	raw := map[string]any{"message": map[string]any{"text": "hello"}}
	if !m.Match(ctx, raw) {
		t.Error("Text matcher should match 'hello'")
	}
	raw = map[string]any{"message": map[string]any{"text": "world"}}
	if m.Match(ctx, raw) {
		t.Error("Text matcher should not match 'world'")
	}
	// Any text
	mAny := NewText("", false, "MessageUpdates")
	if !mAny.Match(ctx, map[string]any{"message": map[string]any{"text": "anything"}}) {
		t.Error("Any-text matcher should match")
	}
}

func TestCommandMatcher(t *testing.T) {
	ctx := context.Background()
	m := NewCommand([]string{"start"}, []string{"/"}, false, false, "MessageUpdates")
	raw := map[string]any{"message": map[string]any{"text": "/start now"}}
	if !m.Match(ctx, raw) {
		t.Error("Command matcher should match '/start now'")
	}
	raw = map[string]any{"message": map[string]any{"text": "/help"}}
	if m.Match(ctx, raw) {
		t.Error("Command matcher should not match '/help'")
	}
}

func TestIsGroup(t *testing.T) {
	ctx := context.Background()
	m := NewIsGroup("MessageUpdates")
	raw := map[string]any{"type": "Group"}
	if !m.Match(ctx, raw) {
		t.Error("IsGroup should match Group")
	}
	raw = map[string]any{"type": "Channel"}
	if m.Match(ctx, raw) {
		t.Error("IsGroup should not match Channel")
	}
}

func TestAndOrNot(t *testing.T) {
	ctx := context.Background()
	raw := map[string]any{
		"type":    "Group",
		"message": map[string]any{"text": "hello"},
	}
	and := NewAnd("MessageUpdates", NewIsGroup("MessageUpdates"), NewText("hello", false, "MessageUpdates"))
	if !and.Match(ctx, raw) {
		t.Error("And should match")
	}
	or := NewOr("MessageUpdates", NewIsChannel("MessageUpdates"), NewText("hello", false, "MessageUpdates"))
	if !or.Match(ctx, raw) {
		t.Error("Or should match (text matches)")
	}
	not := NewNot("MessageUpdates", NewIsChannel("MessageUpdates"))
	if !not.Match(ctx, raw) {
		t.Error("Not(Channel) should match for a Group update")
	}
}
