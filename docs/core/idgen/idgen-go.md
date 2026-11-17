# `idgen` — Package idgen provides ID-generation helpers. The Rubika API

_Source: [`idgen.go`](../../../ParsRubika/core/idgen/idgen.go)_

## Symbols

### `Secret`

```go
func Secret(length int) (string, error)
```

Secret returns a cryptographically random lowercase string of
the requested length.

### `Rnd`

```go
func Rnd() int64 { return mathrand.Int63n(1_000_000) + 1 }
```

Rnd returns a random int64 in the range [1, 1_000_000] — used as
the `rnd` field in message-sending payloads.

### `DeviceHash`

```go
func DeviceHash() string
```

DeviceHash returns a random device-hash string (26 digits).

### `RandomBytes`

```go
func RandomBytes(length int) ([]byte, error)
```

RandomBytes returns cryptographically random bytes of the
requested length.

### `HexID`

```go
func HexID(byteLen int) string
```

HexID returns a random hex string of the given byte length (so
HexID(8) returns a 16-char hex string). If crypto/rand fails the
function panics — a predictable all-zero ID would be a security
hazard and silently returning it is worse than crashing.
