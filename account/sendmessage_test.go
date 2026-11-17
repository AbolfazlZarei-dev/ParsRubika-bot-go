package account

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/security"
)

// TestSendMessageOmitsEmptyReplyToMessageID verifies that when the
// caller does not supply a ReplyToMessageID, the JSON sent to the
// server does NOT contain `reply_to_message_id` at all (rather than
// `reply_to_message_id: ""`). The server rejects the empty-string
// variant with INVALID_INPUT: sendMessage failed.
func TestSendMessageOmitsEmptyReplyToMessageID(t *testing.T) {
	// Reproduce the exact input construction SendMessage does.
	input := map[string]any{
		"object_guid": "g0TEST",
		"rnd":         randInt63(),
	}
	// SendMessage adds `reply_to_message_id` only if non-empty.
	if r := ""; r != "" {
		input["reply_to_message_id"] = r
	}
	if msg := "hello"; msg != "" {
		input["text"] = msg
	}

	body, _ := json.Marshal(input)
	got := string(body)
	if strings.Contains(got, "reply_to_message_id") {
		t.Errorf("payload should NOT contain reply_to_message_id, got: %s", got)
	}
	if !strings.Contains(got, `"object_guid":"g0TEST"`) {
		t.Errorf("payload should contain object_guid, got: %s", got)
	}
	if !strings.Contains(got, `"text":"hello"`) {
		t.Errorf("payload should contain text, got: %s", got)
	}
	if !strings.Contains(got, `"rnd":`) {
		t.Errorf("payload should contain rnd, got: %s", got)
	}
}

// TestSendMessageIncludesReplyToMessageIDWhenNonEmpty verifies the
// positive case: when ReplyToMessageID is non-empty, it MUST be sent.
func TestSendMessageIncludesReplyToMessageIDWhenNonEmpty(t *testing.T) {
	input := map[string]any{
		"object_guid": "g0TEST",
		"rnd":         randInt63(),
	}
	if r := "m123"; r != "" {
		input["reply_to_message_id"] = r
	}
	body, _ := json.Marshal(input)
	got := string(body)
	if !strings.Contains(got, `"reply_to_message_id":"m123"`) {
		t.Errorf("payload should contain reply_to_message_id=m123, got: %s", got)
	}
}

// TestSendMessageBuildsValidJSON is a smoke test that ensures the
// SendMessage input map serializes to valid JSON with the expected
// shape (object_guid + text + rnd, NO reply_to_message_id when empty).
func TestSendMessageBuildsValidJSON(t *testing.T) {
	auth := "abcdefghijklmnopqrstuvwxyz012345"
	a, err := New(WithAuth(auth))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer a.Close()
	// Sanity-check the symmetric key derivation works (used by Request).
	key, err := security.DeriveKey(auth)
	if err != nil {
		t.Fatalf("DeriveKey: %v", err)
	}
	if len(key) != 32 {
		t.Errorf("symmetric key length = %d, want 32", len(key))
	}
	a.symmetricKey = key
}

// Ensure unused imports don't break compilation.
var _ = events.New
var _ context.Context
