# `storage` — Crypto

_Source: [`crypto.go`](../../ParsRubika/storage/crypto.go)_

## Symbols

### `ErrNoPassphrase`

```go
var ErrNoPassphrase = errors.New("storage: session field is encrypted but no passphrase was supplied")
```

ErrNoPassphrase is returned when a sealed field is encountered on read
but no passphrase was configured.

### `SealSecret`

```go
func SealSecret(passphrase, plain string) (string, error)
```

SealSecret encrypts `plain` with a key derived from `passphrase` and
returns the "PRSE1.<base64>" envelope. When `passphrase` is empty the
plain text is returned unchanged (legacy / opt-out mode).

### `OpenSecret`

```go
func OpenSecret(passphrase, encoded string) (string, bool, error)
```

OpenSecret reverses SealSecret. A field without the "PRSE1." prefix is
treated as legacy plain text and returned verbatim (second return is
false, meaning "was not encrypted"). This keeps reads backward
compatible with sessions written by older releases.
