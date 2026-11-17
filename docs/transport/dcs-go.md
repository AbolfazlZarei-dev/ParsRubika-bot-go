# `transport` — Dcs

_Source: [`dcs.go`](../../ParsRubika/transport/dcs.go)_

## Symbols

### `DCDiscoveryURLs`

```go
var DCDiscoveryURLs = []string{
```

DCDiscoveryURLs is the list of endpoints tried, in order, to discover
the active Rubika data-centers. If one fails (DNS, TLS, timeout,
HTTP 5xx) the next is tried. The primary endpoint has historically
been the most reliable; the fallbacks cover regional outages.

### `DCConfig`

```go
type DCConfig struct
```

DCConfig holds the discovered API map + WebSocket URL.

### `DiscoverDCs`

```go
func DiscoverDCs(ctx context.Context, t *Transport) (*DCConfig, error)
```

DiscoverDCs fetches the API/WS endpoints from the data-center
discovery service. It tries every URL in `DCDiscoveryURLs` until one
returns a complete configuration. When all endpoints fail, the
returned error contains a short diagnostic for each endpoint so the
caller can tell whether the failure was DNS, TLS, timeout, or a
server-side 5xx — instead of the opaque "all DC discovery endpoints
failed" message.

### `PostFormFile`

```go
func PostFormFile(ctx context.Context, t *Transport, url, fieldName, fileName string, body io.Reader) ([]byte, int, error)
```

PostFormFile uploads a single file via multipart/form-data. It is
used by the file-upload endpoint returned by `requestSendFile`.

Errors from the pipe-writer goroutine (e.g. the file reader fails
mid-stream) are propagated via CloseWithError so the HTTP request
fails with a real error instead of completing with truncated data.

### `DCMaxRetries`

```go
const DCMaxRetries = 3
```

Constants exported for callers that need to tune the discovery loop.

### `DCBackoffDelay`

```go
const DCBackoffDelay = time.Second
```

Constants exported for callers that need to tune the discovery loop.
