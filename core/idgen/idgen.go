// Package idgen provides ID-generation helpers. The Rubika API
// requires several kinds of random/sequential IDs (rnd for
// messages, device_hash, tmp_session, etc.).
package idgen

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	mathrand "math/rand"
)

// Lowercase alphabet for secret generation.
const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Secret returns a cryptographically random lowercase string of
// the requested length.
func Secret(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}
	buf := make([]byte, length)
	if _, err := cryptorand.Read(buf); err != nil {
		return "", err
	}
	out := make([]byte, length)
	for i, b := range buf {
		out[i] = alphabet[int(b)%len(alphabet)]
	}
	return string(out), nil
}

// Rnd returns a random int64 in the range [1, 1_000_000] — used as
// the `rnd` field in message-sending payloads.
func Rnd() int64 { return mathrand.Int63n(1_000_000) + 1 }

// DeviceHash returns a random device-hash string (26 digits).
func DeviceHash() string {
	var b [26]byte
	for i := range b {
		b[i] = byte('0' + mathrand.Intn(10))
	}
	return string(b[:])
}

// RandomBytes returns cryptographically random bytes of the
// requested length.
func RandomBytes(length int) ([]byte, error) {
	buf := make([]byte, length)
	_, err := cryptorand.Read(buf)
	return buf, err
}

// HexID returns a random hex string of the given byte length (so
// HexID(8) returns a 16-char hex string). If crypto/rand fails the
// function panics — a predictable all-zero ID would be a security
// hazard and silently returning it is worse than crashing.
func HexID(byteLen int) string {
	b, err := RandomBytes(byteLen)
	if err != nil {
		panic("idgen: crypto/rand failed: " + err.Error())
	}
	return hex.EncodeToString(b)
}
