package events

import (
	"encoding/json"
	"testing"
)

// TestIntAcceptsString verifies that Int() now accepts JSON string
// values that parse as decimal integers. This is needed because
// some server fields (notably `dc_id` returned by `requestSendFile`)
// come back as JSON strings like `"480"` instead of bare numbers.
// Without this fix we would silently return 0 and the server would
// then reject the subsequent sendMessage with INVALID_INPUT.
func TestIntAcceptsString(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		key  string
		want int
	}{
		{"int", `{"dc_id":480}`, "dc_id", 480},
		{"float", `{"dc_id":480.0}`, "dc_id", 480},
		{"string decimal", `{"dc_id":"480"}`, "dc_id", 480},
		{"string negative", `{"dc_id":"-12"}`, "dc_id", -12},
		{"string zero", `{"dc_id":"0"}`, "dc_id", 0},
		{"string non-numeric stays 0", `{"dc_id":"abc"}`, "dc_id", 0},
		{"missing", `{"other":1}`, "dc_id", 0},
		{"null", `{"dc_id":null}`, "dc_id", 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var m map[string]any
			if err := json.Unmarshal([]byte(c.raw), &m); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			e := New(m)
			got := e.Int(c.key)
			if got != c.want {
				t.Errorf("Int(%q) = %d, want %d", c.key, got, c.want)
			}
		})
	}
}

// TestInt64AcceptsString is the same regression test for Int64().
func TestInt64AcceptsString(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		key  string
		want int64
	}{
		{"int", `{"size":12345}`, "size", 12345},
		{"string decimal", `{"size":"12345"}`, "size", 12345},
		{"string large", `{"size":"9007199254740993"}`, "size", 9007199254740993},
		{"string non-numeric stays 0", `{"size":"abc"}`, "size", 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var m map[string]any
			if err := json.Unmarshal([]byte(c.raw), &m); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			e := New(m)
			got := e.Int64(c.key)
			if got != c.want {
				t.Errorf("Int64(%q) = %d, want %d", c.key, got, c.want)
			}
		})
	}
}
