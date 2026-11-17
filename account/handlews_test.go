package account

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account/listeners"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account/matchers"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/security"
)

// TestHandleWSMessageDispatch is a small end-to-end test that:
// 1. Sets up an Account with a known auth (so we know the symmetric key).
// 2. Registers a /start listener.
// 3. Builds a fake "decrypted" payload (the same shape the server sends),
// AES-encrypts it with the session key, and feeds it to handleWSMessage.
// 4. Asserts the listener fires.
//
// This proves the dispatch chain works: ws -> decrypt -> dispatch -> listener.
func TestHandleWSMessageDispatch(t *testing.T) {
	// Use a 32-char auth so DeriveKey succeeds.
	auth := "abcdefghijklmnopqrstuvwxyz012345"

	a, err := New(WithAuth(auth))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer a.Close()

	// Derive the symmetric key the same way Start does.
	key, err := security.DeriveKey(auth)
	if err != nil {
		t.Fatalf("DeriveKey: %v", err)
	}
	a.symmetricKey = key

	// Make dispatch sequential so we can deterministically check the
	// callback was invoked.
	a.sequential = true

	var got int32
	AddListener(a, func(ctx context.Context, e *events.Event) error {
		if e.Text() != "/start" {
			t.Errorf("unexpected text: %q", e.Text())
		}
		atomic.AddInt32(&got, 1)
		return nil
	}, listeners.NewMessages(matchers.NewText("/start", false, "MessageUpdates")))

	// Build the decrypted payload the same way the server would.
	plain := map[string]any{
		"user_guid": "u0ME",
		"message_updates": []any{
			map[string]any{
				"object_guid": "g0CHAT",
				"message":     map[string]any{"text": "/start", "message_id": "m1"},
			},
		},
		"chat_updates": []any{
			map[string]any{"object_guid": "g0CHAT2"},
		},
	}

	// AES-encrypt it with the session key (same params as server).
	raw, _ := json.Marshal(plain)
	enc, err := security.Encrypt(raw, []byte(key))
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	// Feed it as if it came from the WebSocket.
	a.handleWSMessage(context.Background(), map[string]any{
		"data_enc": enc,
		"type":     "Update",
	})

	if got != 1 {
		t.Fatalf("listener should fire exactly once for /start, got %d", got)
	}
}
