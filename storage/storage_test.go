package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMemoryStoreRoundTrip(t *testing.T) {
	s := NewMemoryStore()
	if rec, _ := s.Load(); rec != nil {
		t.Fatal("expected nil on empty store")
	}
	want := Record{Auth: "auth", Phone: "98912", GUID: "u0GUID", UserAgent: "ua", PrivateKey: "pk"}
	if err := s.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got == nil || got.Auth != want.Auth || got.GUID != want.GUID {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestSQLiteStoreRoundTrip(t *testing.T) {
	skipWithoutSQLite(t)
	tmpDir, err := os.MkdirTemp("", "parsrubika-sqlite-")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	name := filepath.Join(tmpDir, "test_session")
	s, err := NewSQLiteStore(name)
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer s.Close()

	if rec, _ := s.Load(); rec != nil {
		t.Fatal("expected nil on empty store")
	}
	want := Record{Auth: "auth", Phone: "98912", GUID: "u0GUID", UserAgent: "ua", PrivateKey: "pk"}
	if err := s.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got == nil || got.Auth != want.Auth || got.GUID != want.GUID || got.PrivateKey != want.PrivateKey {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestFileStoreRoundTrip(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsrubika-file-")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	name := filepath.Join(tmpDir, "test_session")
	s, err := NewFileStore(name)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}
	defer s.Close()

	if rec, _ := s.Load(); rec != nil {
		t.Fatal("expected nil on empty store")
	}
	want := Record{Auth: "auth", Phone: "98912", GUID: "u0GUID", UserAgent: "ua", PrivateKey: "pk"}
	if err := s.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got == nil || got.Auth != want.Auth || got.GUID != want.GUID || got.PrivateKey != want.PrivateKey {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	// The file should be base64-encoded (starts with 'W' since the JSON
	// array starts with '[').
	data, err := os.ReadFile(name + FileSuffix)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.HasPrefix(string(data), "W") {
		t.Fatalf("expected base64-encoded JSON (starts with 'W'), got %s", data)
	}
}
