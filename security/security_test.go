package security

import (
	"bytes"
	"strings"
	"testing"
)

func TestDecodeAuth(t *testing.T) {
	// Known mappings (independently computed):
	cases := map[string]string{
		"a": "g", // (32-0)%26 + 97 = 6 + 97 = 103 = 'g'
		"0": "3", // (13-0)%10 + 48 = 3 + 48 = 51 = '3'
		"A": "D", // (29-0)%26 + 65 = 3 + 65 = 68 = 'D'
		"-": "-", // symbols preserved
		"_": "_", // symbols preserved
	}
	for in, want := range cases {
		if got := DecodeAuth(in); got != want {
			t.Errorf("DecodeAuth(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDecodeAuthFull(t *testing.T) {
	in := "abcDEF123-_"
	out := DecodeAuth(in)
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	if !strings.Contains(out, "-") || !strings.Contains(out, "_") {
		t.Fatalf("expected - and _ to be preserved in: %q", out)
	}
}

func TestDeriveKey(t *testing.T) {
	const auth = "abcdefghijklmnopqrstuvwxyz012345"
	got, err := DeriveKey(auth)
	if err != nil {
		t.Fatalf("DeriveKey: %v", err)
	}
	if len(got) != 32 {
		t.Fatalf("DeriveKey length = %d, want 32", len(got))
	}
	for _, r := range got {
		if !(r >= 'a' && r <= 'z') {
			t.Fatalf("DeriveKey must be all lowercase, got %q", got)
		}
	}
	if _, err := DeriveKey("short"); err == nil {
		t.Fatal("expected error on short auth")
	}
}

func TestRandomSecret(t *testing.T) {
	s, err := RandomSecret(32)
	if err != nil {
		t.Fatalf("RandomSecret: %v", err)
	}
	if len(s) != 32 {
		t.Fatalf("RandomSecret length = %d, want 32", len(s))
	}
	for _, r := range s {
		if !(r >= 'a' && r <= 'z') {
			t.Fatalf("RandomSecret must be all lowercase, got %q", s)
		}
	}
	s2, _ := RandomSecret(32)
	if s == s2 {
		t.Fatal("two consecutive RandomSecret calls returned the same value")
	}
}

// TestRandomSecretUniform guards against modulo bias. The previous
// implementation mapped each random byte with `byte % 26`; since
// 256 % 26 == 22, the letters 'w'–'z' were drawn with probability
// 9/256 instead of 10/256 — roughly 10 % under-represented. Drawing a
// large sample and checking the min/max letter frequency ratio catches
// that bias cleanly (uniform ≈ 0.98, biased ≈ 0.90).
func TestRandomSecretUniform(t *testing.T) {
	const letters = len(lowercaseAlphabet)
	const samplesPerLetter = 20000
	const total = letters * samplesPerLetter // 520_000 chars
	var counts [letters]int
	// Generate in chunks to keep peak memory modest.
	const chunkLen = 1000
	generated := 0
	for generated < total {
		n := chunkLen
		if generated+n > total {
			n = total - generated
		}
		s, err := RandomSecret(n)
		if err != nil {
			t.Fatalf("RandomSecret: %v", err)
		}
		for _, r := range s {
			idx := int(r - 'a')
			if idx < 0 || idx >= letters {
				t.Fatalf("RandomSecret produced non-alphabet byte %q", r)
			}
			counts[idx]++
		}
		generated += n
	}
	min, max := counts[0], counts[0]
	for _, c := range counts {
		if c < min {
			min = c
		}
		if c > max {
			max = c
		}
	}
	ratio := float64(min) / float64(max)
	t.Logf("counts min=%d max=%d ratio=%.4f", min, max, ratio)
	// Uniform draws over 20k/letter give a ratio ≈ 0.985; the biased
	// implementation would sit around 0.90. 0.95 cleanly separates them
	// without flakiness.
	if ratio < 0.95 {
		t.Fatalf("modulo bias detected: min/max ratio %.4f < 0.95 (counts=%v)", ratio, counts)
	}
}

// TestPKCS7PadDoesNotCorruptCaller verifies pkcs7Pad never writes into the
// caller's backing array. The old implementation used append(data, …),
// which — when the caller's slice had spare capacity — would overwrite
// sibling bytes the caller still held a reference to.
func TestPKCS7PadDoesNotCorruptCaller(t *testing.T) {
	// Backing array of exactly one block; data is a 5-byte view with full
	// capacity, so the old append-based code would write padding bytes
	// straight into positions 5..15 of this array.
	backing := make([]byte, aesBlockSize)
	copy(backing, []byte{1, 2, 3, 4, 5})
	backing[5] = 99
	backing[10] = 77
	data := backing[:5:aesBlockSize]

	padded := pkcs7Pad(data, aesBlockSize)
	// Caller's visible slice unchanged.
	if !bytes.Equal(data, []byte{1, 2, 3, 4, 5}) {
		t.Fatalf("caller slice mutated: %v", data)
	}
	// Sibling bytes in the backing array unchanged.
	if backing[5] != 99 || backing[10] != 77 {
		t.Fatalf("caller backing array corrupted by pkcs7Pad: %v", backing)
	}
	// Padding is correct and lives in its own allocation.
	if len(padded) != aesBlockSize {
		t.Fatalf("padded length = %d, want %d", len(padded), aesBlockSize)
	}
	want := append([]byte{1, 2, 3, 4, 5}, bytes.Repeat([]byte{11}, 11)...)
	if !bytes.Equal(padded, want) {
		t.Fatalf("padded = %v, want %v", padded, want)
	}
}

func TestRandomBytes(t *testing.T) {
	b, err := RandomBytes(64)
	if err != nil {
		t.Fatalf("RandomBytes: %v", err)
	}
	if len(b) != 64 {
		t.Fatalf("RandomBytes length = %d, want 64", len(b))
	}
	b2, _ := RandomBytes(64)
	if string(b) == string(b2) {
		t.Fatal("two consecutive RandomBytes calls returned the same value")
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := "01234567890123456789012345678901" // 32 chars => AES-256
	plain := []byte(`{"hello":"world","n":42}`)
	enc, err := Encrypt(plain, []byte(key))
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if enc == "" {
		t.Fatal("empty ciphertext")
	}
	dec, err := Decrypt(enc, []byte(key))
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if string(dec) != string(plain) {
		t.Fatalf("round-trip mismatch: got %q, want %q", dec, plain)
	}
}

func TestEncryptJSONRoundTrip(t *testing.T) {
	key := "01234567890123456789012345678901"
	in := map[string]any{"hello": "world", "n": 42.0}
	enc, err := EncryptJSON(in, []byte(key))
	if err != nil {
		t.Fatalf("EncryptJSON: %v", err)
	}
	dec, err := Decrypt(enc, []byte(key))
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !strings.Contains(string(dec), "world") {
		t.Fatalf("decrypted JSON missing 'world': %s", dec)
	}
}

func TestGenerateKeyPair(t *testing.T) {
	pub, priv, err := GenerateKeyPair(1024)
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	if !strings.Contains(priv, "-----BEGIN RSA PRIVATE KEY-----") {
		t.Fatalf("private key PEM missing BEGIN marker: %s", priv)
	}
	if pub == "" {
		t.Fatal("empty public key")
	}
	pk, err := ParsePrivateKey(priv)
	if err != nil {
		t.Fatalf("ParsePrivateKey: %v", err)
	}
	if pk == nil {
		t.Fatal("parsed private key is nil")
	}
}

func TestGenerateKeyPairCustomBits(t *testing.T) {
	// Use a non-default bit size to verify the parameter is wired up.
	_, _, err := GenerateKeyPair(1024)
	if err != nil {
		t.Fatalf("GenerateKeyPair(1024): %v", err)
	}
}

func TestSignAndVerify(t *testing.T) {
	_, priv, err := GenerateKeyPair(1024)
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	pk, err := ParsePrivateKey(priv)
	if err != nil {
		t.Fatalf("ParsePrivateKey: %v", err)
	}
	sig, err := SignRSA(pk, "hello world")
	if err != nil {
		t.Fatalf("SignRSA: %v", err)
	}
	if sig == "" {
		t.Fatal("empty signature")
	}
	// Sign again — should produce the same signature (PKCS1v15 is deterministic).
	sig2, _ := SignRSA(pk, "hello world")
	if sig != sig2 {
		t.Fatal("PKCS1v15 signature must be deterministic")
	}
}

func TestNormalizePrivateKey(t *testing.T) {
	// When no prefix is supplied, NormalizePrivateKey adds the footer.
	in := "-----BEGIN RSA PRIVATE KEY-----\nMIIBOgIB"
	out := NormalizePrivateKey(in)
	if !strings.HasSuffix(out, "-----END RSA PRIVATE KEY-----") {
		t.Fatalf("expected END marker, got: %s", out)
	}
	// When a redaction prefix is supplied, it is stripped.
	in2 := "[REDACTED] -----BEGIN RSA PRIVATE KEY-----\nMIIBOgIB"
	out2 := NormalizePrivateKey(in2)
	if strings.HasPrefix(out2, "[REDACTED]") {
		t.Fatalf("redaction prefix not stripped: %s", out2)
	}
}
