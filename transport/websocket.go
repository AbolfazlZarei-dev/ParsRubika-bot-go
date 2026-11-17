package transport

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/net/proxy"

	coreerrors "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/errors"
)

// MaxWSMessage caps the size of a single inbound WebSocket frame. A
// malicious or buggy peer could otherwise stream an unbounded frame and
// exhaust memory; gorilla rejects the read once it exceeds the limit.
const MaxWSMessage = 8 << 20 // 8 MiB

// wsReadDeadline bounds how long a single Read may block waiting for the
// next frame. Without it a peer that opens the connection but never
// sends (and never closes) can hold the read goroutine forever. The
// keep-alive loop (which sends a frame every interval) normally keeps
// the server responsive well within this window, so the deadline only
// fires on a genuinely dead connection — at which point surfacing the
// error and reconnecting is the correct behaviour.
const wsReadDeadline = 120 * time.Second

// WebSocket wraps a gorilla/websocket.Conn and adds:
//
// - a background keep-alive ping loop
// - a mutex-protected write path (so multiple goroutines can send)
// - a read size limit + read deadline (resource-exhaustion / hang defense)
// - optional proxy support (HTTP/HTTPS and SOCKS5)
// - automatic reconnection support (managed by the caller)
type WebSocket struct {
	url      string
	conn     *websocket.Conn
	connMu   sync.Mutex
	headers  map[string]string
	proxyURL string
}

// NewWebSocket constructs a WebSocket handle. The connection is not
// dialled yet — call Dial to open it.
func NewWebSocket(url string, headers map[string]string) *WebSocket {
	return &WebSocket{url: url, headers: headers}
}

// SetProxy configures an HTTP/HTTPS or SOCKS5 proxy for the next Dial.
// Must be called before Dial. Supports the same URL formats as the HTTP
// transport (http://, https://, socks5://, socks5h://).
func (w *WebSocket) SetProxy(proxyURL string) { w.proxyURL = proxyURL }

// Dial opens the connection and returns the initial handshake response
// (so the caller can send the handShake frame). The headers supplied
// to NewWebSocket (User-Agent, auth, ...) are forwarded to the server
// as part of the HTTP upgrade request.
func (w *WebSocket) Dial(ctx context.Context) error {
	dialer := websocket.Dialer{
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: false},
		HandshakeTimeout: 10 * time.Second,
		// Cap inbound frame size at the transport level too, so an
		// oversized frame is rejected during the read rather than after.
		ReadBufferSize: 4096,
	}
	if err := applyWSProxy(&dialer, w.proxyURL); err != nil {
		return coreerrors.Wrap(coreerrors.KindNetwork, "ws proxy", err)
	}
	header := http.Header{}
	for k, v := range w.headers {
		header.Set(k, v)
	}
	conn, _, err := dialer.DialContext(ctx, w.url, header)
	if err != nil {
		return coreerrors.Wrap(coreerrors.KindNetwork, "ws dial", err)
	}
	// Bound the size of any single inbound message.
	conn.SetReadLimit(MaxWSMessage)
	w.connMu.Lock()
	w.conn = conn
	w.connMu.Unlock()
	return nil
}

// applyWSProxy configures the gorilla dialer for an HTTP/HTTPS or SOCKS5
// proxy. Returns an error on a malformed URL so the caller fails closed
// instead of dialing directly.
func applyWSProxy(d *websocket.Dialer, proxyURL string) error {
	if proxyURL == "" {
		return nil
	}
	u, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("parse proxy url: %w", err)
	}
	switch strings.ToLower(u.Scheme) {
	case "socks5", "socks5h":
		var auth *proxy.Auth
		if u.User != nil {
			pw, _ := u.User.Password()
			auth = &proxy.Auth{User: u.User.Username(), Password: pw}
		}
		sd, err := proxy.SOCKS5("tcp", u.Host, auth, &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("socks5 dialer: %w", err)
		}
		d.NetDialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return sd.Dial(network, addr)
		}
	default:
		d.Proxy = http.ProxyURL(u)
	}
	return nil
}

// SendJSON sends a JSON message over the WebSocket.
func (w *WebSocket) SendJSON(v any) error {
	w.connMu.Lock()
	defer w.connMu.Unlock()
	if w.conn == nil {
		return coreerrors.New(coreerrors.KindNetwork, "ws not connected")
	}
	return w.conn.WriteJSON(v)
}

// Read blocks until the next message arrives. It returns the raw bytes
// (text frames are already decoded; binary frames are passed as-is).
//
// A read deadline (wsReadDeadline) is set before each read so a peer
// that goes silent without closing the connection cannot hold the read
// goroutine indefinitely; the deadline fires as a network error and the
// caller's reconnect loop takes over.
func (w *WebSocket) Read() ([]byte, error) {
	w.connMu.Lock()
	conn := w.conn
	w.connMu.Unlock()
	if conn == nil {
		return nil, coreerrors.New(coreerrors.KindNetwork, "ws not connected")
	}
	_ = conn.SetReadDeadline(time.Now().Add(wsReadDeadline))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Close terminates the WebSocket connection.
func (w *WebSocket) Close() error {
	w.connMu.Lock()
	defer w.connMu.Unlock()
	if w.conn == nil {
		return nil
	}
	err := w.conn.Close()
	w.conn = nil
	return err
}

// KeepAlive sends an empty JSON object every `interval`. It blocks
// until `ctx` is cancelled.
func (w *WebSocket) KeepAlive(ctx context.Context, interval time.Duration, onTick func() error) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.SendJSON(map[string]any{}); err != nil {
				return err
			}
			if onTick != nil {
				if err := onTick(); err != nil {
					return err
				}
			}
		}
	}
}

// Handshake sends the initial handShake frame required by the server.
func (w *WebSocket) Handshake(auth string) error {
	return w.SendJSON(map[string]any{
		"method":      "handShake",
		"auth":        auth,
		"api_version": "6",
		"data":        "",
	})
}

// DecodeMessage parses a WebSocket text frame into a generic map.
func DecodeMessage(raw []byte) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("transport: ws decode: %w", err)
	}
	return m, nil
}
