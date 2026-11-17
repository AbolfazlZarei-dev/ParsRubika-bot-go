package events

import (
	"encoding/json"
	"testing"
)

func TestFindKeySimple(t *testing.T) {
	u := New(map[string]any{
		"object_guid": "g0GROUP",
		"message": map[string]any{
			"text":               "hello",
			"message_id":         "m1",
			"author_object_guid": "u0USER",
		},
	})
	if u.ObjectGUID() != "g0GROUP" {
		t.Errorf("ObjectGUID = %q, want g0GROUP", u.ObjectGUID())
	}
	if u.Text() != "hello" {
		t.Errorf("Text = %q, want hello", u.Text())
	}
	if u.AuthorGUID() != "u0USER" {
		t.Errorf("AuthorGUID = %q, want u0USER", u.AuthorGUID())
	}
	if u.MessageID() != "m1" {
		t.Errorf("MessageID = %q, want m1", u.MessageID())
	}
}

func TestFindKeyDeepNested(t *testing.T) {
	u := New(map[string]any{
		"outer": map[string]any{
			"middle": map[string]any{
				"inner": "found",
			},
		},
	})
	if got := u.Find("inner"); got != "found" {
		t.Errorf("Find(inner) = %v, want found", got)
	}
}

func TestGUIDType(t *testing.T) {
	cases := map[string]string{
		"c0CHANNEL": "Channel",
		"g0GROUP":   "Group",
		"u0USER":    "User",
		"":          "User",
	}
	for in, want := range cases {
		if got := GUIDType(in); got != want {
			t.Errorf("GUIDType(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestJSONRoundTrip(t *testing.T) {
	u := New(map[string]any{
		"object_guid": "g0GROUP",
		"text":        "hello",
		"n":           42,
	})
	js := u.JSON("")
	if !contains(js, "g0GROUP") || !contains(js, "hello") {
		t.Fatalf("JSON missing expected content: %s", js)
	}
	var roundTripped map[string]any
	if err := json.Unmarshal([]byte(js), &roundTripped); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestExternalData(t *testing.T) {
	u := New(map[string]any{})
	u.SetExternal("foo", 42)
	v, ok := u.GetExternal("foo")
	if !ok || v != 42 {
		t.Fatalf("GetExternal = %v (ok=%v), want 42", v, ok)
	}
}

func TestPropertyAccessors(t *testing.T) {
	u := New(map[string]any{
		"type": "Group",
		"message": map[string]any{
			"text":           "hello",
			"is_edited":      true,
			"forwarded_from": map[string]any{"type_from": "User"},
		},
	})
	if !u.IsGroup() {
		t.Error("IsGroup should be true")
	}
	if u.IsChannel() {
		t.Error("IsChannel should be false")
	}
	if u.IsPrivate() {
		t.Error("IsPrivate should be false")
	}
	if !u.IsEdited() {
		t.Error("IsEdited should be true")
	}
	if !u.IsForward() {
		t.Error("IsForward should be true")
	}
}
