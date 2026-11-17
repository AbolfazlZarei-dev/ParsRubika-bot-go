package account

import (
	"encoding/json"
	"strings"
)

// redactKeys are the JSON keys whose values are treated as sensitive and
// masked in debug logs. The list covers credentials (auth, private_key,
// sign, data_enc), personally identifying data (phone_number, bio,
// password), and Rubika file/session handles (access_hash_*, file_id,
// dc_id) which could be abused to fetch a user's media.
var redactKeys = []string{
	"text", "message", "bio", "password", "phone_number", "phone",
	"auth", "private_key", "sign", "data_enc",
	"file_id", "dc_id", "access_hash", "access_hash_send", "access_hash_rec",
	"hash_file_request", "hash_file_receive",
}

// redactForLog returns a JSON-marshaled copy of v in which the values of
// sensitive keys (see redactKeys) are replaced with "[REDACTED]".
// Non-sensitive fields are preserved so the log line is still useful for
// debugging. Nested maps and slices are walked recursively.
func redactForLog(v any) string {
	cleaned := redactAny(v)
	jb, err := json.Marshal(cleaned)
	if err != nil {
		return "<unmarshalable>"
	}
	return string(jb)
}

func redactAny(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, val := range x {
			if isSensitiveKey(k) {
				out[k] = "[REDACTED]"
				continue
			}
			out[k] = redactAny(val)
		}
		return out
	case []any:
		out := make([]any, len(x))
		for i, val := range x {
			out[i] = redactAny(val)
		}
		return out
	default:
		return v
	}
}

func isSensitiveKey(k string) bool {
	lk := strings.ToLower(k)
	for _, s := range redactKeys {
		if lk == s {
			return true
		}
		// Match prefixes like "access_hash_send" → covered by exact entries,
		// but also catch future variants such as "private_key_pem".
		if strings.HasPrefix(lk, s+"_") {
			return true
		}
	}
	return false
}

// summarizeSendInput produces a compact, non-sensitive summary of a
// sendMessage input map for debug logging: the chat GUID prefix, the
// text length (NOT the text itself), and whether a reply/file is set.
func summarizeSendInput(input map[string]any) string {
	if input == nil {
		return "<nil>"
	}
	chat := stringify(input["object_guid"])
	if len(chat) > 10 {
		chat = chat[:10] + "…"
	}
	textLen := 0
	if t, ok := input["text"].(string); ok {
		textLen = len(t)
	}
	_, hasReply := input["reply_to_message_id"]
	_, hasFile := input["file_inline"]
	return "chat=" + chat + " text_len=" + itoa(textLen) +
		" reply=" + boolStr(hasReply) + " file=" + boolStr(hasFile)
}

func stringify(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func itoa(n int) string {
	// Avoid strconv import churn in this small helper file.
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
