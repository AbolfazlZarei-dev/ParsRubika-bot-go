// Package http3 is an optional HTTP/3 (QUIC) client for ParsRubika.
//
// It is a separate module so that the core library stays buildable with
// Go 1.20; quic-go itself needs Go 1.23. Pass the result of NewHTTP3Client
// to bot.WithHTTPClient or account.WithHTTPClient.
package http3

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	qh3 "github.com/quic-go/quic-go/http3"
)

// http3Once guards the lazy construction of the shared HTTP/3
// round-tripper. QUIC connections are expensive to set up (TLS 1.3
// handshake + 0-RTT), so we cache a single *qh3.Transport for the
// whole process and reuse its connection pool across calls.
var (
	http3Once sync.Once
	http3RT   *qh3.Transport
)

// sharedHTTP3Transport returns the process-wide *qh3.Transport. The
// transport is created lazily on first use so importing the package
// never pays the QUIC setup cost unless HTTP/3 is actually used.
func sharedHTTP3Transport() *qh3.Transport {
	http3Once.Do(func() {
		http3RT = &qh3.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS13,
			},
		}
	})
	return http3RT
}

// NewHTTP3Client builds an *http.Client that uses HTTP/3 (QUIC) for
// requests to HTTP/3-capable servers. It is an opt-in alternative to
// the default HTTP/2 client; callers that want to try HTTP/3 should
// pass it via `WithHTTP3Client`.
//
// Note: the QUIC handshake requires a TLS config with at least
// `MinVersion: tls.VersionTLS13`.
func NewHTTP3Client(timeout time.Duration, insecure bool) *http.Client {
	rt := &qh3.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
			MinVersion:         tls.VersionTLS13,
		},
	}
	return &http.Client{
		Transport: rt,
		Timeout:   timeout,
	}
}

// maxResponseBody caps how many bytes of an HTTP/3 response body are
// read into memory, so a malicious server cannot exhaust RAM by
// streaming an enormous body.
const maxResponseBody = 16 << 20 // 16 MiB

// maxErrorSnippet caps how much of a response body is embedded in an
// error string, so echoed credentials never leak through an error.
const maxErrorSnippet = 512

func readBody3(r io.Reader) ([]byte, bool, error) {
	b, err := io.ReadAll(io.LimitReader(r, maxResponseBody+1))
	if err != nil {
		return nil, false, err
	}
	return b, len(b) > maxResponseBody, nil
}

func sanitizeForError3(body []byte) string {
	if len(body) > maxErrorSnippet {
		body = body[:maxErrorSnippet]
	}
	s := string(body)
	if i := strings.Index(s, "-----BEGIN"); i >= 0 {
		s = s[:i] + "[redacted pem]"
	}
	return s
}

// Post3 sends a POST request over HTTP/3. It uses a process-wide cached
// *qh3.Transport so QUIC connections are reused across calls
// (re-dialing + re-handshaking per call would defeat the main benefit
// of HTTP/3).
func Post3(ctx context.Context, url string, headers map[string]string, body []byte, timeout time.Duration) ([]byte, int, error) {
	rt := sharedHTTP3Transport()
	client := &http.Client{
		Transport: rt,
		// Rubika endpoints never redirect; disable cross-host redirects so
		// the auth-bearing JSON body is never re-sent to a different host.
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if timeout > 0 {
		client.Timeout = timeout
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("transport: http3: %w", err)
	}
	defer resp.Body.Close()
	respBody, tooLarge, err := readBody3(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("transport: http3: read response: %w", err)
	}
	if tooLarge {
		return nil, resp.StatusCode, fmt.Errorf("transport: http3: response body exceeds %d bytes", maxResponseBody)
	}
	if resp.StatusCode >= 400 {
		return respBody, resp.StatusCode, fmt.Errorf("transport: http3 http %d: %s", resp.StatusCode, sanitizeForError3(respBody))
	}
	return respBody, resp.StatusCode, nil
}
