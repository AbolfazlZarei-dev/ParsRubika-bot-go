package account

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/security"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/transport"
)

// --- MessageQueue: concurrent Submit + Stop must not panic ---

func TestMessageQueueConcurrentSubmitStopNoPanic(t *testing.T) {
	a, err := New(WithAuth("abcdefghijklmnopqrstuvwxyz012345"))
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	// Point the transport at a local server that fails fast so workers
	// do not hit the real network.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()
	a.dcConfig = &transport.DCConfig{APIURL: srv.URL, WSSURL: "wss://x.invalid/"}

	mq := NewMessageQueue(a, 64, 4)
	var wg sync.WaitGroup
	stop := make(chan struct{})
	// Hammer Submit from many goroutines while Stop runs concurrently.
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				res := mq.Submit("guid", "hi")
				select {
				case <-res:
				case <-time.After(50 * time.Millisecond):
				}
			}
		}()
	}
	// Let some submits land, then stop.
	time.Sleep(30 * time.Millisecond)
	close(stop)
	mq.Stop()
	wg.Wait()
	// Submit after Stop must return ErrQueueClosed, not panic.
	res := mq.Submit("guid", "late")
	select {
	case e := <-res:
		if e == nil {
			t.Fatal("post-Stop Submit should fail with ErrQueueClosed")
		}
	case <-time.After(time.Second):
		t.Fatal("post-Stop Submit blocked")
	}
}

// --- Public Uploader/Downloader must not be built with a nil HTTP client ---

func TestPublicUploaderDownloaderNonNil(t *testing.T) {
	a, err := New(WithAuth("abcdefghijklmnopqrstuvwxyz012345"))
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	if a.Uploader() == nil {
		t.Fatal("public Uploader() is nil")
	}
	if a.Downloader() == nil {
		t.Fatal("public Downloader() is nil")
	}
	// The transport's HTTP client is what the uploader/downloader should
	// now be wired with (previously nil → nil-panic on Upload/Download).
	if a.transport == nil || a.transport.HTTPClient() == nil {
		t.Fatal("transport HTTP client is nil; uploader/downloader would nil-panic")
	}
}

// --- RequestOptions.Auth must be honoured on the v6 non-tmpSession path ---

func TestRequestOptionsAuthHonoured(t *testing.T) {
	a, err := New(WithAuth("abcdefghijklmnopqrstuvwxyz012345"))
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()

	var captured struct {
		Auth string `json:"auth"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &captured)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"OK","status_det":"OK","data":{}}`))
	}))
	defer srv.Close()
	a.dcConfig = &transport.DCConfig{APIURL: srv.URL, WSSURL: "wss://x.invalid/"}

	const customAuth = "zyxwvutsrqponmlkjihgfedcba987654"
	_, err = a.Request(context.Background(), RequestOptions{
		Method:  "getUserInfo",
		Auth:    customAuth,
		Encrypt: false,
	})
	if err != nil {
		t.Fatalf("Request: %v", err)
	}
	want := security.DecodeAuth(customAuth)
	if captured.Auth != want {
		t.Fatalf("payload auth = %q, want %q (DecodeAuth of opts.Auth)", captured.Auth, want)
	}
}

// --- Stop/Close set the stopped flag so Listen will not redial ---

func TestStopSetsStoppedFlag(t *testing.T) {
	a, err := New(WithAuth("abcdefghijklmnopqrstuvwxyz012345"))
	if err != nil {
		t.Fatal(err)
	}
	if a.stopped.Load() {
		t.Fatal("stopped flag should be false on a fresh account")
	}
	a.Stop()
	if !a.stopped.Load() {
		t.Fatal("Stop() must set the stopped flag")
	}
}

func TestCloseSetsStoppedFlag(t *testing.T) {
	a, err := New(WithAuth("abcdefghijklmnopqrstuvwxyz012345"))
	if err != nil {
		t.Fatal(err)
	}
	a.Close()
	if !a.stopped.Load() {
		t.Fatal("Close() must set the stopped flag")
	}
}
