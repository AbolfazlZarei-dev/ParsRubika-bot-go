package account

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// The upload service has been seen to answer with the hash but no
// status/status_det — the exact shape that makes the stock
// files.Uploader fail with `upload failed: /`.
func TestPushImageChunksAcceptsHashWithoutStatus(t *testing.T) {
	var gotFileID, gotHashSend string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotFileID = r.Header.Get("file-id")
		gotHashSend = r.Header.Get("access-hash-send")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"byte_receive":1,"access_hash_rec":"hash-no-status"}}`))
	}))
	defer srv.Close()

	hash, err := pushImageChunks(context.Background(), srv.Client(), srv.URL,
		"auth", "ua", "f1", "send1", bytes.Repeat([]byte{7}, 300))
	if err != nil {
		t.Fatalf("pushImageChunks: %v", err)
	}
	if hash != "hash-no-status" {
		t.Fatalf("hash = %q, want hash-no-status", hash)
	}
	if gotFileID != "f1" || gotHashSend != "send1" {
		t.Fatalf("headers: file-id=%q access-hash-send=%q", gotFileID, gotHashSend)
	}
}

// Numeric status/status_det and a numeric hash must also parse.
func TestPushImageChunksAcceptsNumericFields(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":200,"status_det":200,"data":{"access_hash_rec":123456789}}`))
	}))
	defer srv.Close()

	hash, err := pushImageChunks(context.Background(), srv.Client(), srv.URL,
		"auth", "ua", "f1", "send1", []byte("hello"))
	if err != nil {
		t.Fatalf("pushImageChunks: %v", err)
	}
	if hash != "123456789" {
		t.Fatalf("hash = %q, want 123456789", hash)
	}
}

// ERROR_TRY_AGAIN must surface as errTryAgain so the caller re-requests
// the slot instead of giving up.
func TestPushImageChunksSurfacesTryAgain(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ERROR_TRY_AGAIN","status_det":"ERROR_TRY_AGAIN"}`))
	}))
	defer srv.Close()

	_, err := pushImageChunks(context.Background(), srv.Client(), srv.URL,
		"auth", "ua", "f1", "send1", []byte("hello"))
	if !errors.Is(err, errTryAgain) {
		t.Fatalf("err = %v, want errTryAgain", err)
	}
}

// A final chunk with no hash anywhere is a hard error.
func TestPushImageChunksFailsWithoutHash(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"OK","status_det":"OK"}`))
	}))
	defer srv.Close()

	_, err := pushImageChunks(context.Background(), srv.Client(), srv.URL,
		"auth", "ua", "f1", "send1", []byte("hello"))
	if err == nil {
		t.Fatal("want an error when no access_hash_rec is returned")
	}
}

// avatarJSONStr must render every JSON scalar kind the upload service
// has been seen to use.
func TestAvatarJSONStr(t *testing.T) {
	cases := map[string]struct {
		in   any
		want string
	}{
		"string":      {"OK", "OK"},
		"jsonNumber":  {json.Number("129951015778986"), "129951015778986"},
		"float":       {float64(484), "484"},
		"bool":        {true, "true"},
		"nil":         {nil, ""},
		"unsupported": {[]any{1}, ""},
	}
	for name, c := range cases {
		if got := avatarJSONStr(c.in); got != c.want {
			t.Errorf("%s: avatarJSONStr = %q, want %q", name, got, c.want)
		}
	}
}

// avatarHashRec must find the hash at the top level and one level of
// data nesting.
func TestAvatarHashRecNesting(t *testing.T) {
	top := map[string]any{"access_hash_rec": "abc"}
	if got := avatarHashRec(top); got != "abc" {
		t.Fatalf("top: got %q", got)
	}
	nested := map[string]any{"data": map[string]any{"access_hash_rec": json.Number("42")}}
	if got := avatarHashRec(nested); got != "42" {
		t.Fatalf("nested: got %q", got)
	}
	if got := avatarHashRec(nil); got != "" {
		t.Fatalf("nil: got %q", got)
	}
}

// HangiProfile must refuse to run without a GUID when the account is
// not signed in.
func TestHangiProfileRequiresGUID(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer a.Close()
	if _, err := a.HangiProfile(context.Background(), "", "1.jpg", "2.jpg"); err == nil {
		t.Fatal("want an error when no GUID is available")
	}
}

// WithHangi* options must apply their values and keep the safe
// defaults for the zero value.
func TestHangiOptionsDefaults(t *testing.T) {
	o := &HangiOptions{Rounds: 1, Offset: 300 * time.Millisecond, Gap: 700 * time.Millisecond}
	WithHangiRounds(3)(o)
	WithHangiOffset(time.Second)(o)
	WithHangiGap(2 * time.Second)(o)
	WithHangiClean()(o)
	WithHangiKeepOld()(o)
	var msgs []string
	WithHangiProgress(func(msg string) { msgs = append(msgs, msg) })(o)
	if o.Rounds != 3 || o.Offset != time.Second || o.Gap != 2*time.Second || !o.Clean || !o.KeepOld {
		t.Fatalf("options not applied: %+v", o)
	}
	if len(msgs) != 0 {
		t.Fatalf("progress should be silent until used")
	}
	// A non-positive rounds override must not disable the run.
	WithHangiRounds(0)(o)
	if o.Rounds != 3 {
		t.Fatalf("WithHangiRounds(0) must be a no-op, got %d", o.Rounds)
	}
}
