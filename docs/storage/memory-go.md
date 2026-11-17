# `storage` — Memory

_Source: [`memory.go`](../../ParsRubika/storage/memory.go)_

## Symbols

### `MemoryStore`

```go
type MemoryStore struct
```

MemoryStore is an in-memory Store. It is the default when no other
storage is configured. Records are lost when the process exits.

### `NewMemoryStore`

```go
func NewMemoryStore() *MemoryStore { return &MemoryStore{} }
```

NewMemoryStore returns a fresh MemoryStore.

### `Load`

```go
func (m *MemoryStore) Load() (*Record, error)
```

Load returns the stored record (or nil if empty).

### `Save`

```go
func (m *MemoryStore) Save(rec Record) error
```

Save upserts the record (a copy is kept so callers can mutate the
original safely).

### `Close`

```go
func (m *MemoryStore) Close() error { return nil }
```

Close is a no-op.
