package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// StringSessionVersion is the magic byte prefix used to identify a
// ParsRubika StringSession blob. The first byte of the encoded string
// is always this value, so callers can quickly tell whether a string
// is encrypted or plain.
const StringSessionVersion byte = 1

// StringSession is a portable, copy-pasteable session that holds the
// same Record tuple as the file/sqlite backends. The encoded form is
// `v1.<base64>` where the base64 decodes to:
//
// [1-byte version][1-byte flags][N bytes payload]
//
// Flags:
//
// bit 0 (0x01): payload is encrypted with AES-256-GCM
// bit 1 (0x02): payload is zlib-compressed (reserved, not yet used)
//
// When encrypted, the payload layout is:
//
// [16-byte salt][12-byte nonce][ciphertext + 16-byte GCM tag]
//
// (Go's gcm.Seal appends the authentication tag to the ciphertext.)
// The AES-256 key is derived from the user passphrase via PBKDF2-SHA256
// with 200,000 iterations. When no passphrase is supplied, the payload
// is stored as plain JSON.
type StringSession struct {
	rec *Record
}

// NewStringSession returns an empty StringSession.
func NewStringSession() *StringSession { return &StringSession{} }

// LoadString decodes an encoded StringSession. When `passphrase` is
// non-empty it is used to decrypt the payload.
func LoadString(encoded, passphrase string) (*StringSession, error) {
	if encoded == "" {
		return &StringSession{}, nil
	}
	// Strip the optional "v1." prefix.
	body := encoded
	if strings.HasPrefix(body, "v1.") {
		body = body[3:]
	}
	// Cap the encoded input before base64 decode so a maliciously huge
	// string cannot make LoadString allocate unbounded memory. The real
	// payload is a few hundred bytes; 256 KiB is a generous ceiling.
	if len(body) > maxStringSessionEncoded {
		return nil, fmt.Errorf("storage: string session too large (%d > %d)", len(body), maxStringSessionEncoded)
	}
	raw, err := base64.RawStdEncoding.DecodeString(body)
	if err != nil {
		raw, err = base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, fmt.Errorf("storage: invalid string session: %w", err)
		}
	}
	if len(raw) < 2 {
		return nil, fmt.Errorf("storage: string session too short")
	}
	if raw[0] != StringSessionVersion {
		return nil, fmt.Errorf("storage: unsupported string session version %d", raw[0])
	}
	flags := raw[1]
	payload := raw[2:]

	if flags&0x01 != 0 {
		if passphrase == "" {
			return nil, fmt.Errorf("storage: string session is encrypted but no passphrase was supplied")
		}
		if len(payload) > maxStringSessionPayload {
			return nil, fmt.Errorf("storage: string session ciphertext too large (%d > %d)", len(payload), maxStringSessionPayload)
		}
		decrypted, err := decryptBlob(payload, []byte(passphrase))
		if err != nil {
			return nil, err
		}
		if len(decrypted) > maxStringSessionPayload {
			return nil, fmt.Errorf("storage: string session payload too large (%d > %d)", len(decrypted), maxStringSessionPayload)
		}
		payload = decrypted
	}

	var arr []string
	if err := json.Unmarshal(payload, &arr); err != nil {
		return nil, fmt.Errorf("storage: string session json: %w", err)
	}
	rec := &Record{}
	if len(arr) >= 1 {
		rec.Phone = arr[0]
	}
	if len(arr) >= 2 {
		rec.Auth = arr[1]
	}
	if len(arr) >= 3 {
		rec.GUID = arr[2]
	}
	if len(arr) >= 4 {
		rec.UserAgent = arr[3]
	}
	if len(arr) >= 5 {
		rec.PrivateKey = arr[4]
	}
	return &StringSession{rec: rec}, nil
}

// Load returns the stored record (or nil).
func (s *StringSession) Load() (*Record, error) {
	if s.rec == nil {
		return nil, nil
	}
	c := *s.rec
	return &c, nil
}

// Save upserts the record.
func (s *StringSession) Save(rec Record) error {
	s.rec = &rec
	return nil
}

// Close is a no-op.
func (s *StringSession) Close() error { return nil }

// Encode returns the encoded string. When `passphrase` is non-empty
// the payload is encrypted with AES-256-GCM (PBKDF2-SHA256, 200k iters).
func (s *StringSession) Encode(passphrase string) (string, error) {
	if s.rec == nil {
		return "", nil
	}
	arr := []string{s.rec.Phone, s.rec.Auth, s.rec.GUID, s.rec.UserAgent, s.rec.PrivateKey}
	raw, err := json.Marshal(arr)
	if err != nil {
		return "", fmt.Errorf("storage: encode json: %w", err)
	}
	var flags byte
	payload := raw
	if passphrase != "" {
		encrypted, err := encryptBlob(raw, []byte(passphrase))
		if err != nil {
			return "", err
		}
		payload = encrypted
		flags |= 0x01
	}
	// Build the final blob: [version][flags][payload].
	blob := make([]byte, 2+len(payload))
	blob[0] = StringSessionVersion
	blob[1] = flags
	copy(blob[2:], payload)
	return "v1." + base64.RawStdEncoding.EncodeToString(blob), nil
}

// FromSQLite copies the record stored in `other` into a fresh
// StringSession.
func FromSQLite(other Store) (*StringSession, error) {
	rec, err := other.Load()
	if err != nil {
		return nil, err
	}
	return &StringSession{rec: rec}, nil
}

// ToSQLite writes the current record into a SQLiteStore with the
// supplied file name.
//
// The target store is created WITHOUT a passphrase, so the Auth and
// PrivateKey fields are written in plain text (SQLiteStore.Save emits a
// one-time warning when this happens). To preserve at-rest encryption —
// for example when the source StringSession was loaded from an encrypted
// blob — use ToSQLiteWithPassphrase instead.
func (s *StringSession) ToSQLite(name string) (*SQLiteStore, error) {
	return s.ToSQLiteWithPassphrase(name, "")
}

// ToSQLiteWithPassphrase writes the current record into an encrypted
// SQLiteStore. Pass the same passphrase that protects (or should
// protect) the session; a non-empty value enables AES-256-GCM at-rest
// encryption for the Auth and PrivateKey columns.
func (s *StringSession) ToSQLiteWithPassphrase(name, passphrase string) (*SQLiteStore, error) {
	store, err := NewSQLiteStoreWithPassphrase(name, passphrase)
	if err != nil {
		return nil, err
	}
	if s.rec != nil {
		if err := store.Save(*s.rec); err != nil {
			store.Close()
			return nil, err
		}
	}
	return store, nil
}

// ===== AES-GCM helpers ================================================

const (
	pbkdf2Iter   = 200_000
	pbkdf2KeyLen = 32
	saltLen      = 16
	nonceLen     = 12

	// maxStringSessionPayload caps the size of a decrypted StringSession
	// payload. The stored tuple is a tiny JSON string array (~200 bytes);
	// the cap exists purely as a guard so a maliciously crafted encrypted
	// string cannot make LoadString allocate unbounded memory.
	maxStringSessionPayload = 1 << 16 // 64 KiB

	// maxStringSessionEncoded caps the base64-encoded input accepted by
	// LoadString before it is decoded. It is comfortably larger than
	// maxStringSessionPayload's base64 expansion (~1.33x) while still
	// bounding the allocation an attacker can force.
	maxStringSessionEncoded = 1 << 18 // 256 KiB
)

func deriveKey(passphrase, salt []byte) []byte {
	return pbkdf2.Key(passphrase, salt, pbkdf2Iter, pbkdf2KeyLen, sha256.New)
}

func encryptBlob(plain, passphrase []byte) ([]byte, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	key := deriveKey(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, nonceLen)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nil, nonce, plain, nil)
	// Layout: [salt][nonce][ciphertext (with tag appended)].
	out := make([]byte, 0, saltLen+nonceLen+len(ciphertext))
	out = append(out, salt...)
	out = append(out, nonce...)
	out = append(out, ciphertext...)
	return out, nil
}

func decryptBlob(blob, passphrase []byte) ([]byte, error) {
	if len(blob) < saltLen+nonceLen+aes.BlockSize {
		return nil, fmt.Errorf("storage: encrypted blob too short")
	}
	salt := blob[:saltLen]
	nonce := blob[saltLen : saltLen+nonceLen]
	ciphertext := blob[saltLen+nonceLen:]
	key := deriveKey(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("storage: decryption failed (wrong passphrase?): %w", err)
	}
	return plain, nil
}

// Ensure bytes is imported.
