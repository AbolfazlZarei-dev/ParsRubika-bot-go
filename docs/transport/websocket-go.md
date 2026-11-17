# `transport` — Websocket

_Source: [`websocket.go`](../../ParsRubika/transport/websocket.go)_

## Symbols

### `MaxWSMessage`

```go
const MaxWSMessage = 8 << 20 // 8 MiB
```

MaxWSMessage caps the size of a single inbound WebSocket frame. A
malicious or buggy peer could otherwise stream an unbounded frame and
exhaust memory; gorilla rejects the read once it exceeds the limit.

### `WebSocket`

```go
type WebSocket struct
```

WebSocket wraps a gorilla/websocket.Conn and adds:

- a background keep-alive ping loop
- a mutex-protected write path (so multiple goroutines can send)
- a read size limit + read deadline (resource-exhaustion / hang defense)
- optional proxy support (HTTP/HTTPS and SOCKS5)
- automatic reconnection support (managed by the caller)

### `NewWebSocket`

```go
func NewWebSocket(url string, headers map[string]string) *WebSocket
```

NewWebSocket constructs a WebSocket handle. The connection is not
dialled yet — call Dial to open it.

### `SetProxy`

```go
func (w *WebSocket) SetProxy(proxyURL string) { w.proxyURL = proxyURL }
```

SetProxy configures an HTTP/HTTPS or SOCKS5 proxy for the next Dial.
Must be called before Dial. Supports the same URL formats as the HTTP
transport (http://, https://, socks5://, socks5h://).

### `Dial`

```go
func (w *WebSocket) Dial(ctx context.Context) error
```

Dial opens the connection and returns the initial handshake response
(so the caller can send the handShake frame). The headers supplied
to NewWebSocket (User-Agent, auth, ...) are forwarded to the server
as part of the HTTP upgrade request.

### `SendJSON`

```go
func (w *WebSocket) SendJSON(v any) error
```

SendJSON sends a JSON message over the WebSocket.

### `Read`

```go
func (w *WebSocket) Read() ([]byte, error)
```

Read blocks until the next message arrives. It returns the raw bytes
(text frames are already decoded; binary frames are passed as-is).

A read deadline (wsReadDeadline) is set before each read so a peer
that goes silent without closing the connection cannot hold the read
goroutine indefinitely; the deadline fires as a network error and the
caller's reconnect loop takes over.

### `Close`

```go
func (w *WebSocket) Close() error
```

Close terminates the WebSocket connection.

### `KeepAlive`

```go
func (w *WebSocket) KeepAlive(ctx context.Context, interval time.Duration, onTick func() error) error
```

KeepAlive sends an empty JSON object every `interval`. It blocks
until `ctx` is cancelled.

### `Handshake`

```go
func (w *WebSocket) Handshake(auth string) error
```

Handshake sends the initial handShake frame required by the server.

### `DecodeMessage`

```go
func DecodeMessage(raw []byte) (map[string]any, error)
```

DecodeMessage parses a WebSocket text frame into a generic map.
