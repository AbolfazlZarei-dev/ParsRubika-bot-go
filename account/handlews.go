package account

import (
	"context"
	"encoding/json"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/security"
)

// handleWSMessage decrypts and dispatches a single WebSocket message.
//
// Wire format:
//
//	{"data_enc": "<base64-urlsafe AES-256-CBC ciphertext>", "type": "..."}
//
// After AES-decrypting `data_enc` with the session key, the resulting
// JSON looks like:
//
//	{
//	 "user_guid": "<self guid>",
//	 "message_updates": [ { <update>, ... },
//	 "chat_updates": [ ... ],
//	 "show_activities": [ ... ],
//	 "show_notifications": [ ... ],
//	 "remove_notifications": [ ... ]
//	}
//
// For every list-valued key, each individual update is augmented with
// `user_guid` and dispatched under the CamelCase name
// (capitalize("message_updates") -> "MessageUpdates").
func (a *Account) handleWSMessage(ctx context.Context, msg map[string]any) {
	// Always recover from panics — the WebSocket read loop launches
	// this in a goroutine and a panic here would otherwise crash the
	// whole process.
	defer func() {
		if r := recover(); r != nil {
			a.errorf("ws handler panic: %v", r)
		}
	}()

	enc, ok := msg["data_enc"].(string)
	if !ok || enc == "" {
		// Status / keepalive frames (e.g. {"status":"OK","status_det":"OK"})
		// have no data_enc — silently ignore them.
		return
	}

	// Decrypt with the session key. After sign-in, `a.symmetricKey`
	// is `DeriveKey(auth)` — the same key the server uses to encrypt.
	// Use the cached []byte form (set during Request init) so this hot
	// WS frame path does not allocate per frame.
	key := []byte(a.symmetricKey)
	if key == nil && a.symmetricKey != "" {
		key = []byte(a.symmetricKey)
	}
	decrypted, err := security.Decrypt(enc, key)
	if err != nil {
		a.warnf("ws decrypt failed (len=%d): %v", len(enc), err)
		return
	}

	var data map[string]any
	if err := json.Unmarshal(decrypted, &data); err != nil {
		a.warnf("ws json decode failed: %v", err)
		a.debugf("ws undecodable payload: %q", string(decrypted))
		return
	}

	// Verbose debug logging of decrypted payload keys — only when
	// debug is explicitly enabled.
	if a.debug {
		a.debugf("DECRYPTED keys: %v", mapKeys(data))
	}

	// `user_guid` is popped and re-attached to each individual update.
	userGuid, _ := data["user_guid"].(string)
	delete(data, "user_guid")

	for name, updates := range data {
		list, ok := updates.([]any)
		if !ok {
			continue
		}
		for _, upd := range list {
			u, ok := upd.(map[string]any)
			if !ok {
				continue
			}
			u["user_guid"] = userGuid
			// Log the update name + type only when debug is enabled —
			// otherwise this fires on every WS frame and floods the log.
			if a.debug {
				updType, _ := u["type"].(string)
				a.debugf("ws update name=%s type=%s", name, updType)
			}
			a.dispatch(ctx, name, u)
		}
	}
}
