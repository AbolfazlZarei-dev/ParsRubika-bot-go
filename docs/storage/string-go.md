# `storage` — String

_Source: [`string.go`](../../ParsRubika/storage/string.go)_

## Symbols

### `StringSession`

```go
type StringSession struct
```

StringSession is a portable, copy-pasteable session that holds the
same Record tuple as the file/sqlite backends. The encoded form is
`v1.<base64>` where the base64 decodes to:

[1-byte version][1-byte flags][N bytes payload]

Flags:

bit 0 (0x01): payload is encrypted with AES-256-GCM
bit 1 (0x02): payload is zlib-compressed (reserved, not yet used)

When encrypted, the payload layout is:

[16-byte salt][12-byte nonce][ciphertext + 16-byte GCM tag]

(Go's gcm.Seal appends the authentication tag to the ciphertext.)
The AES-256 key is derived from the user passphrase via PBKDF2-SHA256
with 200,000 iterations. When no passphrase is supplied, the payload
is stored as plain JSON.

### `NewStringSession`

```go
func NewStringSession() *StringSession { return &StringSession{} }
```

NewStringSession returns an empty StringSession.

### `LoadString`

```go
func LoadString(encoded, passphrase string) (*StringSession, error)
```

LoadString decodes an encoded StringSession. When `passphrase` is
non-empty it is used to decrypt the payload.

### `Load`

```go
func (s *StringSession) Load() (*Record, error)
```

Load returns the stored record (or nil).

### `Save`

```go
func (s *StringSession) Save(rec Record) error
```

Save upserts the record.

### `Close`

```go
func (s *StringSession) Close() error { return nil }
```

Close is a no-op.

### `Encode`

```go
func (s *StringSession) Encode(passphrase string) (string, error)
```

Encode returns the encoded string. When `passphrase` is non-empty
the payload is encrypted with AES-256-GCM (PBKDF2-SHA256, 200k iters).

### `FromSQLite`

```go
func FromSQLite(other Store) (*StringSession, error)
```

FromSQLite copies the record stored in `other` into a fresh
StringSession.

### `ToSQLite`

```go
func (s *StringSession) ToSQLite(name string) (*SQLiteStore, error)
```

ToSQLite writes the current record into a SQLiteStore with the
supplied file name.

The target store is created WITHOUT a passphrase, so the Auth and
PrivateKey fields are written in plain text (SQLiteStore.Save emits a
one-time warning when this happens). To preserve at-rest encryption —
for example when the source StringSession was loaded from an encrypted
blob — use ToSQLiteWithPassphrase instead.

### `ToSQLiteWithPassphrase`

```go
func (s *StringSession) ToSQLiteWithPassphrase(name, passphrase string) (*SQLiteStore, error)
```

ToSQLiteWithPassphrase writes the current record into an encrypted
SQLiteStore. Pass the same passphrase that protects (or should
protect) the session; a non-empty value enables AES-256-GCM at-rest
encryption for the Auth and PrivateKey columns.
