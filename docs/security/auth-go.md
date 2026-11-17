# `security` — Auth

_Source: [`auth.go`](../../ParsRubika/security/auth.go)_

## Symbols

### `DecodeAuth`

```go
func DecodeAuth(auth string) string
```

DecodeAuth applies the character-substitution used to derive the
`decode_auth` payload. The mapping is:

- lowercase letters: shifted by (32 - (c - 'a')) mod 26
- uppercase letters: shifted by (29 - (c - 'A')) mod 26
- digits: shifted by (13 - (c - '0')) mod 10

All other characters (symbols, dashes, underscores) are kept as-is.
Go's `%` operator follows the sign of the dividend, so we use the
`((x % m) + m) % m` idiom to obtain an always-positive
modulo.

### `DeriveKey`

```go
func DeriveKey(auth string) (string, error)
```

DeriveKey derives the AES-256 symmetric key from a 32-character auth
string by:

1. Splitting auth into 4 chunks of 8 characters.
2. Reordering them as chunks[2] + chunks[0] + chunks[3] + chunks[1].
3. For each character: applying `((c - 'a' + 9) mod 26) + 'a'`.

The result is the 32-character AES-256 key (the underlying bytes are
the ASCII codes of these characters, exactly as the server does).
