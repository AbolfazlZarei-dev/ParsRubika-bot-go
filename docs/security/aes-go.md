# `security` — Aes

_Source: [`aes.go`](../../ParsRubika/security/aes.go)_

## Symbols

### `Encrypt`

```go
func Encrypt(data, key []byte) (string, error)
```

Encrypt AES-256-CBC-encrypts `data` with the supplied key (32 ASCII
chars → AES-256). The output is base64 standard-encoded so it can be
embedded directly in a JSON payload.

### `EncryptJSON`

```go
func EncryptJSON(v any, key []byte) (string, error)
```

EncryptJSON is a convenience helper that JSON-encodes `v` before
calling Encrypt.

### `Decrypt`

```go
func Decrypt(data string, key []byte) ([]byte, error)
```

Decrypt reverses Encrypt. It tries URL-safe base64 first (the wire
format), then standard base64, then URL-safe without padding, so it
is resilient to small variations in the server response.
