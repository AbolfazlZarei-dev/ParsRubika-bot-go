package storage

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---- OpenSecret corrupted-data paths ----

func sealedEnvelope(t *testing.T, pass, plain string) []byte {
	t.Helper()
	s, err := SealSecret(pass, plain)
	if err != nil {
		t.Fatalf("SealSecret: %v", err)
	}
	raw, err := base64.RawStdEncoding.DecodeString(s[len(secretPrefix):])
	if err != nil {
		t.Fatalf("decode sealed envelope: %v", err)
	}
	return raw
}

func resealed(t *testing.T, raw []byte) string {
	t.Helper()
	return secretPrefix + base64.RawStdEncoding.EncodeToString(raw)
}

func TestOpenSecretTruncatedEnvelopeFails(t *testing.T) {
	raw := sealedEnvelope(t, "pw", "secret")
	// Keep only the version byte + a couple of salt bytes — way below the
	// minimum (1 + salt + nonce + block).
	short := resealed(t, raw[:3])
	if _, _, err := OpenSecret("pw", short); err == nil {
		t.Fatal("OpenSecret must reject a truncated envelope")
	}
}

func TestOpenSecretBadVersionFails(t *testing.T) {
	raw := sealedEnvelope(t, "pw", "secret")
	raw[0] = 99 // unsupported version
	if _, _, err := OpenSecret("pw", resealed(t, raw)); err == nil {
		t.Fatal("OpenSecret must reject an unsupported version byte")
	}
}

func TestOpenSecretFlippedTagFails(t *testing.T) {
	raw := sealedEnvelope(t, "pw", "secret")
	// Flip the last byte (part of the GCM authentication tag).
	raw[len(raw)-1] ^= 0x01
	if _, _, err := OpenSecret("pw", resealed(t, raw)); err == nil {
		t.Fatal("OpenSecret must reject a flipped GCM tag (auth failure)")
	}
}

func TestOpenSecretFlippedCiphertextFails(t *testing.T) {
	raw := sealedEnvelope(t, "pw", "secret")
	// Flip a ciphertext byte (just past salt+nonce).
	idx := 1 + secretSaltLen + secretNonceLen
	raw[idx] ^= 0x01
	if _, _, err := OpenSecret("pw", resealed(t, raw)); err == nil {
		t.Fatal("OpenSecret must reject altered ciphertext")
	}
}

func TestOpenSecretOversizedEncodedRejected(t *testing.T) {
	huge := secretPrefix + strings.Repeat("A", secretMaxEncoded)
	if _, _, err := OpenSecret("pw", huge); err == nil {
		t.Fatal("OpenSecret must reject an oversized encoded envelope")
	}
}

// ---- LoadString corrupted-data paths ----

func TestLoadStringTruncatedFails(t *testing.T) {
	// version byte only, no flags/payload.
	body := base64.RawStdEncoding.EncodeToString([]byte{StringSessionVersion})
	if _, err := LoadString("v1."+body, ""); err == nil {
		t.Fatal("LoadString must reject a truncated blob")
	}
}

func TestLoadStringBadVersionFails(t *testing.T) {
	blob := []byte{StringSessionVersion + 1, 0x00}
	body := base64.RawStdEncoding.EncodeToString(blob)
	if _, err := LoadString("v1."+body, ""); err == nil {
		t.Fatal("LoadString must reject an unsupported version")
	}
}

func TestLoadStringOversizedEncodedRejected(t *testing.T) {
	body := strings.Repeat("A", maxStringSessionEncoded+1)
	if _, err := LoadString("v1."+body, ""); err == nil {
		t.Fatal("LoadString must reject an oversized encoded string before decode")
	}
}

func TestLoadStringEncryptedFlippedTagFails(t *testing.T) {
	sess := NewStringSession()
	if err := sess.Save(sampleRecord()); err != nil {
		t.Fatal(err)
	}
	enc, err := sess.Encode("pass")
	if err != nil {
		t.Fatal(err)
	}
	body := enc[3:]
	raw, err := base64.RawStdEncoding.DecodeString(body)
	if err != nil {
		t.Fatal(err)
	}
	// Flip the final GCM tag byte of the payload.
	raw[len(raw)-1] ^= 0x01
	corrupt := "v1." + base64.RawStdEncoding.EncodeToString(raw)
	if _, err := LoadString(corrupt, "pass"); err == nil {
		t.Fatal("LoadString must reject a flipped tag in an encrypted payload")
	}
}

func TestLoadStringEncryptedNoPassphraseFails(t *testing.T) {
	sess := NewStringSession()
	if err := sess.Save(sampleRecord()); err != nil {
		t.Fatal(err)
	}
	enc, _ := sess.Encode("pass")
	if _, err := LoadString(enc, ""); err == nil {
		t.Fatal("LoadString must require a passphrase for an encrypted payload")
	}
}

// ---- FileStore hardening ----

func TestFileStoreWrittenFileMode0600(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFileStoreWithPassphrase(filepath.Join(dir, "enc"), "pass")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Save(sampleRecord()); err != nil {
		t.Fatalf("Save: %v", err)
	}
	fi, err := os.Stat(s.Path())
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("file mode = %v, want 0600", fi.Mode().Perm())
	}
}

func TestFileStoreNoTempLeftover(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFileStore(filepath.Join(dir, "s"))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Save(sampleRecord()); err != nil {
		t.Fatalf("Save: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), ".tmp") {
			t.Fatalf("temp leftover present after Save: %s", e.Name())
		}
	}
}

func TestFileStoreRejectsOversizedFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "big.bin")
	// Write a file larger than the load cap.
	if err := os.WriteFile(path, make([]byte, maxFileStoreBytes+1), 0o600); err != nil {
		t.Fatal(err)
	}
	s, err := NewFileStore(path)
	if err != nil {
		// NewFileStore appends .bin; point it at the exact file instead.
		t.Fatal(err)
	}
	if _, err := s.Load(); err == nil {
		t.Fatal("Load must reject a file larger than maxFileStoreBytes")
	}
}
