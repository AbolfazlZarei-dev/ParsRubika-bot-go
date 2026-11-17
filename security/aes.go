package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

// pkcs7Pad appends PKCS#7 padding so the result is a multiple of
// `blockSize` bytes long. It always returns a freshly allocated slice
// and never mutates the caller's `data` backing array (a previous
// implementation used append, which could alias and overwrite the
// caller's slice when it had spare capacity).
func pkcs7Pad(data []byte, blockSize int) []byte {
	pad := blockSize - (len(data) % blockSize)
	out := make([]byte, len(data)+pad)
	copy(out, data)
	for i := len(data); i < len(out); i++ {
		out[i] = byte(pad)
	}
	return out
}

// pkcs7Unpad reverses pkcs7Pad. It returns an error if the padding is
// malformed (constant-time comparison is used for the padding bytes).
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	n := len(data)
	if n == 0 || n%blockSize != 0 {
		return nil, errors.New("security: invalid padded data length")
	}
	pad := int(data[n-1])
	if pad == 0 || pad > blockSize {
		return nil, errors.New("security: invalid padding value")
	}
	// Constant-time check: every padding byte must equal `pad`.
	var mismatch byte
	for i := n - pad; i < n; i++ {
		mismatch |= data[i] ^ byte(pad)
	}
	if mismatch != 0 {
		return nil, errors.New("security: inconsistent padding bytes")
	}
	return data[:n-pad], nil
}

// Encrypt AES-256-CBC-encrypts `data` with the supplied key (32 ASCII
// chars → AES-256). The output is base64 standard-encoded so it can be
// embedded directly in a JSON payload.
func Encrypt(data, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("security: invalid aes key: %w", err)
	}
	padded := pkcs7Pad(data, aesBlockSize)
	ciphertext := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, aesIV).CryptBlocks(ciphertext, padded)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// EncryptJSON is a convenience helper that JSON-encodes `v` before
// calling Encrypt.
func EncryptJSON(v any, key []byte) (string, error) {
	raw, err := jsonMarshal(v)
	if err != nil {
		return "", fmt.Errorf("security: json marshal: %w", err)
	}
	return Encrypt(raw, key)
}

// Decrypt reverses Encrypt. It tries URL-safe base64 first (the wire
// format), then standard base64, then URL-safe without padding, so it
// is resilient to small variations in the server response.
func Decrypt(data string, key []byte) ([]byte, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		ciphertext, err = base64.StdEncoding.DecodeString(data)
		if err != nil {
			ciphertext, err = base64.RawURLEncoding.DecodeString(data)
			if err != nil {
				return nil, fmt.Errorf("security: base64 decode: %w", err)
			}
		}
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("security: invalid aes key: %w", err)
	}
	if len(ciphertext) == 0 || len(ciphertext)%aesBlockSize != 0 {
		return nil, errors.New("security: ciphertext is not a multiple of the block size")
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, aesIV).CryptBlocks(plaintext, ciphertext)
	return pkcs7Unpad(plaintext, aesBlockSize)
}

// jsonMarshal is indirected through a helper so tests can stub it out.
func jsonMarshal(v any) ([]byte, error) {
	return jsonEncode(v)
}
