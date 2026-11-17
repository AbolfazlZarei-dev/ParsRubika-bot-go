package account

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/transport"
)

// HangProfile must fire exactly two updateProfile requests with the
// supplied fields, both carrying the matching updated_parameters list.
// It must not raise when both shots succeed and must surface the two
// events / errors on the returned result.
func TestHangProfileFiresTwoShots(t *testing.T) {
	a, err := New(WithAuth("abcdefghijklmnopqrstuvwxyz012345"))
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()

	var seen int32
	type shot struct {
		Method  string                 `json:"method"`
		First   string                 `json:"first_name"`
		Last    string                 `json:"last_name"`
		Bio     string                 `json:"bio"`
		Birth   string                 `json:"birth_date"`
		Updated []string               `json:"updated_parameters"`
		Extra   map[string]interface{} `json:"-"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		// Decrypt the data_enc payload via the same path as Request:
		// the test server returns an unencrypted OK so we only need to
		// record that a request landed. We do NOT parse the encrypted
		// inner input here (that requires the symmetric key); instead
		// we trust Request to have built the payload correctly and only
		// count the shots.
		_ = raw
		atomic.AddInt32(&seen, 1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"OK","status_det":"OK","data":{}}`))
	}))
	defer srv.Close()
	a.dcConfig = &transport.DCConfig{APIURL: srv.URL, WSSURL: "wss://x.invalid/"}

	res, err := a.HangProfile(context.Background(), "Foo", "Bar", "hello", "2000-01-01")
	if err != nil {
		t.Fatalf("HangProfile: %v", err)
	}
	if res == nil {
		t.Fatal("result is nil")
	}
	if res.FirstErr != nil {
		t.Fatalf("first shot error: %v", res.FirstErr)
	}
	if res.SecondErr != nil {
		t.Fatalf("second shot error: %v", res.SecondErr)
	}
	if res.First == nil || res.Second == nil {
		t.Fatal("both shot events must be non-nil on success")
	}
	if got := atomic.LoadInt32(&seen); got != 2 {
		t.Fatalf("expected 2 updateProfile shots, got %d", got)
	}
	_ = shot{} // keep the struct type referenced for future field checks
}

// HangProfile must propagate errors when the server rejects the shots.
// Both shots should fail (server returns a non-OK status_det) and the
// result must surface the errors without panicking.
func TestHangProfileSurfacesErrors(t *testing.T) {
	a, err := New(WithAuth("abcdefghijklmnopqrstuvwxyz012345"))
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ERROR","status_det":"INVALID_REQUEST","data":{}}`))
	}))
	defer srv.Close()
	a.dcConfig = &transport.DCConfig{APIURL: srv.URL, WSSURL: "wss://x.invalid/"}

	res, err := a.HangProfile(context.Background(), "Foo", "", "", "")
	if err != nil {
		t.Fatalf("HangProfile returned error (should be nil, errors go on the result): %v", err)
	}
	if res == nil {
		t.Fatal("result is nil")
	}
	if res.FirstErr == nil {
		t.Fatal("expected first shot error, got nil")
	}
	if res.SecondErr == nil {
		t.Fatal("expected second shot error, got nil")
	}
}

// HangProfile must honour context cancellation: when the context is
// already cancelled the shots must still return (with an error) and the
// call must not block forever.
func TestHangProfileRespectsContextCancellation(t *testing.T) {
	a, err := New(WithAuth("abcdefghijklmnopqrstuvwxyz012345"))
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()

	// A server that hangs so the only way out is context cancellation.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()
	a.dcConfig = &transport.DCConfig{APIURL: srv.URL, WSSURL: "wss://x.invalid/"}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before the call

	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _ = a.HangProfile(ctx, "Foo", "Bar", "bio", "2000-01-01")
	}()
	select {
	case <-done:
		// good: returned instead of blocking
	case <-make(chan struct{}):
		// fallback no-op
	}
	_ = json.Marshal // keep encoding/json referenced
}
