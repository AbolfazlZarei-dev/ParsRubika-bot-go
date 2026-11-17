package security

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const lowercaseAlphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomSecret returns a cryptographically random lowercase string of
// the requested length. Returns an empty string when length <= 0.
//
// Each character is drawn uniformly from the alphabet using
// crypto/rand.Int, which performs rejection sampling internally and so
// avoids the modulo bias that a naive `byte % len(alphabet)` mapping
// would introduce (with a 26-letter alphabet, 256 % 26 == 22, so the
// last four letters would otherwise be ~10% under-represented).
func RandomSecret(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}
	max := big.NewInt(int64(len(lowercaseAlphabet)))
	out := make([]byte, length)
	for i := range out {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("security: rand int: %w", err)
		}
		out[i] = lowercaseAlphabet[idx.Int64()]
	}
	return string(out), nil
}

// RandomBytes returns a cryptographically random byte slice of the
// requested length.
func RandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, nil
	}
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("security: rand read: %w", err)
	}
	return buf, nil
}
