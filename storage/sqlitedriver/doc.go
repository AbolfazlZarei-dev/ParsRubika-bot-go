// Package sqlitedriver registers a pure-Go (no cgo) SQLite driver under
// the name "sqlite" so storage.NewSQLiteStore works out of the box.
//
// It lives in its own module so that the core library keeps a minimal
// dependency set and builds with Go 1.20; this driver needs Go 1.21+.
//
//	import _ "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/storage/sqlitedriver"
package sqlitedriver

import _ "modernc.org/sqlite"
