package events

import (
	"encoding/json"
	"testing"
)

// TestAnyStr verifies that AnyStr coerces all common JSON types
// to their string form. This is the helper used by sendFileType
// to read `file_id` / `access_hash_send` from the requestSendFile
// response, where the server may return either a string or a number.
func TestAnyStr(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		key  string
		want string
	}{
		{"string", `{"id":"abc123"}`, "id", "abc123"},
		{"int", `{"id":12345}`, "id", "12345"},
		{"float whole", `{"id":12345.0}`, "id", "12345"},
		{"float fractional", `{"id":1.5}`, "id", "1.5"},
		{"negative int", `{"id":-42}`, "id", "-42"},
		{"missing key", `{"other":1}`, "id", ""},
		{"null", `{"id":null}`, "id", ""},
		{"nested", `{"data":{"id":999}}`, "id", "999"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var m map[string]any
			if err := json.Unmarshal([]byte(c.raw), &m); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			e := New(m)
			got := e.AnyStr(c.key)
			if got != c.want {
				t.Errorf("AnyStr(%q) = %q, want %q", c.key, got, c.want)
			}
		})
	}
}
