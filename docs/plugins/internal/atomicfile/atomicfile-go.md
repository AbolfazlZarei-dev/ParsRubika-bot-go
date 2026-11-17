# `atomicfile` — Package atomicfile writes files so that a crash never leaves a

_Source: [`atomicfile.go`](../../../../ParsRubika/plugins/internal/atomicfile/atomicfile.go)_

## Symbols

### `WriteJSON`

```go
func WriteJSON(path string, v any) error
```

WriteJSON marshals v and atomically replaces path (mode 0600).

### `ReadJSON`

```go
func ReadJSON(path string, v any) error
```

ReadJSON loads path into v. A missing file is not an error.
