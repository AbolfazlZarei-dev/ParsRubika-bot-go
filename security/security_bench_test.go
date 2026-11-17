package security

import (
	"crypto/rsa"
	"testing"
)

// BenchmarkEncrypt measures the cost of AES-256-CBC encrypting a
// 1 KB payload.
func BenchmarkEncrypt(b *testing.B) {
	key := []byte("01234567890123456789012345678901")
	plain := make([]byte, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Encrypt(plain, key)
	}
}

// BenchmarkDecrypt measures the cost of AES-256-CBC decrypting a 1 KB
// payload.
func BenchmarkDecrypt(b *testing.B) {
	key := []byte("01234567890123456789012345678901")
	plain := make([]byte, 1024)
	enc, _ := Encrypt(plain, key)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decrypt(enc, key)
	}
}

// BenchmarkEncryptJSON measures the cost of JSON-marshalling + AES
// encrypting a small dict (a typical sendMessage payload).
func BenchmarkEncryptJSON(b *testing.B) {
	key := []byte("01234567890123456789012345678901")
	payload := map[string]any{
		"client": map[string]any{"app_name": "Main", "platform": "PWA"},
		"method": "sendMessage",
		"input":  map[string]any{"object_guid": "g0GROUP", "text": "hello world"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = EncryptJSON(payload, key)
	}
}

// BenchmarkDecodeAuth measures the cost of the auth-key character
// substitution over a 32-character string.
func BenchmarkDecodeAuth(b *testing.B) {
	auth := "abcdefghijklmnopqrstuvwxyz012345"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DecodeAuth(auth)
	}
}

// BenchmarkDeriveKey measures the AES key derivation cost.
func BenchmarkDeriveKey(b *testing.B) {
	auth := "abcdefghijklmnopqrstuvwxyz012345"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DeriveKey(auth)
	}
}

// BenchmarkSignRSA measures the PKCS1v15 SHA-256 signature cost for a
// typical request payload.
func BenchmarkSignRSA(b *testing.B) {
	_, priv, _ := GenerateKeyPair(1024)
	pk := priv // keep a reference so the compiler doesn't optimize it away
	data := "eyJjbGllbnQiOnsiYXBwX25hbWUiOiJNYWluIn0sIm1ldGhvZCI6InNlbmRNZXNzYWdlIn0="
	var rsaKey *rsa.PrivateKey
	if k, err := ParsePrivateKey(pk); err == nil {
		rsaKey = k
	}
	if rsaKey == nil {
		b.Fatal("could not parse key")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SignRSA(rsaKey, data)
	}
}
