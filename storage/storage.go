// Package storage defines the pluggable session storage used by the
// ParsRubika account client. Three implementations are provided:
//
// - Memory : in-memory (lost when the process exits)
// - SQLite : file-backed, schema-compatible with the original `.rp`
// files so existing sessions can be re-used.
// - File : a JSON file encrypted with the user's auth key.
//
// All three implement the Store interface so callers can swap them
// freely or supply their own implementation (e.g. Redis, Vault, ...).
package storage

// Record holds the persistent identity of a logged-in user.
type Record struct {
	Phone      string
	Auth       string
	GUID       string
	UserAgent  string
	PrivateKey string
}
