// Package httpx is a tiny shared HTTP shim for the provider packages
// whose wire format is NOT OpenAI chat-completions (Anthropic Messages,
// Google Gemini generateContent, GapGPT Responses). The OpenAI-shaped
// providers use the dedicated oai package instead.
//
// It does one thing: marshal a JSON body, POST it with the given
// headers, cap the response at 8 MiB and hand back the status code and
// raw bytes. Each provider parses its own response shape.
package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MaxBodyBytes caps a response body so a misbehaving endpoint cannot
// exhaust memory.
const MaxBodyBytes = 8 << 20

// PostJSON POSTs the marshalled body to url with the given headers and
// returns the HTTP status code and the raw response body. A nil client
// falls back to http.DefaultClient. The Content-Type header defaults to
// application/json when not set in headers.
func PostJSON(ctx context.Context, client *http.Client, url string, headers map[string]string, body any) (int, []byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return 0, nil, fmt.Errorf("httpx: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return 0, nil, fmt.Errorf("httpx: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("httpx: request: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, MaxBodyBytes))
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("httpx: read body: %w", err)
	}
	return resp.StatusCode, raw, nil
}

// Truncate returns at most n bytes of b as a string, with an ellipsis
// when it was shortened. Used to keep error messages readable.
func Truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "…"
}
