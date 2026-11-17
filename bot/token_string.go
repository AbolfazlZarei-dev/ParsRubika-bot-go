package bot

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// TokenStringVersion is the current format version of the encrypted
// token string.
const TokenStringVersion byte = 1

// TokenStringPrefix is the 4-character marker that identifies a token
// string produced by ExportTokenString.
const TokenStringPrefix = "PRTS" // ParsRubika Token String

// maxTokenDecompressed caps how large the zlib-decompressed token-string
// payload may be. A token payload is a tiny JSON map (~120 bytes); the
// cap exists purely as a decompression-bomb guard so a maliciously
// crafted (or corrupted) ciphertext cannot make the importer allocate
// unbounded memory. The cap is generous to leave room for future fields.
const maxTokenDecompressed = 1 << 16 // 64 KiB

// ExportTokenString encrypts the bot token (along with optional metadata)
// into a portable, password-protected string. The format is:
//
// PRTS<base64(version | salt | nonce | ciphertext | checksum)>
//
// The ciphertext is AES-256-GCM, the key is derived via PBKDF2-SHA256
// with 200,000 iterations, so tampering is detected on import.
//
// If password is empty, a deterministic but insecure "no-password" mode
// is used — the resulting string is still obfuscated but NOT encrypted,
// and a one-time warning is logged. Always supply a strong password for
// any real use; treat the empty-password form as a development convenience.
func (b *Bot) ExportTokenString(password string) (string, error) {
	if b.token == "" {
		return "", errors.New("bot: no token to export")
	}
	if password == "" {
		log.Printf("[bot] WARNING: ExportTokenString called with an EMPTY password — the resulting " +
			"token string is obfuscated only, NOT encrypted. Anyone who obtains it can recover the bot token.")
	}
	payload := map[string]any{
		"token":      b.token,
		"base_url":   b.baseURL,
		"user_agent": b.opts.UserAgent,
		"parse_mode": string(b.opts.ParseMode),
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("bot: marshal token payload: %w", err)
	}
	// zlib-compress the payload to reduce size.
	var comp bytes.Buffer
	zw := zlib.NewWriter(&comp)
	if _, err := zw.Write(raw); err != nil {
		return "", err
	}
	if err := zw.Close(); err != nil {
		return "", err
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	nonce := make([]byte, 12) // GCM standard nonce size
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	key := deriveTokenKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nil, nonce, comp.Bytes(), nil)
	// Build the binary blob: version(1) | salt(16) | nonce(12) | ciphertext
	blob := make([]byte, 0, 1+len(salt)+len(nonce)+len(ciphertext))
	blob = append(blob, TokenStringVersion)
	blob = append(blob, salt...)
	blob = append(blob, nonce...)
	blob = append(blob, ciphertext...)
	// Append a SHA-256 checksum of everything so far for integrity
	// detection on import.
	checksum := sha256.Sum256(blob)
	blob = append(blob, checksum[:]...)
	encoded := base64.RawURLEncoding.EncodeToString(blob)
	return TokenStringPrefix + encoded, nil
}

// ImportTokenString decodes a token string produced by ExportTokenString
// and returns the payload map. The password must match the one used at
// export time.
func ImportTokenString(tokenString, password string) (map[string]any, error) {
	if !strings.HasPrefix(tokenString, TokenStringPrefix) {
		return nil, fmt.Errorf("bot: invalid token string (missing %q prefix)", TokenStringPrefix)
	}
	encoded := strings.TrimPrefix(tokenString, TokenStringPrefix)
	blob, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("bot: decode token string: %w", err)
	}
	if len(blob) < 1+16+12+32+32 {
		return nil, errors.New("bot: token string too short")
	}
	// Verify checksum (last 32 bytes).
	storedChecksum := blob[len(blob)-32:]
	payload := blob[:len(blob)-32]
	actualChecksum := sha256.Sum256(payload)
	if !bytes.Equal(storedChecksum, actualChecksum[:]) {
		return nil, errors.New("bot: token string checksum mismatch (corrupted)")
	}
	version := payload[0]
	if version != TokenStringVersion {
		return nil, fmt.Errorf("bot: unsupported token string version %d", version)
	}
	salt := payload[1 : 1+16]
	nonce := payload[1+16 : 1+16+12]
	ciphertext := payload[1+16+12:]
	key := deriveTokenKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("bot: decrypt token string (wrong password?): %w", err)
	}
	// zlib-decompress with a decompression-bomb guard: a valid token
	// payload is tiny (~120 bytes), so cap the reader well below the
	// 64 KiB ceiling and refuse anything that tries to exceed it.
	zr, err := zlib.NewReader(bytes.NewReader(plaintext))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	decompressed, err := io.ReadAll(io.LimitReader(zr, maxTokenDecompressed+1))
	if err != nil {
		return nil, err
	}
	if len(decompressed) > maxTokenDecompressed {
		return nil, fmt.Errorf("bot: token string decompressed payload exceeds %d bytes (possible bomb)", maxTokenDecompressed)
	}
	var result map[string]any
	if err := json.Unmarshal(decompressed, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// deriveTokenKey derives a 32-byte AES key from a password and salt
// using PBKDF2-SHA256 with 200,000 iterations. When the password is
// empty, a fixed all-zero key is used (obfuscation only — NOT secure).
func deriveTokenKey(password string, salt []byte) []byte {
	if password == "" {
		// Obfuscation-only mode: fixed key. This is NOT encryption.
		return make([]byte, 32)
	}
	return pbkdf2.Key([]byte(password), salt, 200000, 32, sha256.New)
}

// NewFromTokenString constructs a new Bot from an encrypted token string.
// It decodes the token string, extracts the token, and calls New with the
// saved options (parse_mode, user_agent, etc.).
func NewFromTokenString(tokenString, password string, opts ...Options) (*Bot, error) {
	payload, err := ImportTokenString(tokenString, password)
	if err != nil {
		return nil, err
	}
	token, _ := payload["token"].(string)
	if token == "" {
		return nil, errors.New("bot: token string has no token field")
	}
	// Merge saved options with caller-provided opts (caller wins).
	merged := DefaultOptions()
	if pm, ok := payload["parse_mode"].(string); ok && pm != "" {
		merged.ParseMode = ParseMode(pm)
	}
	if ua, ok := payload["user_agent"].(string); ok && ua != "" {
		merged.UserAgent = ua
	}
	for _, o := range opts {
		if o.BaseURL != "" {
			merged.BaseURL = o.BaseURL
		}
		if o.Timeout > 0 {
			merged.Timeout = o.Timeout
		}
		if o.PollInterval > 0 {
			merged.PollInterval = o.PollInterval
		}
		if o.MaxRetries > 0 {
			merged.MaxRetries = o.MaxRetries
		}
		if o.UserAgent != "" {
			merged.UserAgent = o.UserAgent
		}
		if o.ParseMode != "" {
			merged.ParseMode = o.ParseMode
		}
		if o.HTTPClient != nil {
			merged.HTTPClient = o.HTTPClient
		}
	}
	return New(token, merged), nil
}
