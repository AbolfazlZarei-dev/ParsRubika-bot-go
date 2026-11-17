package matchers

import (
	"context"
	"testing"
)

func TestOperatorEqual(t *testing.T) {
	op := Operator{Kind: OpEqual, Value: "hello"}
	if !op.Match("hello") {
		t.Error("Equal should match")
	}
	if op.Match("world") {
		t.Error("Equal should not match")
	}
}

func TestOperatorNumeric(t *testing.T) {
	op := Operator{Kind: OpGreater, Value: 10}
	if !op.Match(15) {
		t.Error("15 > 10 should match")
	}
	if op.Match(5) {
		t.Error("5 > 10 should not match")
	}
	op = Operator{Kind: OpLessEqual, Value: 10}
	if !op.Match(10) {
		t.Error("10 <= 10 should match")
	}
}

func TestNewMessageID(t *testing.T) {
	m := NewMessageID(OpEqual, "m1")
	if m.Kind() != "MessageUpdates" {
		t.Errorf("Kind = %q, want MessageUpdates", m.Kind())
	}
	raw := map[string]any{
		"message": map[string]any{
			"message_id": "m1",
		},
	}
	if !m.Match(context.Background(), raw) {
		t.Error("Should match m1")
	}
}

func TestIsVideoMatcher(t *testing.T) {
	m := NewIsVideo("MessageUpdates")
	// The Event wrapper's IsFileInline() checks message.type, and
	// FileInline() reads the nested file_inline block.
	raw := map[string]any{
		"message": map[string]any{
			"type": "FileInline",
			"file_inline": map[string]any{
				"type": "Video",
			},
		},
	}
	if !m.Match(context.Background(), raw) {
		t.Error("Should match Video")
	}
}

func TestEventType(t *testing.T) {
	m := NewEventType("AddedGroupMembers", "MessageUpdates")
	raw := map[string]any{
		"message": map[string]any{
			"type": "Event",
			"event_data": map[string]any{
				"type": "AddedGroupMembers",
			},
		},
	}
	if !m.Match(context.Background(), raw) {
		t.Error("Should match event type")
	}
}
