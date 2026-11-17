package rubino

import "encoding/json"

// jsonUnmarshal is a small indirection so the studio package doesn't
// have to import encoding/json in every file.
func jsonUnmarshal(data []byte, out *map[string]any) error {
	return json.Unmarshal(data, out)
}
