package sqlitedriver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/storage"
)

// This test file lives in the sqlitedriver module so the pure-Go
// "sqlite" driver is linked into the test binary. It exercises the
// storage.SQLiteStore behaviour that the core storage package can only
// test with `t.Skip()` because it ships without a driver:
//
//   - file permission 0600,
//   - deterministic latest-row selection (ORDER BY rowid DESC),
//   - AES-256-GCM encrypted round-trip + wrong-passphrase failure,
//   - StringSession.ToSQLiteWithPassphrase preserving encryption.

func tmpDB(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.session")
}

func TestSQLiteDriverLinked(t *testing.T) {
	// Sanity: the driver really is registered by this package's doc.go.
	// NewSQLiteStore returns ErrNoSQLiteDriver when no driver is linked,
	// so a successful open on a temp path proves registration.
	store, err := storage.NewSQLiteStore(tmpDB(t))
	if err != nil {
		t.Fatalf("sqlite driver not registered; NewSQLiteStore: %v", err)
	}
	store.Close()
}

func TestSQLiteRoundTripPlain(t *testing.T) {
	store, err := storage.NewSQLiteStore(tmpDB(t))
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer store.Close()

	rec := storage.Record{Phone: "+981", Auth: "auth-secret-1", GUID: "g1", UserAgent: "ua"}
	if err := store.Save(rec); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got == nil || got.Auth != "auth-secret-1" || got.Phone != "+981" {
		t.Fatalf("round-trip mismatch: %+v", got)
	}
}

func TestSQLiteFileMode0600(t *testing.T) {
	path := tmpDB(t)
	store, err := storage.NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	rec := storage.Record{Phone: "+982", Auth: "secret", GUID: "g"}
	if err := store.Save(rec); err != nil {
		t.Fatalf("Save: %v", err)
	}
	store.Close()

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	// Mask down to the permission bits (ignore file type bits).
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("sqlite file mode = %v, want 0600", fi.Mode().Perm())
	}
}

func TestSQLiteLoadReturnsLatestRow(t *testing.T) {
	store, err := storage.NewSQLiteStore(tmpDB(t))
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer store.Close()

	// Insert two rows with distinct phones; the most recent must win.
	if err := store.Save(storage.Record{Phone: "old", Auth: "old-auth", GUID: "g-old"}); err != nil {
		t.Fatalf("Save old: %v", err)
	}
	if err := store.Save(storage.Record{Phone: "new", Auth: "new-auth", GUID: "g-new"}); err != nil {
		t.Fatalf("Save new: %v", err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got == nil || got.Phone != "new" || got.Auth != "new-auth" {
		t.Fatalf("Load returned %+v, want the latest row (phone=new)", got)
	}
}

func TestSQLiteEncryptedRoundTrip(t *testing.T) {
	const pass = "correct horse battery staple"
	store, err := storage.NewSQLiteStoreWithPassphrase(tmpDB(t), pass)
	if err != nil {
		t.Fatalf("NewSQLiteStoreWithPassphrase: %v", err)
	}
	defer store.Close()

	rec := storage.Record{Phone: "+983", Auth: "super-secret-auth", PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nFAKE\n-----END RSA PRIVATE KEY-----", GUID: "g"}
	if err := store.Save(rec); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Auth != rec.Auth || got.PrivateKey != rec.PrivateKey {
		t.Fatalf("encrypted round-trip mismatch: auth=%q key=%q", got.Auth, got.PrivateKey)
	}
}

func TestSQLiteWrongPassphraseFails(t *testing.T) {
	path := tmpDB(t)
	store, err := storage.NewSQLiteStoreWithPassphrase(path, "right-pass")
	if err != nil {
		t.Fatalf("NewSQLiteStoreWithPassphrase: %v", err)
	}
	if err := store.Save(storage.Record{Phone: "+984", Auth: "secret"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	store.Close()

	store2, err := storage.NewSQLiteStoreWithPassphrase(path, "wrong-pass")
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer store2.Close()
	if _, err := store2.Load(); err == nil {
		t.Fatal("Load with wrong passphrase must fail (GCM tag mismatch)")
	}
}

func TestSQLiteEncryptedColumnNotPlain(t *testing.T) {
	// The on-disk auth column must not contain the plain secret.
	path := tmpDB(t)
	store, err := storage.NewSQLiteStoreWithPassphrase(path, "pass")
	if err != nil {
		t.Fatalf("NewSQLiteStoreWithPassphrase: %v", err)
	}
	const plainAuth = "PLAIN_AUTH_VALUE_THAT_MUST_NOT_LEAK"
	if err := store.Save(storage.Record{Phone: "+985", Auth: plainAuth}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	store.Close()

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if strings.Contains(string(raw), plainAuth) {
		t.Fatal("plain auth value found inside the sqlite file; encryption did not seal the column")
	}
}

func TestStringSessionToSQLiteWithPassphrasePreservesEncryption(t *testing.T) {
	const pass = "export-pass"
	sess := storage.NewStringSession()
	if err := sess.Save(storage.Record{Phone: "+986", Auth: "to-export-auth", PrivateKey: "k"}); err != nil {
		t.Fatalf("sess.Save: %v", err)
	}
	path := tmpDB(t)
	store, err := sess.ToSQLiteWithPassphrase(path, pass)
	if err != nil {
		t.Fatalf("ToSQLiteWithPassphrase: %v", err)
	}
	store.Close()

	// Reopen with the right passphrase and confirm the secrets survived.
	store2, err := storage.NewSQLiteStoreWithPassphrase(path, pass)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer store2.Close()
	got, err := store2.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Auth != "to-export-auth" {
		t.Fatalf("ToSQLiteWithPassphrase lost the auth: %q", got.Auth)
	}

	// And the wrong passphrase must not decrypt.
	store3, err := storage.NewSQLiteStoreWithPassphrase(path, "nope")
	if err != nil {
		t.Fatalf("reopen wrong: %v", err)
	}
	defer store3.Close()
	if _, err := store3.Load(); err == nil {
		t.Fatal("wrong passphrase must not decrypt the exported session")
	}
}
