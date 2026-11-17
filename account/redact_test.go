package account

import (
	"strings"
	"testing"
)

func TestRedactForLogMasksSecrets(t *testing.T) {
	in := map[string]any{
		"object_guid":         "u0abcdef",
		"text":                "this is a private message body",
		"reply_to_message_id": "12345",
		"file_inline": map[string]any{
			"file_id":         "FILE123",
			"dc_id":           3,
			"access_hash_rec": "ACC456",
		},
		"auth":         "abcdefghijklmnopqrstuvwxyz012345",
		"private_key":  "-----BEGIN RSA PRIVATE KEY-----\nDATA\n-----END RSA PRIVATE KEY-----",
		"sign":         "base64signature",
		"data_enc":     "base64ciphertext",
		"phone_number": "+989121234567",
		"bio":          "user biography",
		"password":     "user-2fa-password",
		"status":       "OK",
		"status_det":   "OK",
	}
	out := redactForLog(in)

	// No sensitive value may appear verbatim in the redacted output.
	for _, secret := range []string{
		"this is a private message body",
		"FILE123", "ACC456",
		"abcdefghijklmnopqrstuvwxyz012345",
		"BEGIN RSA PRIVATE KEY",
		"base64signature", "base64ciphertext",
		"+989121234567", "user biography", "user-2fa-password",
	} {
		if strings.Contains(out, secret) {
			t.Errorf("redactForLog leaked %q:\n%s", secret, out)
		}
	}
	// Non-sensitive fields must remain visible so the log is useful.
	for _, want := range []string{"object_guid", "u0abcdef", "status", "OK"} {
		if !strings.Contains(out, want) {
			t.Errorf("redactForLog dropped non-sensitive %q:\n%s", want, out)
		}
	}
	// Sensitive keys must be masked with the marker.
	for _, k := range []string{"text", "auth", "private_key", "sign", "data_enc", "password", "bio", "phone_number"} {
		if !strings.Contains(out, "\""+k+"\":\"[REDACTED]\"") {
			t.Errorf("redactForLog did not mask %q:\n%s", k, out)
		}
	}
}

func TestSummarizeSendInputNoLeak(t *testing.T) {
	in := map[string]any{
		"object_guid":         "u0someuserguid123",
		"text":                "secret text content",
		"reply_to_message_id": "99",
		"file_inline":         map[string]any{"file_id": "X"},
	}
	summary := summarizeSendInput(in)
	if strings.Contains(summary, "secret text content") {
		t.Errorf("summarizeSendInput leaked text: %s", summary)
	}
	if !strings.Contains(summary, "text_len=") {
		t.Errorf("summary missing text_len: %s", summary)
	}
	if !strings.Contains(summary, "reply=true") || !strings.Contains(summary, "file=true") {
		t.Errorf("summary missing flags: %s", summary)
	}
}
