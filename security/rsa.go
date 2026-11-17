package security

import (
	stdcrypto "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

// sha256Hash is the crypto.Hash identifier required by rsa.SignPKCS1v15.
const sha256Hash = stdcrypto.SHA256

// SignRSA returns the base64-encoded PKCS#1 v1.5 signature of `data`
// under the supplied RSA private key. SHA-256 is used as the digest.
func SignRSA(priv *rsa.PrivateKey, data string) (string, error) {
	if priv == nil {
		return "", errors.New("security: nil private key")
	}
	h := sha256.Sum256([]byte(data))
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, sha256Hash, h[:])
	if err != nil {
		return "", fmt.Errorf("security: rsa sign: %w", err)
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

// DecryptOAEP decrypts `data` (base64-encoded) using RSA-OAEP with
// SHA-1 (matching PyCryptodome's PKCS1_OAEP default).
func DecryptOAEP(priv *rsa.PrivateKey, data string) (string, error) {
	if priv == nil {
		return "", errors.New("security: nil private key")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("security: base64 decode: %w", err)
	}
	// PyCryptodome PKCS1_OAEP uses SHA-1 by default, NOT SHA-256.
	plaintext, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("security: rsa oaep decrypt: %w", err)
	}
	return string(plaintext), nil
}

// GenerateKeyPair generates a fresh RSA keypair with the supplied bit
// size. It returns:
//
// - publicKey : PEM-encoded RSA public key in SubjectPublicKeyInfo
// (PKIX) format — exactly what PyCryptodome's
// `key.publickey().export_key()` produces by default
// (`-----BEGIN PUBLIC KEY-----`). The Rubika server's RSA encryptor
// expects this format; using PKCS#1 (`-----BEGIN RSA PUBLIC KEY-----`)
// here causes the server to encrypt the returned auth with a key
// derived from the wrong ASN.1 structure, which surfaces as
// `security: invalid padding value` during DecryptOAEP.
// The PEM is then passed through DecodeAuth.
// - privateKey : PEM-encoded PKCS#1 RSA private key, as a string.
// - err : any error from key generation / encoding.
func GenerateKeyPair(bits int) (publicKey string, privateKey string, err error) {
	if bits <= 0 {
		bits = 1024
	}
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", fmt.Errorf("security: rsa generate: %w", err)
	}
	privDER := x509.MarshalPKCS1PrivateKey(priv)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privDER})
	// Use PKIX (SubjectPublicKeyInfo) format — this is what PyCryptodome
	// produces by default with key.publickey().export_key().
	pubDER, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("security: marshal public key: %w", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})
	pubB64 := base64.StdEncoding.EncodeToString(pubPEM)
	return DecodeAuth(pubB64), string(privPEM), nil
}

// ParsePrivateKey parses a PEM-encoded RSA private key (PKCS#1 or
// PKCS#8). The input may include the `[REDACTED:...]` marker added by
// some clients — it is stripped automatically.
func ParsePrivateKey(s string) (*rsa.PrivateKey, error) {
	// Strip everything before the first "-----BEGIN".
	if i := strings.Index(s, "-----BEGIN"); i > 0 {
		s = s[i:]
	}
	block, _ := pem.Decode([]byte(s))
	if block == nil {
		return nil, errors.New("security: failed to parse PEM private key")
	}
	if k, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return k, nil
	}
	if k, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rk, ok := k.(*rsa.PrivateKey); ok {
			return rk, nil
		}
	}
	return nil, errors.New("security: unsupported private key format")
}

// ParsePublicKey parses a PEM-encoded RSA public key (PKCS#1 or PKCS#8).
func ParsePublicKey(s string) (*rsa.PublicKey, error) {
	if i := strings.Index(s, "-----BEGIN"); i > 0 {
		s = s[i:]
	}
	block, _ := pem.Decode([]byte(s))
	if block == nil {
		return nil, errors.New("security: failed to parse PEM public key")
	}
	if k, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		return k, nil
	}
	// PKIX form (SubjectPublicKeyInfo).
	if k, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		if rk, ok := k.(*rsa.PublicKey); ok {
			return rk, nil
		}
	}
	return nil, errors.New("security: unsupported public key format")
}
