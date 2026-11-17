package storage

import (
	"strings"
	"testing"
)

// BenchmarkSealSecret measures the cost of sealing a small secret field
// (an auth key) with AES-256-GCM + PBKDF2 (200k iters). PBKDF2 dominates.
// This is the per-Save cost for an encrypted FileStore/SQLiteStore.
func BenchmarkSealSecret(b *testing.B) {
	pass := "benchmark-passphrase"
	plain := "abcdefghijklmnopqrstuvwxyz012345" // 32-char auth key
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := SealSecret(pass, plain); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkOpenSecret measures the decrypt path (also PBKDF2-bound).
func BenchmarkOpenSecret(b *testing.B) {
	pass := "benchmark-passphrase"
	plain := "abcdefghijklmnopqrstuvwxyz012345"
	sealed, _ := SealSecret(pass, plain)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := OpenSecret(pass, sealed); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFileStoreSaveEncrypted covers the full Save round-trip with
// at-rest encryption enabled (seal two fields + base64 + atomic write).
func BenchmarkFileStoreSaveEncrypted(b *testing.B) {
	dir := b.TempDir()
	s, _ := NewFileStoreWithPassphrase(dir+"/bench", "bench-pass")
	rec := Record{
		Phone:      "+989121234567",
		Auth:       "abcdefghijklmnopqrstuvwxyz012345",
		GUID:       "u0" + strings.Repeat("a", 30),
		UserAgent:  "ParsRubika/bench",
		PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nDATA\n-----END RSA PRIVATE KEY-----\n",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := s.Save(rec); err != nil {
			b.Fatal(err)
		}
	}
}
