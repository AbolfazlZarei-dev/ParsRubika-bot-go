# `security` — Random

_Source: [`random.go`](../../ParsRubika/security/random.go)_

## Symbols

### `RandomSecret`

```go
func RandomSecret(length int) (string, error)
```

RandomSecret returns a cryptographically random lowercase string of
the requested length. Returns an empty string when length <= 0.

Each character is drawn uniformly from the alphabet using
crypto/rand.Int, which performs rejection sampling internally and so
avoids the modulo bias that a naive `byte % len(alphabet)` mapping
would introduce (with a 26-letter alphabet, 256 % 26 == 22, so the
last four letters would otherwise be ~10% under-represented).

### `RandomBytes`

```go
func RandomBytes(length int) ([]byte, error)
```

RandomBytes returns a cryptographically random byte slice of the
requested length.
