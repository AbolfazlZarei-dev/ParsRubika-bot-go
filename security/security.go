// Package security implements every cryptographic primitive used by
// ParsRubika. It is split across several small files so each concern is
// easy to inspect independently:
//
// - aes.go : AES-256-CBC encryption / decryption + PKCS#7 padding
// - rsa.go : RSA-PKCS1v15 signing + RSA-OAEP decryption + key gen
// - auth.go : decode_auth + passphrase derivation
// - random.go : cryptographically secure random generation
// - keys.go : key parsing / formatting helpers
package security

// AES block size in bytes (16 for AES-128/192/256).
const aesBlockSize = 16

// aesIV is the all-zero IV used by the wire protocol. It must remain
// constant so that messages stay inter-operable with the server.
var aesIV = make([]byte, aesBlockSize)
