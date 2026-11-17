package account

import (
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/security"
)

// generateKeyPairImpl wraps security.GenerateKeyPair so methods_auth2fa.go
// doesn't need to import security directly.
func generateKeyPairImpl() (pub, priv string, err error) {
	return security.GenerateKeyPair(1024)
}
