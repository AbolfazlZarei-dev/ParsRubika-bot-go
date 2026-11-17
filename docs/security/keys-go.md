# `security` — Keys

_Source: [`keys.go`](../../ParsRubika/security/keys.go)_

## Symbols

### `NormalizePrivateKey`

```go
func NormalizePrivateKey(s string) string
```

NormalizePrivateKey cleans up a user-supplied RSA private-key string:

- Strips any "[REDACTED:ssh_private_key]" redaction marker that may
have been pasted in (the marker itself is removed, not just the
text before it).
- Ensures the PEM body has a proper "-----BEGIN …-----" / "-----END …-----"
envelope, appending the footer only if it is genuinely missing.
- Validates the result by parsing it as a PEM-encoded private key.

Returns the cleaned string. If the input cannot be parsed as a PEM
private key the original (cleaned) string is returned unchanged so
that the caller's subsequent ParsePrivateKey call surfaces a clear
error instead of a silently-mangled key.
