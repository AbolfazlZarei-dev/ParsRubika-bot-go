# `bot` — Token string

_Source: [`token_string.go`](../../ParsRubika/bot/token_string.go)_

## Symbols

### `TokenStringPrefix`

```go
const TokenStringPrefix = "PRTS" // ParsRubika Token String
```

TokenStringPrefix is the 4-character marker that identifies a token
string produced by ExportTokenString.

### `ExportTokenString`

```go
func (b *Bot) ExportTokenString(password string) (string, error)
```

ExportTokenString encrypts the bot token (along with optional metadata)
into a portable, password-protected string. The format is:

PRTS<base64(version | salt | nonce | ciphertext | checksum)>

The ciphertext is AES-256-GCM, the key is derived via PBKDF2-SHA256
with 200,000 iterations, so tampering is detected on import.

If password is empty, a deterministic but insecure "no-password" mode
is used — the resulting string is still obfuscated but NOT encrypted,
and a one-time warning is logged. Always supply a strong password for
any real use; treat the empty-password form as a development convenience.

### `ImportTokenString`

```go
func ImportTokenString(tokenString, password string) (map[string]any, error)
```

ImportTokenString decodes a token string produced by ExportTokenString
and returns the payload map. The password must match the one used at
export time.

### `NewFromTokenString`

```go
func NewFromTokenString(tokenString, password string, opts ...Options) (*Bot, error)
```

NewFromTokenString constructs a new Bot from an encrypted token string.
It decodes the token string, extracts the token, and calls New with the
saved options (parse_mode, user_agent, etc.).
