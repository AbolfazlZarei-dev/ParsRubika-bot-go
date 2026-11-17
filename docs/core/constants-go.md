# `core` — Constants

_Source: [`constants.go`](../../ParsRubika/core/constants.go)_

## Symbols

### `Identity`

```go
type Identity struct
```

Identity holds the client-side identity sent inside the `client` field
of every request body. The exact values are tuned to mimic the
official web / Android clients so the server returns the correct
feature flags.

### `WebIdentity`

```go
func WebIdentity(lang string) Identity
```

WebIdentity returns the identity used by PWA (web) clients.

### `AndroidIdentity`

```go
func AndroidIdentity(lang string) Identity
```

AndroidIdentity returns the identity used by Android clients.

### `Platform`

```go
type Platform string
```

Platform is the top-level selector used by the account client.

### `PlatformWeb`

```go
const PlatformWeb Platform = "PWA"
```

### `PlatformAndroid`

```go
const PlatformAndroid Platform = "Android"
```

### `AllDCDiscoveryURLs`

```go
var AllDCDiscoveryURLs = []string{
```

AllDCDiscoveryURLs is the list of endpoints tried, in order, to
discover the active Rubika data-centers. The transport package is the canonical
source (core.AllDCDiscoveryURLs is kept for backwards compatibility).

### `DefaultWebSocketPingInterval`

```go
const DefaultWebSocketPingInterval = 10 // seconds
```

DefaultWebSocketPingInterval controls how often the keep-alive ping
is sent over the WebSocket connection.

### `MaxMessageTextLength`

```go
const MaxMessageTextLength = 4200
```

MaxMessageTextLength is the server-imposed maximum size of a single
text message; longer texts are chunked into multiple messages.

### `RSAKeyBits`

```go
const RSAKeyBits = 2048
```

RSAKeyBits is the size (in bits) of the RSA key generated at sign-in
time. The Rubika server accepts 2048-bit keys (the official MaxRubika
client uses RSA-2048); 2048 is the minimum modern baseline and is
resistant to feasible factorisation.

Existing on-disk sessions keep using whatever key they were created
with (the private key is loaded verbatim); only NEW sign-ins generate
a 2048-bit key. See MIGRATION.md for details.

### `AuthLength`

```go
const AuthLength = 32
```

AuthLength is the length of the auth key string (32 chars).

### `PasskeyLength`

```go
const PasskeyLength = 32
```

PasskeyLength is the length of the derived AES key (32 chars → AES-256).
