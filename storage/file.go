package storage

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// FileSuffix is the file extension used for JSON file sessions.
const FileSuffix = ".bin"

// FileStore persists the session record as a base64-encoded JSON array
// inside a small file. It is the lightweight alternative to SQLite.
//
// Security: when a passphrase is configured (NewFileStoreWithPassphrase
// or (*FileStore).SetPassphrase), the two secret fields — Auth and
// PrivateKey — are sealed with AES-256-GCM (PBKDF2-SHA256, 200k iters)
// before they reach disk. Phone, GUID and UserAgent are not sensitive
// and stay as plain text so the file stays human-inspectable for
// debugging. Reads auto-detect the legacy plain-text format and migrate
// transparently on the next Save.
type FileStore struct {
	path       string
	passphrase string
	// warnedOnce prevents the "writing plain-text secrets" notice from
	// flooding the log on every Save when no passphrase is set.
	warnedOnce bool
	mu         sync.Mutex
}

// NewFileStore opens (or creates) a JSON file session. If `name` does
// not end in `.bin` the suffix is appended.
//
// The returned store writes secret fields in plain text. To enable
// at-rest encryption use NewFileStoreWithPassphrase or call
// (*FileStore).SetPassphrase before the first Save.
func NewFileStore(name string) (*FileStore, error) {
	if name == "" {
		return nil, fmt.Errorf("storage: empty session name")
	}
	if filepath.Ext(name) != FileSuffix {
		name = name + FileSuffix
	}
	return &FileStore{path: name}, nil
}

// NewFileStoreWithPassphrase is like NewFileStore but enables AES-256-GCM
// at-rest encryption for the Auth and PrivateKey fields. Pass `passphrase`
// as a high-entropy secret that is NOT stored alongside the session file.
func NewFileStoreWithPassphrase(name, passphrase string) (*FileStore, error) {
	s, err := NewFileStore(name)
	if err != nil {
		return nil, err
	}
	s.passphrase = passphrase
	return s, nil
}

// SetPassphrase enables (or changes) at-rest encryption. Must be called
// before the first Save for the new key to take effect. Load must be
// called with the same passphrase that was used to Save.
func (s *FileStore) SetPassphrase(passphrase string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.passphrase = passphrase
}

// Load reads the record from disk. Returns nil (no error) when the
// file does not exist yet. Secret fields are decrypted when a passphrase
// is configured; legacy plain-text fields are returned as-is.
func (s *FileStore) Load() (*Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("storage: read file: %w", err)
	}
	if len(data) > maxFileStoreBytes {
		return nil, fmt.Errorf("storage: file store too large (%d > %d)", len(data), maxFileStoreBytes)
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(data)))
	if err != nil {
		return nil, fmt.Errorf("storage: base64 decode: %w", err)
	}
	var arr []string
	if err := json.Unmarshal(decoded, &arr); err != nil {
		return nil, fmt.Errorf("storage: json decode: %w", err)
	}
	rec := &Record{}
	if len(arr) >= 1 {
		rec.Phone = arr[0]
	}
	if len(arr) >= 2 {
		auth, _, err := OpenSecret(s.passphrase, arr[1])
		if err != nil {
			return nil, fmt.Errorf("storage: auth field: %w", err)
		}
		rec.Auth = auth
	}
	if len(arr) >= 3 {
		rec.GUID = arr[2]
	}
	if len(arr) >= 4 {
		rec.UserAgent = arr[3]
	}
	if len(arr) >= 5 {
		pk, _, err := OpenSecret(s.passphrase, arr[4])
		if err != nil {
			return nil, fmt.Errorf("storage: private_key field: %w", err)
		}
		rec.PrivateKey = pk
	}
	return rec, nil
}

// Save writes the record to disk (base64-encoded JSON array). The
// write is atomic: we write to a temp file first, then rename it over
// the target so a crash mid-write cannot corrupt the existing store.
//
// When a passphrase is configured the Auth and PrivateKey fields are
// sealed with AES-256-GCM. When no passphrase is configured the secret
// fields are written in plain text (a one-time warning is emitted).
func (s *FileStore) Save(rec Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	encAuth, err := SealSecret(s.passphrase, rec.Auth)
	if err != nil {
		return fmt.Errorf("storage: seal auth: %w", err)
	}
	encKey, err := SealSecret(s.passphrase, rec.PrivateKey)
	if err != nil {
		return fmt.Errorf("storage: seal private_key: %w", err)
	}
	if s.passphrase == "" && !s.warnedOnce {
		s.warnedOnce = true
		log.Printf("[storage] WARNING: FileStore(%q) is writing the auth key and RSA private key in PLAIN TEXT. "+
			"Use NewFileStoreWithPassphrase / SetPassphrase to enable AES-256-GCM at-rest encryption.", s.path)
	}
	arr := []string{rec.Phone, encAuth, rec.GUID, rec.UserAgent, encKey}
	raw, err := json.Marshal(arr)
	if err != nil {
		return fmt.Errorf("storage: json encode: %w", err)
	}
	encoded := base64.StdEncoding.EncodeToString(raw)
	// Atomic + durable write: a random temp file (no symlink/TOCTOU race),
	// fsync, chmod 0600, rename over the target, then fsync the parent
	// directory. Mode 0600 keeps the file readable only by the owning user
	// (auth + RSA private key live inside, even when encrypted).
	if err := writeFileAtomic(s.path, []byte(encoded), 0o600); err != nil {
		return err
	}
	return nil
}

// Close is a no-op for file stores.
func (s *FileStore) Close() error { return nil }

// Path returns the on-disk file path.
func (s *FileStore) Path() string { return s.path }
