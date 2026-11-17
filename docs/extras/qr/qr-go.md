# `qr` — Package qr provides a simple QR-code generator. It produces a

_Source: [`qr.go`](../../../ParsRubika/extras/qr/qr.go)_

## Symbols

### `Level`

```go
type Level int
```

Level is the error-correction level.

### `LevelL`

```go
const LevelL Level = iota // 7% recovery
```

### `LevelM`

```go
const LevelM // 15% recovery
```

### `LevelQ`

```go
const LevelQ // 25% recovery
```

### `LevelH`

```go
const LevelH // 30% recovery
```

### `Encode`

```go
func Encode(text string, size int, level Level) ([]byte, error)
```

Encode generates a QR code PNG image and returns it as a byte
slice.

### `EncodeBase64`

```go
func EncodeBase64(text string, size int, level Level) (string, error)
```

EncodeBase64 generates a QR code and returns it as a base64
string (useful for embedding in HTML/markdown).

### `SaveFile`

```go
func SaveFile(text, path string, size int, level Level) error
```

SaveFile generates a QR code and writes it to the given file.

### `EncodeBuffer`

```go
func EncodeBuffer(text string, size int, level Level, buf *bytes.Buffer) error
```

EncodeBuffer generates a QR code and writes it to the supplied
buffer.
