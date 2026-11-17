package storage

// Store is the contract used by the account client to load and persist
// the active user-account credentials.
type Store interface {
	// Load returns the stored record, or nil if no session has been
	// saved yet.
	Load() (*Record, error)
	// Save upserts the record.
	Save(rec Record) error
	// Close releases any underlying resources.
	Close() error
}
