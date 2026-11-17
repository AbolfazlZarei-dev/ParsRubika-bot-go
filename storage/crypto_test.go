package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sampleRecord() Record {
	return Record{
		Phone:      "+989121234567",
		Auth:       "abcdefghijklmnopqrstuvwxyz012345",
		GUID:       "u0" + strings.Repeat("a", 30),
		UserAgent:  "ParsRubika/test",
		PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nFAKEKEYDATA\n-----END RSA PRIVATE KEY-----\n",
	}
}

func TestSealOpenRoundTrip(t *testing.T) {
	pass := "correct horse battery staple"
	for _, plain := range []string{"", "x", "auth-key-value", strings.Repeat("k", 4096)} {
		sealed, err := SealSecret(pass, plain)
		if err != nil {
			t.Fatalf("SealSecret(%q): %v", plain, err)
		}
		if pass != "" && plain != "" && sealed == plain {
			t.Fatalf("SealSecret did not encrypt %q", plain)
		}
		got, wasEnc, err := OpenSecret(pass, sealed)
		if err != nil {
			t.Fatalf("OpenSecret(%q): %v", sealed, err)
		}
		if got != plain {
			t.Fatalf("round-trip mismatch: got %q want %q", got, plain)
		}
		if plain != "" && !wasEnc {
			t.Fatalf("expected encrypted envelope for %q", plain)
		}
	}
}

func TestSealEmptyPassphraseIsPlaintext(t *testing.T) {
	sealed, err := SealSecret("", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if sealed != "secret" {
		t.Fatalf("empty passphrase must pass through, got %q", sealed)
	}
	got, wasEnc, err := OpenSecret("", "legacy-plaintext-auth")
	if err != nil {
		t.Fatal(err)
	}
	if got != "legacy-plaintext-auth" || wasEnc {
		t.Fatalf("legacy plaintext must pass through, got %q wasEnc=%v", got, wasEnc)
	}
}

func TestOpenSealedWithoutPassphraseFails(t *testing.T) {
	sealed, _ := SealSecret("pw", "secret")
	if _, _, err := OpenSecret("", sealed); err == nil {
		t.Fatal("OpenSecret on sealed field with no passphrase must error")
	}
}

func TestOpenWithWrongPassphraseFails(t *testing.T) {
	sealed, _ := SealSecret("correct", "secret")
	if _, _, err := OpenSecret("wrong", sealed); err == nil {
		t.Fatal("OpenSecret with wrong passphrase must fail (GCM auth)")
	}
}

func TestSealRejectsOversized(t *testing.T) {
	big := strings.Repeat("x", secretMaxPlain+1)
	if _, err := SealSecret("pw", big); err == nil {
		t.Fatal("SealSecret must reject oversized plaintext")
	}
}

func TestFileStoreEncryptedRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "enc")
	s, err := NewFileStoreWithPassphrase(path, "super-secret-pass")
	if err != nil {
		t.Fatal(err)
	}
	rec := sampleRecord()
	if err := s.Save(rec); err != nil {
		t.Fatalf("Save: %v", err)
	}
	// Read the raw file: the auth and private_key must NOT appear in clear text.
	raw, err := os.ReadFile(s.Path())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), rec.Auth) {
		t.Fatalf("auth key leaked into file in clear text:\n%s", raw)
	}
	if strings.Contains(string(raw), "BEGIN RSA PRIVATE KEY") {
		t.Fatalf("private key leaked into file in clear text:\n%s", raw)
	}
	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Auth != rec.Auth || got.PrivateKey != rec.PrivateKey || got.GUID != rec.GUID {
		t.Fatalf("record mismatch: %+v vs %+v", got, rec)
	}
}

func TestFileStoreLegacyMigration(t *testing.T) {
	dir := t.TempDir()
	// 1) write a legacy (plain-text) file with no passphrase.
	legacy, err := NewFileStore(filepath.Join(dir, "legacy"))
	if err != nil {
		t.Fatal(err)
	}
	rec := sampleRecord()
	if err := legacy.Save(rec); err != nil {
		t.Fatal(err)
	}
	// 2) open the SAME file with a passphrase: legacy plain text must still load.
	migrated, err := NewFileStoreWithPassphrase(filepath.Join(dir, "legacy"), "pass")
	if err != nil {
		t.Fatal(err)
	}
	got, err := migrated.Load()
	if err != nil {
		t.Fatalf("migrated Load: %v", err)
	}
	if got.Auth != rec.Auth || got.PrivateKey != rec.PrivateKey {
		t.Fatalf("migration Load mismatch: %+v", got)
	}
	// 3) re-save -> file is now encrypted on disk.
	if err := migrated.Save(rec); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(migrated.Path())
	if strings.Contains(string(raw), rec.Auth) {
		t.Fatalf("post-migration file still contains plain-text auth:\n%s", raw)
	}
	// 4) loading the now-encrypted file with the correct passphrase works.
	got2, err := migrated.Load()
	if err != nil {
		t.Fatalf("post-migration Load: %v", err)
	}
	if got2.Auth != rec.Auth {
		t.Fatalf("post-migration Auth mismatch")
	}
}
