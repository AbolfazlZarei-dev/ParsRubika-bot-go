package transport

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkPostHotPath measures the per-request cost of the
// Transport.Post path with a tiny payload against a local server.
// This is the most-called function in any bot: every API method
// (SendMessage, GetUpdates, ...) funnels through it.
func BenchmarkPostHotPath(b *testing.B) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"OK","status_det":"OK","data":{}}`))
	}))
	defer srv.Close()

	t := New()
	ctx := context.Background()
	payload := map[string]any{
		"api_version": "6",
		"auth":        "abcdefghijklmnopqrstuvwxyz123456",
		"client": map[string]any{
			"app_name":    "Main",
			"app_version": "3.8.2",
			"platform":    "PWA",
			"package":     "m.rubika.ir",
			"lang_code":   "fa",
		},
		"data":   map[string]any{"text": "hello"},
		"method": "sendMessage",
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := t.Post(ctx, srv.URL, payload)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPostLargeBody measures the per-request cost with a larger
// payload (10 KiB message body), to verify the buffer-pool reuse
// actually pays off for non-trivial payloads.
func BenchmarkPostLargeBody(b *testing.B) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"OK","status_det":"OK","data":{}}`))
	}))
	defer srv.Close()

	t := New()
	ctx := context.Background()
	// 10 KiB of text — realistic for a long broadcast message.
	big := make([]byte, 10*1024)
	for i := range big {
		big[i] = 'x'
	}
	payload := map[string]any{
		"api_version": "6",
		"auth":        "abcdefghijklmnopqrstuvwxyz123456",
		"data":        map[string]any{"text": string(big)},
		"method":      "sendMessage",
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := t.Post(ctx, srv.URL, payload)
		if err != nil {
			b.Fatal(err)
		}
	}
}
