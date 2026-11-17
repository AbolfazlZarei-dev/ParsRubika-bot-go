// Package convert provides type-conversion helpers for the loose
// `map[string]any` payloads returned by the Rubika API. These
// helpers are allocation-free where possible.
package convert

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ToString converts any scalar to its string form.
func ToString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		if float64(int64(x)) == x {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case bool:
		if x {
			return "true"
		}
		return "false"
	case nil:
		return ""
	default:
		b, _ := json.Marshal(x)
		return string(b)
	}
}

// ToInt converts any numeric scalar to int.
func ToInt(v any) int {
	switch x := v.(type) {
	case int:
		return x
	case int64:
		return int(x)
	case float64:
		return int(x)
	case string:
		n, _ := strconv.Atoi(x)
		return n
	case json.Number:
		n, _ := x.Int64()
		return int(n)
	}
	return 0
}

// ToInt64 converts any numeric scalar to int64.
func ToInt64(v any) int64 {
	switch x := v.(type) {
	case int:
		return int64(x)
	case int64:
		return x
	case float64:
		return int64(x)
	case string:
		n, _ := strconv.ParseInt(x, 10, 64)
		return n
	case json.Number:
		n, _ := x.Int64()
		return n
	}
	return 0
}

// ToBool converts any scalar to bool.
func ToBool(v any) bool {
	switch x := v.(type) {
	case bool:
		return x
	case string:
		return x == "true" || x == "1"
	case int:
		return x != 0
	case int64:
		return x != 0
	case float64:
		return x != 0
	}
	return false
}

// ToJSON marshals `v` to a JSON string.
func ToJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// FromJSON unmarshals `s` into `out`.
func FromJSON(s string, out any) error {
	return json.Unmarshal([]byte(s), out)
}

// MapToJSON converts a map to a JSON string.
func MapToJSON(m map[string]any) string { return ToJSON(m) }

// JSONToMap parses a JSON string into a map[string]any.
func JSONToMap(s string) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, fmt.Errorf("convert: json: %w", err)
	}
	return m, nil
}
