# `storage` — Sqlite

_Source: [`sqlite.go`](../../ParsRubika/storage/sqlite.go)_

## Symbols

### `SQLiteDriverName`

```go
const SQLiteDriverName = "sqlite"
```

SQLiteDriverName is the database/sql driver name used for sessions.

### `ErrNoSQLiteDriver`

```go
var ErrNoSQLiteDriver = errors.New(`storage: no "sqlite" database/sql driver registered; import storage/sqlitedriver or use NewFileStore`)
```

ErrNoSQLiteDriver is returned when no "sqlite" database/sql driver
has been linked into the program. Add a blank import of

_ "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/storage/sqlitedriver"

(or any pure-Go/cgo driver registered as "sqlite"), or use the file
or in-memory stores that need no driver.

### `SQLiteSuffix`

```go
const SQLiteSuffix = ".session"
```

SQLiteSuffix is the file extension used for SQLite sessions.

### `SQLiteStore`

```go
type SQLiteStore struct
```

SQLiteStore is a file-backed Store backed by SQLite. Existing session
files created by earlier releases stay compatible.

Security: when a passphrase is configured (NewSQLiteStoreWithPassphrase
or (*SQLiteStore).SetPassphrase), the Auth and PrivateKey columns are
sealed with AES-256-GCM (PBKDF2-SHA256, 200k iters) before being
written. Reads auto-detect legacy plain-text columns and migrate
transparently on the next Save.

### `NewSQLiteStore`

```go
func NewSQLiteStore(name string) (*SQLiteStore, error)
```

NewSQLiteStore opens (or creates) a SQLite-backed session file. If
`name` does not end in `.session` the suffix is appended automatically.

The returned store writes secret columns in plain text. To enable
at-rest encryption use NewSQLiteStoreWithPassphrase or call
(*SQLiteStore).SetPassphrase before the first Save.

### `NewSQLiteStoreWithPassphrase`

```go
func NewSQLiteStoreWithPassphrase(name, passphrase string) (*SQLiteStore, error)
```

NewSQLiteStoreWithPassphrase is like NewSQLiteStore but enables
AES-256-GCM at-rest encryption for the auth and private_key columns.

### `SetPassphrase`

```go
func (s *SQLiteStore) SetPassphrase(passphrase string)
```

SetPassphrase enables (or changes) at-rest encryption. Must be called
before the first Save for the new key to take effect.

### `Load`

```go
func (s *SQLiteStore) Load() (*Record, error)
```

Load returns the single stored session row, or nil if empty. Secret
columns are decrypted when a passphrase is configured; legacy
plain-text columns are returned as-is (and re-sealed on next Save).

### `Save`

```go
func (s *SQLiteStore) Save(rec Record) error
```

Save upserts the row identified by rec.Phone (or "(default)" when empty).
When a passphrase is configured the auth and private_key columns are
sealed with AES-256-GCM; otherwise they are written in plain text (a
one-time warning is emitted).

### `Close`

```go
func (s *SQLiteStore) Close() error
```

Close releases the underlying SQLite connection.

### `Path`

```go
func (s *SQLiteStore) Path() string { return s.path }
```

Path returns the on-disk file path.
