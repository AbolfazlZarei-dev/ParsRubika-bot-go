package storage

import "sync"

// MemoryStore is an in-memory Store. It is the default when no other
// storage is configured. Records are lost when the process exits.
type MemoryStore struct {
	mu  sync.RWMutex
	rec *Record
}

// NewMemoryStore returns a fresh MemoryStore.
func NewMemoryStore() *MemoryStore { return &MemoryStore{} }

// Load returns the stored record (or nil if empty).
func (m *MemoryStore) Load() (*Record, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.rec == nil {
		return nil, nil
	}
	rec := *m.rec
	return &rec, nil
}

// Save upserts the record (a copy is kept so callers can mutate the
// original safely).
func (m *MemoryStore) Save(rec Record) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c := rec
	m.rec = &c
	return nil
}

// Close is a no-op.
func (m *MemoryStore) Close() error { return nil }
