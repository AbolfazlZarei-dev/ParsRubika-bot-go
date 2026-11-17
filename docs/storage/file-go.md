# `storage` — File

_Source: [`file.go`](../../ParsRubika/storage/file.go)_

## Symbols

### `FileSuffix`

```go
const FileSuffix = ".bin"
```

FileSuffix is the file extension used for JSON file sessions.

### `FileStore`

```go
type FileStore struct
```

FileStore persists the session record as a base64-encoded JSON array
inside a small file. It is the lightweight alternative to SQLite.

Security: when a passphrase is configured (NewFileStoreWithPassphrase
or (*FileStore).SetPassphrase), the two secret fields — Auth and
PrivateKey — are sealed with AES-256-GCM (PBKDF2-SHA256, 200k iters)
before they reach disk. Phone, GUID and UserAgent are not sensitive
and stay as plain text so the file stays human-inspectable for
debugging. Reads auto-detect the legacy plain-text format and migrate
transparently on the next Save.

### `NewFileStore`

```go
func NewFileStore(name string) (*FileStore, error)
```

NewFileStore opens (or creates) a JSON file session. If `name` does
not end in `.bin` the suffix is appended.

The returned store writes secret fields in plain text. To enable
at-rest encryption use NewFileStoreWithPassphrase or call
(*FileStore).SetPassphrase before the first Save.

### `NewFileStoreWithPassphrase`

```go
func NewFileStoreWithPassphrase(name, passphrase string) (*FileStore, error)
```

NewFileStoreWithPassphrase is like NewFileStore but enables AES-256-GCM
at-rest encryption for the Auth and PrivateKey fields. Pass `passphrase`
as a high-entropy secret that is NOT stored alongside the session file.

### `SetPassphrase`

```go
func (s *FileStore) SetPassphrase(passphrase string)
```

SetPassphrase enables (or changes) at-rest encryption. Must be called
before the first Save for the new key to take effect. Load must be
called with the same passphrase that was used to Save.

### `Load`

```go
func (s *FileStore) Load() (*Record, error)
```

Load reads the record from disk. Returns nil (no error) when the
file does not exist yet. Secret fields are decrypted when a passphrase
is configured; legacy plain-text fields are returned as-is.

### `Save`

```go
func (s *FileStore) Save(rec Record) error
```

Save writes the record to disk (base64-encoded JSON array). The
write is atomic: we write to a temp file first, then rename it over
the target so a crash mid-write cannot corrupt the existing store.

When a passphrase is configured the Auth and PrivateKey fields are
sealed with AES-256-GCM. When no passphrase is configured the secret
fields are written in plain text (a one-time warning is emitted).

### `Close`

```go
func (s *FileStore) Close() error { return nil }
```

Close is a no-op for file stores.

### `Path`

```go
func (s *FileStore) Path() string { return s.path }
```

Path returns the on-disk file path.
