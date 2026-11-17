# `netutil` — Package netutil provides SSRF defenses for outbound HTTP requests the

_Source: [`netutil.go`](../../ParsRubika/netutil/netutil.go)_

## Symbols

### `ErrBlocked`

```go
var ErrBlocked = errors.New("netutil: address blocked by SSRF policy")
```

ErrBlocked is returned when a URL or address is refused by the SSRF
policy (private / loopback / link-local / metadata IP, or a forbidden
scheme).

### `AllowedSchemes`

```go
var AllowedSchemes = map[string]bool{"http": true, https: true}
```

AllowedSchemes are the URL schemes permitted for outbound fetches.

### `IsBlockedIP`

```go
func IsBlockedIP(ip net.IP) bool
```

IsBlockedIP reports whether ip points at a host the library must never
contact on the user's behalf: loopback, private RFC-1918 ranges,
link-local (incl. the cloud metadata endpoint 169.254.169.254),
unspecified, multicast, and other non-global-unicast addresses.

### `ValidateURL`

```go
func ValidateURL(rawURL string) (*url.URL, error)
```

ValidateURL parses rawURL and applies the syntactic SSRF policy:
http/https scheme only, non-empty host, no localhost, and an IP-literal
host (if any) that is not private. It does NOT perform DNS resolution,
so a hostname that resolves to a private IP is only caught by
AssertURL or NewSafeTransport.

### `AssertURL`

```go
func AssertURL(ctx context.Context, rawURL string) (*url.URL, error)
```

AssertURL runs ValidateURL and then resolves the host name (when it is
not already an IP literal) and verifies that NONE of the resolved
addresses are private. The check uses the default resolver with the
supplied context so it honours cancellation / timeouts.

AssertURL is a pre-flight check. For full defense against DNS
rebinding, also route the request through NewSafeTransport.

### `ControlSSRF`

```go
func ControlSSRF(network, address string, _ syscall.RawConn) error
```

ControlSSRF is the Dialer.Control hook: it receives the resolved
IP:port just before connect and refuses private addresses. This runs
AFTER name resolution, so it sees the real address being contacted.

It is exported so the main API/WebSocket transport can install the
same SSRF defense used by NewSafeTransport, without having to replace
its whole http.Transport.

### `NewSafeTransport`

```go
func NewSafeTransport(timeout time.Duration) *http.Transport
```

NewSafeTransport returns an *http.Transport that blocks dials to
private / loopback / link-local / metadata IPs via a Control hook.
Use it as the transport of any http.Client that fetches user-supplied
URLs (uploads, scrapes, JSON fetches).

### `NewSafeClient`

```go
func NewSafeClient(timeout time.Duration) *http.Client
```

NewSafeClient returns an http.Client whose transport blocks SSRF
targets and whose overall timeout is bounded.

### `MaxFetchBytes`

```go
const MaxFetchBytes = 50 << 20 // 50 MiB
```

MaxFetchBytes caps how many bytes a user-URL fetch may read into memory
(UploadFromURL, FetchJSON, ScrapeChannelPostLink). It prevents a
malicious remote server from exhausting RAM by streaming an enormous
response.
