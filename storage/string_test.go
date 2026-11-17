package storage

import (
	"strings"
	"testing"
)

func TestStringSessionPlain(t *testing.T) {
	rec := Record{
		Phone:      "989121234567",
		Auth:       "abcdefghijklmnopqrstuvwxyz012345",
		GUID:       "u0GUID",
		UserAgent:  "ParsRubika/4.0.0",
		PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMII...\n-----END RSA PRIVATE KEY-----",
	}
	s := NewStringSession()
	if err := s.Save(rec); err != nil {
		t.Fatalf("Save: %v", err)
	}
	encoded, err := s.Encode("")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if encoded == "" {
		t.Fatal("empty encoded")
	}
	if !strings.HasPrefix(encoded, "v1.") {
		t.Fatalf("expected v1. prefix, got %s", encoded)
	}
	loaded, err := LoadString(encoded, "")
	if err != nil {
		t.Fatalf("LoadString: %v", err)
	}
	got, err := loaded.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got == nil || got.Auth != rec.Auth || got.GUID != rec.GUID || got.PrivateKey != rec.PrivateKey {
		t.Fatalf("round-trip mismatch: %+v", got)
	}
}

func TestStringSessionEncrypted(t *testing.T) {
	rec := Record{
		Phone: "989121234567",
		Auth:  "abcdefghijklmnopqrstuvwxyz012345",
		GUID:  "u0GUID",
	}
	s := NewStringSession()
	_ = s.Save(rec)
	encoded, err := s.Encode("my-secret-passphrase")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if !strings.HasPrefix(encoded, "v1.") {
		t.Fatalf("expected v1. prefix, got %s", encoded)
	}
	loaded, err := LoadString(encoded, "my-secret-passphrase")
	if err != nil {
		t.Fatalf("LoadString: %v", err)
	}
	got, _ := loaded.Load()
	if got == nil || got.Auth != rec.Auth {
		t.Fatalf("encrypted round-trip mismatch: %+v", got)
	}

	// Wrong passphrase should fail.
	if _, err := LoadString(encoded, "wrong-passphrase"); err == nil {
		t.Fatal("expected error with wrong passphrase")
	}

	// No passphrase on an encrypted session should fail.
	if _, err := LoadString(encoded, ""); err == nil {
		t.Fatal("expected error when no passphrase supplied for encrypted session")
	}
}

func TestStringSessionEmpty(t *testing.T) {
	s := NewStringSession()
	encoded, err := s.Encode("")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if encoded != "" {
		t.Fatalf("empty session should encode to empty string, got %q", encoded)
	}
}

func TestStringSessionFromSQLite(t *testing.T) {
	skipWithoutSQLite(t)
	tmp := t.TempDir()
	sqlite, err := NewSQLiteStore(tmp + "/test")
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer sqlite.Close()
	rec := Record{Auth: "auth", GUID: "u0GUID"}
	if err := sqlite.Save(rec); err != nil {
		t.Fatalf("Save: %v", err)
	}
	s, err := FromSQLite(sqlite)
	if err != nil {
		t.Fatalf("FromSQLite: %v", err)
	}
	got, _ := s.Load()
	if got == nil || got.Auth != rec.Auth || got.GUID != rec.GUID {
		t.Fatalf("FromSQLite mismatch: %+v", got)
	}
}
