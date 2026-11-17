package ai

import "strings"

// extractPrompt strips the command (and optional prefix) from text and
// returns the remainder, trimmed of surrounding whitespace. It is
// case-insensitive on the command token but preserves the original
// casing of the prompt itself.
//
// Examples (Command="ask", Prefixes=["/"]):
//
//	"/ask hello"      -> "hello"
//	"/Ask   world "   -> "world"
//	"ask foo bar"     -> "foo bar"   (prefix-less form is accepted too)
//	"/ask"            -> ""          (no prompt)
func extractPrompt(text, cmd string, prefixes []string) string {
	t := strings.TrimSpace(text)
	if t == "" {
		return ""
	}
	low := strings.ToLower(t)
	for _, p := range prefixes {
		token := p + strings.ToLower(cmd)
		if low == token {
			return ""
		}
		if strings.HasPrefix(low, token) {
			rest := t[len(token):]
			return strings.TrimSpace(rest)
		}
	}
	// Allow the no-prefix form too.
	token := strings.ToLower(cmd)
	if low == token {
		return ""
	}
	if strings.HasPrefix(low, token) {
		rest := t[len(token):]
		return strings.TrimSpace(rest)
	}
	return ""
}
