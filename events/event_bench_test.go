package events

import (
	"encoding/json"
	"testing"
)

// BenchmarkFindKey measures the cost of a recursive key lookup on a
// moderately nested payload (typical message update).
func BenchmarkFindKey(b *testing.B) {
	raw := map[string]any{
		"object_guid": "g0GROUP",
		"message": map[string]any{
			"text":               "hello world",
			"message_id":         "m1",
			"author_object_guid": "u0USER",
			"file_inline": map[string]any{
				"file_id": "f1",
				"dc_id":   1,
				"size":    1024,
			},
		},
	}
	u := New(raw)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = u.Find("file_id")
	}
}

// BenchmarkStr measures the Str() helper cost.
func BenchmarkStr(b *testing.B) {
	u := New(map[string]any{"message": map[string]any{"text": "hello"}})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = u.Str("text")
	}
}

// BenchmarkJSON measures the JSON serialization cost.
func BenchmarkJSON(b *testing.B) {
	u := New(map[string]any{
		"object_guid": "g0GROUP",
		"text":        "hello world",
		"n":           42,
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = u.JSON("")
	}
}

// BenchmarkUnmarshal measures the JSON parsing cost of a typical
// server response.
func BenchmarkUnmarshal(b *testing.B) {
	raw := []byte(`{"object_guid":"g0GROUP","message":{"text":"hello","message_id":"m1","author_object_guid":"u0USER"}}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var m map[string]any
		_ = json.Unmarshal(raw, &m)
	}
}

// BenchmarkNew measures the cost of constructing a new *Event from an
// existing map. With the lazy-External optimisation this should be a
// single allocation (the *Event itself); the External map is no
// longer allocated eagerly.
func BenchmarkNew(b *testing.B) {
	raw := map[string]any{
		"object_guid": "g0GROUP",
		"message":     map[string]any{"text": "hello", "message_id": "m1"},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(raw)
	}
}
