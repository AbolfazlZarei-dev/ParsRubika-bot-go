package storage

import "testing"

// skipWithoutSQLite skips tests that need a linked "sqlite" driver; the
// core module deliberately ships without one (see sqlitedriver).
func skipWithoutSQLite(t *testing.T) {
	t.Helper()
	if !driverRegistered(SQLiteDriverName) {
		t.Skip("no sqlite driver linked; run the tests of storage/sqlitedriver")
	}
}
