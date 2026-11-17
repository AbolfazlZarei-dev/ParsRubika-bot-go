# `transport` — Pool

_Source: [`pool.go`](../../ParsRubika/transport/pool.go)_

## Symbols

### `AcquireBuffer`

```go
func AcquireBuffer() *bytes.Buffer
```

AcquireBuffer returns a reset *bytes.Buffer from the pool. The caller
must return it with ReleaseBuffer.

### `ReleaseBuffer`

```go
func ReleaseBuffer(b *bytes.Buffer)
```

ReleaseBuffer returns a buffer to the pool. The buffer is reset before
being put back so its memory can be reused.

### `AcquireMap`

```go
func AcquireMap() map[string]any
```

AcquireMap returns a cleared map[string]any from the pool.

### `ReleaseMap`

```go
func ReleaseMap(m map[string]any)
```

ReleaseMap returns a map to the pool. The map is cleared first so its
memory can be reused.
