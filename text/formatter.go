// Package text implements the Markdown/HTML → metadata conversion used
// to attach formatting information to outgoing messages.
//
// The wire format requires `meta_data_parts` with UTF-16 offsets, so
// the formatter builds UTF-16 prefix lengths while scanning. The
// formatter is stateless and safe for concurrent use.
package text

// Formatter is the public, stateless converter. Use New() to obtain a
// zero-cost instance (it is just a struct with no fields).
type Formatter struct{}

// New returns a fresh Formatter. It is a zero-allocation constructor
// that exists only so callers don't have to write `text.Formatter{}`
// everywhere.
func New() Formatter { return Formatter{} }
