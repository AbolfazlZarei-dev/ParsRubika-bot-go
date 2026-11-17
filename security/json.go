package security

import "encoding/json"

// jsonEncode is the actual json.Marshal call, kept in a tiny helper so
// it can be replaced in tests if needed.
func jsonEncode(v any) ([]byte, error) {
	return json.Marshal(v)
}
