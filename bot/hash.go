package bot

import "crypto/sha256"

// sha256Sum is a small indirection to avoid importing crypto/sha256 in
// the main client file.
func sha256Sum(s []byte) []byte {
	h := sha256.Sum256(s)
	return h[:]
}
