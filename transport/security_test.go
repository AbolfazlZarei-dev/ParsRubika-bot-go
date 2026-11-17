package transport

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// newTestTransport builds a Transport whose client points at the given
// test server and whose retry backoff is short so tests stay fast.
func newTestTransport(t *testing.T, base string) *Transport {
	t.Helper()
	tr := New()
	tr.httpClient = &http.Client{Timeout: 5 * time.Second}
	return tr
}

// --- POST retry semantics: no duplicate on ambiguous network error ---

func TestPostDoesNotRetryOnNetworkError(t *testing.T) {
	// A TCP listener that accepts exactly one connection and closes it
	// immediately, producing a transport-level error on the client.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	var accepts atomic.Int64
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			accepts.Add(1)
			c.Close() // force a connection error, no HTTP response
		}
	}()

	tr := New(WithTimeout(2 * time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, _, err = tr.Post(ctx, "http://"+ln.Addr().String()+"/", map[string]any{"x": 1})
	if err == nil {
		t.Fatal("expected a transport error, got nil")
	}
	// Exactly one TCP attempt: the POST must NOT have retried on the
	// ambiguous network error (which would risk a duplicate message).
	if got := accepts.Load(); got != 1 {
		t.Fatalf("expected exactly 1 dial on network error (no retry), got %d", got)
	}
}

// --- POST DOES retry on an explicit server-side retryable status ---

func TestPostRetriesOnServerStatusThenSucceeds(t *testing.T) {
	var calls atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if n == 1 {
			http.Error(w, "busy", http.StatusServiceUnavailable) // 503 → retryable
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	tr := New(WithTimeout(5 * time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body, status, err := tr.Post(ctx, srv.URL, map[string]any{"x": 1})
	if err != nil {
		t.Fatalf("Post: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
	if !strings.Contains(string(body), `"ok":true`) {
		t.Fatalf("unexpected body: %s", body)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("expected 2 attempts (1 retry on 503), got %d", got)
	}
}

// --- Response body size limit ---

func TestResponseBodyLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(make([]byte, MaxResponseBody+1))
	}))
	defer srv.Close()

	tr := newTestTransport(t, srv.URL)
	_, _, err := tr.Get(context.Background(), srv.URL)
	if err != ErrResponseBodyTooLarge {
		t.Fatalf("expected ErrResponseBodyTooLarge, got %v", err)
	}
}

// --- Error scrubbing: no PEM / secrets / over-long body in error ---

func TestErrorScrubbing(t *testing.T) {
	secret := "-----BEGIN RSA PRIVATE KEY-----\nSUPERSECRETKEYDATA\n-----END RSA PRIVATE KEY-----"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, strings.Repeat("pad", 400)+secret)
	}))
	defer srv.Close()

	tr := newTestTransport(t, srv.URL)
	_, _, err := tr.Get(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected an error")
	}
	msg := err.Error()
	if strings.Contains(msg, "SUPERSECRETKEYDATA") {
		t.Fatalf("secret leaked into error message: %s", msg)
	}
	if strings.Contains(msg, "-----BEGIN RSA") {
		t.Fatalf("PEM marker leaked into error message: %s", msg)
	}
	// The embedded snippet must be capped.
	if len(msg) > maxErrorSnippet+80 { // allow room for the "transport: http 400: " prefix
		t.Fatalf("error message not truncated: len=%d", len(msg))
	}
}

// --- Proxy fail-closed: bad proxy must NOT fall back to a direct conn ---

func TestProxyFailClosed(t *testing.T) {
	var hits atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cli := NewHTTPClient(2*time.Second, "bad-scheme://%%%", false)
	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	_, err := cli.Do(req)
	if err == nil {
		t.Fatal("expected error because the proxy is unavailable")
	}
	// The server must NOT have been contacted directly.
	if got := hits.Load(); got != 0 {
		t.Fatalf("proxy fail-open: server was contacted %d times directly", got)
	}
}

// --- Redirects disabled on the API client (auth body never re-sent) ---

func TestRedirectsDisabled(t *testing.T) {
	var hits atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		if r.URL.Path == "/start" {
			http.Redirect(w, r, "/elsewhere", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cli := NewHTTPClient(2*time.Second, "", false)
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/start", nil)
	resp, err := cli.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 (redirect NOT followed), got %d", resp.StatusCode)
	}
	if got := hits.Load(); got != 1 {
		t.Fatalf("expected 1 request (no redirect follow), got %d", got)
	}
}

// --- Concurrent SetHeader + applyHeaders under the race detector ---

func TestSetHeaderConcurrent(t *testing.T) {
	tr := New(WithTimeout(2 * time.Second))
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			tr.SetHeader("x-token", strings.Repeat("a", i%8+1))
		}(i)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest(http.MethodGet, "http://example.invalid/", nil)
			tr.applyHeaders(req)
		}()
	}
	wg.Wait()
}

// --- DC URL validation ---

func TestParseDCsRejectsInsecureAPIURL(t *testing.T) {
	// http (not https) API URL must be dropped; wss with a private IP
	// literal must be dropped; a valid https + wss pair must survive.
	body := []byte(`{"data":{"API":{"1":"http://evil.example/","2":"https://api.rubika.ir/"},"socket":{"1":"wss://10.0.0.1/socket","2":"wss://socket.rubika.ir/socket"},"default_api":"2","default_socket":"2"}}`)
	cfg, err := parseDCs(body)
	if err != nil {
		t.Fatalf("parseDCs: %v", err)
	}
	if _, ok := cfg.APIMap["1"]; ok {
		t.Fatal("insecure http API URL was accepted")
	}
	if cfg.APIURL != "https://api.rubika.ir/" {
		t.Fatalf("APIURL = %q, want the https endpoint", cfg.APIURL)
	}
	if cfg.WSSURL != "wss://socket.rubika.ir/socket" {
		t.Fatalf("WSSURL = %q, want the non-private wss endpoint", cfg.WSSURL)
	}
}

func TestValidWSSURLRejectsPrivateAndNonWSS(t *testing.T) {
	if validWSSURL("ws://socket.rubika.ir/s") {
		t.Fatal("ws:// must be rejected (only wss)")
	}
	if validWSSURL("wss://127.0.0.1/s") {
		t.Fatal("loopback wss must be rejected")
	}
	if validWSSURL("wss://169.254.169.254/s") {
		t.Fatal("metadata IP wss must be rejected")
	}
	if !validWSSURL("wss://socket.rubika.ir/s") {
		t.Fatal("legitimate wss must be accepted")
	}
}

// --- WebSocket read limit ---

func TestWebSocketReadLimit(t *testing.T) {
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		// Send a frame larger than MaxWSMessage; the client's read limit
		// must reject it.
		_ = c.WriteMessage(websocket.BinaryMessage, make([]byte, MaxWSMessage+1024))
	}))
	defer srv.Close()

	ws := NewWebSocket("ws"+strings.TrimPrefix(srv.URL, "http"), nil)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := ws.Dial(ctx); err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer ws.Close()
	if _, err := ws.Read(); err == nil {
		t.Fatal("Read of an oversized frame must fail (read limit)")
	}
}
