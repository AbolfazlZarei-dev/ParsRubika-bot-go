package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// This file provides at-rest encryption for the FileStore and SQLiteStore
// backends. The Record holds two genuinely secret fields — Auth (the
// 32-char symmetric-key source) and PrivateKey (the PEM RSA key). With a
// passphrase configured, the two secret fields are sealed with AES-256-GCM
// (PBKDF2-SHA256, 200k iterations) before being written, and unsealed on
// read.
//
// Format of a sealed field (text-safe, ASCII-only — safe for both JSON
// files and SQLite TEXT columns):
//
//	PRSE1.<base64-rawstd([1B ver=1][16B salt][12B nonce][ciphertext + 16B GCM tag])>
//
// Load is backward compatible: a field that does not start with the
// "PRSE1." prefix is treated as legacy plain text and returned as-is, so
// existing on-disk sessions keep working and are transparently re-sealed
// on the next Save (automatic migration).

const (
	secretPrefix        = "PRSE1."
	secretVersion  byte = 1
	secretPBKDF2It      = 200_000
	secretKeyLen        = 32
	secretSaltLen       = 16
	secretNonceLen      = 12
	secretMaxPlain      = 1 << 20 // 1 MiB hard cap on a single sealed field (auth+PEM key are far smaller)
	// secretMaxEncoded caps the base64 envelope accepted by OpenSecret
	// before it is decoded, so a maliciously large stored field cannot
	// force an unbounded allocation. ~2 MiB is comfortably above the
	// base64 expansion of secretMaxPlain.
	secretMaxEncoded = 1 << 21 // 2 MiB
)

// ErrNoPassphrase is returned when a sealed field is encountered on read
// but no passphrase was configured.
var ErrNoPassphrase = errors.New("storage: session field is encrypted but no passphrase was supplied")

// SealSecret encrypts `plain` with a key derived from `passphrase` and
// returns the "PRSE1.<base64>" envelope. When `passphrase` is empty the
// plain text is returned unchanged (legacy / opt-out mode).
func SealSecret(passphrase, plain string) (string, error) {
	if passphrase == "" {
		return plain, nil
	}
	if len(plain) > secretMaxPlain {
		return "", fmt.Errorf("storage: secret field too large (%d > %d)", len(plain), secretMaxPlain)
	}
	salt := make([]byte, secretSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := pbkdf2.Key([]byte(passphrase), salt, secretPBKDF2It, secretKeyLen, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, secretNonceLen)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ct := gcm.Seal(nil, nonce, []byte(plain), nil)
	out := make([]byte, 0, 1+secretSaltLen+secretNonceLen+len(ct))
	out = append(out, secretVersion)
	out = append(out, salt...)
	out = append(out, nonce...)
	out = append(out, ct...)
	return secretPrefix + base64.RawStdEncoding.EncodeToString(out), nil
}

// OpenSecret reverses SealSecret. A field without the "PRSE1." prefix is
// treated as legacy plain text and returned verbatim (second return is
// false, meaning "was not encrypted"). This keeps reads backward
// compatible with sessions written by older releases.
func OpenSecret(passphrase, encoded string) (string, bool, error) {
	if !strings.HasPrefix(encoded, secretPrefix) {
		return encoded, false, nil
	}
	if passphrase == "" {
		return "", false, ErrNoPassphrase
	}
	if len(encoded) > secretMaxEncoded {
		return "", false, fmt.Errorf("storage: sealed field too large (%d > %d)", len(encoded), secretMaxEncoded)
	}
	raw, err := base64.RawStdEncoding.DecodeString(encoded[len(secretPrefix):])
	if err != nil {
		return "", false, fmt.Errorf("storage: sealed field base64: %w", err)
	}
	if len(raw) < 1+secretSaltLen+secretNonceLen+aes.BlockSize {
		return "", false, fmt.Errorf("storage: sealed field too short")
	}
	if raw[0] != secretVersion {
		return "", false, fmt.Errorf("storage: sealed field version %d not supported", raw[0])
	}
	salt := raw[1 : 1+secretSaltLen]
	nonce := raw[1+secretSaltLen : 1+secretSaltLen+secretNonceLen]
	ct := raw[1+secretSaltLen+secretNonceLen:]
	key := pbkdf2.Key([]byte(passphrase), salt, secretPBKDF2It, secretKeyLen, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", false, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", false, err
	}
	// GCM authenticates the ciphertext; a wrong passphrase fails here.
	plain, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", false, fmt.Errorf("storage: sealed field decryption failed (wrong passphrase?): %w", err)
	}
	return string(plain), true, nil
}
