package security

import (
	"fmt"
	"regexp"
	"strings"
)

// DecodeAuth applies the character-substitution used to derive the
// `decode_auth` payload. The mapping is:
//
// - lowercase letters: shifted by (32 - (c - 'a')) mod 26
// - uppercase letters: shifted by (29 - (c - 'A')) mod 26
// - digits: shifted by (13 - (c - '0')) mod 10
//
// All other characters (symbols, dashes, underscores) are kept as-is.
// Go's `%` operator follows the sign of the dividend, so we use the
// `((x % m) + m) % m` idiom to obtain an always-positive
// modulo.
func DecodeAuth(auth string) string {
	var b strings.Builder
	b.Grow(len(auth))
	for _, c := range auth {
		switch {
		case c >= 'a' && c <= 'z':
			shifted := ((32-(int(c)-97))%26 + 26) % 26
			b.WriteRune(rune(shifted + 97))
		case c >= 'A' && c <= 'Z':
			shifted := ((29-(int(c)-65))%26 + 26) % 26
			b.WriteRune(rune(shifted + 65))
		case c >= '0' && c <= '9':
			shifted := ((13-(int(c)-48))%10 + 10) % 10
			b.WriteRune(rune(shifted + 48))
		default:
			b.WriteRune(c)
		}
	}
	return b.String()
}

// chunkRE splits a 32-char auth string into 4 chunks of 8 chars each.
var chunkRE = regexp.MustCompile(`\S{8}`)

// DeriveKey derives the AES-256 symmetric key from a 32-character auth
// string by:
//
// 1. Splitting auth into 4 chunks of 8 characters.
// 2. Reordering them as chunks[2] + chunks[0] + chunks[3] + chunks[1].
// 3. For each character: applying `((c - 'a' + 9) mod 26) + 'a'`.
//
// The result is the 32-character AES-256 key (the underlying bytes are
// the ASCII codes of these characters, exactly as the server does).
func DeriveKey(auth string) (string, error) {
	if len(auth) != 32 {
		return "", fmt.Errorf("security: auth length must be %d, got %d", 32, len(auth))
	}
	chunks := chunkRE.FindAllString(auth, -1)
	if len(chunks) != 4 {
		return "", fmt.Errorf("security: cannot split auth %q into 4 chunks of 8 chars", auth)
	}
	var b strings.Builder
	b.Grow(32)
	for _, c := range chunks[2] + chunks[0] + chunks[3] + chunks[1] {
		shifted := ((int(c)-97+9)%26 + 26) % 26
		b.WriteRune(rune(shifted + 97))
	}
	return b.String(), nil
}
