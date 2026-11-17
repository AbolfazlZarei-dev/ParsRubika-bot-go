package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// SQLiteDriverName is the database/sql driver name used for sessions.
const SQLiteDriverName = "sqlite"

// ErrNoSQLiteDriver is returned when no "sqlite" database/sql driver
// has been linked into the program. Add a blank import of
//
//	_ "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/storage/sqlitedriver"
//
// (or any pure-Go/cgo driver registered as "sqlite"), or use the file
// or in-memory stores that need no driver.
var ErrNoSQLiteDriver = errors.New(`storage: no "sqlite" database/sql driver registered; import storage/sqlitedriver or use NewFileStore`)

// driverRegistered reports whether a database/sql driver is available.
func driverRegistered(name string) bool {
	for _, d := range sql.Drivers() {
		if d == name {
			return true
		}
	}
	return false
}

// SQLiteSuffix is the file extension used for SQLite sessions.
const SQLiteSuffix = ".session"

// SQLiteStore is a file-backed Store backed by SQLite. Existing session
// files created by earlier releases stay compatible.
//
// Security: when a passphrase is configured (NewSQLiteStoreWithPassphrase
// or (*SQLiteStore).SetPassphrase), the Auth and PrivateKey columns are
// sealed with AES-256-GCM (PBKDF2-SHA256, 200k iters) before being
// written. Reads auto-detect legacy plain-text columns and migrate
// transparently on the next Save.
type SQLiteStore struct {
	path       string
	conn       *sql.DB
	passphrase string
	warnedOnce bool
	mu         sync.Mutex
}

// NewSQLiteStore opens (or creates) a SQLite-backed session file. If
// `name` does not end in `.session` the suffix is appended automatically.
//
// The returned store writes secret columns in plain text. To enable
// at-rest encryption use NewSQLiteStoreWithPassphrase or call
// (*SQLiteStore).SetPassphrase before the first Save.
func NewSQLiteStore(name string) (*SQLiteStore, error) {
	if name == "" {
		return nil, fmt.Errorf("storage: empty session name")
	}
	if filepath.Ext(name) != SQLiteSuffix {
		name = name + SQLiteSuffix
	}
	if !driverRegistered(SQLiteDriverName) {
		return nil, ErrNoSQLiteDriver
	}
	conn, err := sql.Open(SQLiteDriverName, name)
	if err != nil {
		return nil, fmt.Errorf("storage: open sqlite: %w", err)
	}
	// Use the classic rollback-journal instead of WAL so the session file
	// has no separate -wal / -shm sidecars that would follow the process
	// umask and potentially be world-readable. The session DB is tiny, so
	// WAL's concurrency benefit is irrelevant here.
	if _, err := conn.Exec(`pragma journal_mode=DELETE`); err != nil {
		conn.Close()
		return nil, fmt.Errorf("storage: set journal mode: %w", err)
	}
	if err := initSchema(conn); err != nil {
		conn.Close()
		return nil, err
	}
	// Restrict the on-disk file (which holds the auth key + RSA private
	// key, even when encrypted) to the owning user. sql.Open honours the
	// umask, so we re-assert 0600 here. Best-effort: a newly created file
	// is chmod'd; an existing one is tightened if it was looser.
	_ = os.Chmod(name, 0o600)
	return &SQLiteStore{path: name, conn: conn}, nil
}

// NewSQLiteStoreWithPassphrase is like NewSQLiteStore but enables
// AES-256-GCM at-rest encryption for the auth and private_key columns.
func NewSQLiteStoreWithPassphrase(name, passphrase string) (*SQLiteStore, error) {
	s, err := NewSQLiteStore(name)
	if err != nil {
		return nil, err
	}
	s.passphrase = passphrase
	return s, nil
}

// SetPassphrase enables (or changes) at-rest encryption. Must be called
// before the first Save for the new key to take effect.
func (s *SQLiteStore) SetPassphrase(passphrase string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.passphrase = passphrase
}

func initSchema(conn *sql.DB) error {
	_, err := conn.Exec(`create table if not exists version (version integer primary key)`)
	if err != nil {
		return fmt.Errorf("storage: create version table: %w", err)
	}
	_, err = conn.Exec(`insert or ignore into version(version) values (1)`)
	if err != nil {
		return fmt.Errorf("storage: seed version: %w", err)
	}
	_, err = conn.Exec(`create table if not exists session (
		phone text primary key,
		auth text,
		guid text,
		agent text,
		private_key text)`)
	if err != nil {
		return fmt.Errorf("storage: create session table: %w", err)
	}
	return nil
}

// Load returns the single stored session row, or nil if empty. Secret
// columns are decrypted when a passphrase is configured; legacy
// plain-text columns are returned as-is (and re-sealed on next Save).
func (s *SQLiteStore) Load() (*Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Pick the most recently inserted row deterministically. The schema
	// allows multiple phones to share a table; a bare `limit 1` would
	// return an implementation-defined row, so order by rowid desc to
	// always select the latest session.
	row := s.conn.QueryRow(`select phone, auth, guid, agent, private_key from session order by rowid desc limit 1`)
	var phone, encAuth, guid, agent, encKey string
	err := row.Scan(&phone, &encAuth, &guid, &agent, &encKey)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("storage: scan row: %w", err)
	}
	rec := &Record{Phone: phone, GUID: guid, UserAgent: agent}
	if encAuth != "" {
		auth, _, err := OpenSecret(s.passphrase, encAuth)
		if err != nil {
			return nil, fmt.Errorf("storage: auth field: %w", err)
		}
		rec.Auth = auth
	}
	if encKey != "" {
		pk, _, err := OpenSecret(s.passphrase, encKey)
		if err != nil {
			return nil, fmt.Errorf("storage: private_key field: %w", err)
		}
		rec.PrivateKey = pk
	}
	return rec, nil
}

// Save upserts the row identified by rec.Phone (or "(default)" when empty).
// When a passphrase is configured the auth and private_key columns are
// sealed with AES-256-GCM; otherwise they are written in plain text (a
// one-time warning is emitted).
func (s *SQLiteStore) Save(rec Record) error {
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
		log.Printf("[storage] WARNING: SQLiteStore(%q) is writing the auth key and RSA private key in PLAIN TEXT. "+
			"Use NewSQLiteStoreWithPassphrase / SetPassphrase to enable AES-256-GCM at-rest encryption.", s.path)
	}
	phone := rec.Phone
	if phone == "" {
		phone = "(default)"
	}
	_, err = s.conn.Exec(`insert or replace into session (phone, auth, guid, agent, private_key) values (?, ?, ?, ?, ?)`,
		phone, encAuth, rec.GUID, rec.UserAgent, encKey)
	if err != nil {
		return fmt.Errorf("storage: insert: %w", err)
	}
	// The SQLite file is created lazily on first write, so re-assert 0600
	// here in case the open-time chmod ran before the file existed.
	_ = os.Chmod(s.path, 0o600)
	return nil
}

// Close releases the underlying SQLite connection.
func (s *SQLiteStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// Path returns the on-disk file path.
func (s *SQLiteStore) Path() string { return s.path }
