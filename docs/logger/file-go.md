# `logger` — File

_Source: [`file.go`](../../ParsRubika/logger/file.go)_

## Symbols

### `RotatingFile`

```go
type RotatingFile struct
```

RotatingFile is an io.WriteCloser that appends to a file and rotates
it once it grows past a size limit, keeping a bounded number of
numbered backups (app.log.1, app.log.2, ...).

### `NewRotatingFile`

```go
func NewRotatingFile(path string, maxSizeMB, maxBackups int) (*RotatingFile, error)
```

NewRotatingFile opens path for appending, creating directories as
needed. maxSizeMB <= 0 disables rotation.

### `Write`

```go
func (r *RotatingFile) Write(p []byte) (int, error)
```

Write implements io.Writer.

### `Close`

```go
func (r *RotatingFile) Close() error
```

Close closes the underlying file.

### `Rotate`

```go
func (r *RotatingFile) Rotate() error
```

Rotate forces an immediate rotation.
