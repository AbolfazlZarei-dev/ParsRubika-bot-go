package security

import (
	"crypto/x509"
	"encoding/pem"
	"strings"
)

// NormalizePrivateKey cleans up a user-supplied RSA private-key string:
//
// - Strips any "[REDACTED:ssh_private_key]" redaction marker that may
// have been pasted in (the marker itself is removed, not just the
// text before it).
// - Ensures the PEM body has a proper "-----BEGIN …-----" / "-----END …-----"
// envelope, appending the footer only if it is genuinely missing.
// - Validates the result by parsing it as a PEM-encoded private key.
//
// Returns the cleaned string. If the input cannot be parsed as a PEM
// private key the original (cleaned) string is returned unchanged so
// that the caller's subsequent ParsePrivateKey call surfaces a clear
// error instead of a silently-mangled key.
func NormalizePrivateKey(s string) string {
	if s == "" {
		return s
	}
	// Remove the redaction markers wherever they appear. Both the
	// specific "[REDACTED:ssh_private_key]" marker and the generic
	// "[REDACTED]" prefix (which some tools prepend) are stripped.
	s = strings.ReplaceAll(s, "[REDACTED:ssh_private_key]", "")
	s = strings.ReplaceAll(s, "[REDACTED]", "")

	// Trim leading whitespace/quotes that often accompany pasted keys.
	s = strings.TrimSpace(s)

	// If the key has a BEGIN header but is missing its END footer, append
	// the standard RSA footer. We only do this when there is no END at
	// all — appending to a key that already has a (possibly differently
	// spaced) END would produce invalid PEM.
	hasBegin := strings.Contains(s, "-----BEGIN")
	hasEnd := strings.Contains(s, "-----END")
	if hasBegin && !hasEnd {
		s += "\n-----END RSA PRIVATE KEY-----"
	}

	// Best-effort validation: if we can parse it as PEM + PKCS1/PKCS8,
	// we know the structure is sound. If not, return the cleaned string
	// and let the caller's parser produce the error.
	if block, _ := pem.Decode([]byte(s)); block != nil {
		if _, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
			return s
		}
		if _, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			return s
		}
	}
	return s
}
