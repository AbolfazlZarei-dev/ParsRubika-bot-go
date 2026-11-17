package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

type nullLogger struct{}

func (nullLogger) Printf(string, ...any) {}

func newTestBot(t *testing.T, opts ...Option) *Bot {
	t.Helper()
	base := []Option{WithPersistOffset(false), WithLogger(nullLogger{}), WithRateLimit(0), WithRetries(2, time.Millisecond), WithSkipBacklog(false)}
	return New("tok", append(base, opts...)...)
}

func msgUpdate(chat, id, text string) *models.Update {
	return &models.Update{
		Type: models.UpdateTypeNewMessage, ChatID: chat, UpdateTime: id,
		NewMessage: &models.Message{MessageID: id, Text: text},
	}
}

func TestMiddlewareCanBlockHandler(t *testing.T) {
	b := newTestBot(t)
	var ran atomic.Bool
	b.On(nil, func(context.Context, *Bot, *Update) error { ran.Store(true); return nil })
	b.Use(func(ctx context.Context, _ *Bot, _ *Update, next func() error) error { return nil })
	if err := b.Dispatch(context.Background(), msgUpdate("c", "1", "hi")); err != nil {
		t.Fatal(err)
	}
	if ran.Load() {
		t.Fatal("handler ran although the middleware did not call next")
	}
}

func TestMiddlewareWrapsHandlerAndSeesErrors(t *testing.T) {
	b := newTestBot(t)
	var order []string
	boom := errors.New("boom")
	b.Use(func(ctx context.Context, _ *Bot, _ *Update, next func() error) error {
		order = append(order, "mw-before")
		err := next()
		order = append(order, "mw-after:"+fmt.Sprint(err))
		return err
	})
	b.On(nil, func(context.Context, *Bot, *Update) error { order = append(order, "handler"); return boom })
	var seen error
	b.opts.ErrorHandler = func(ec *ErrorContext) { seen = ec.Err }
	_ = b.Dispatch(context.Background(), msgUpdate("c", "1", "x"))
	got := strings.Join(order, ",")
	if got != "mw-before,handler,mw-after:boom" {
		t.Fatalf("order = %s", got)
	}
	if !errors.Is(seen, boom) {
		t.Fatalf("error handler got %v", seen)
	}
}

func TestPanicIsRecovered(t *testing.T) {
	var got *ErrorContext
	b := newTestBot(t, WithErrorHandler(func(ec *ErrorContext) { got = ec }))
	b.On(nil, func(context.Context, *Bot, *Update) error { panic("kaboom") })
	if err := b.Dispatch(context.Background(), msgUpdate("c", "1", "x")); err != nil {
		t.Fatalf("Dispatch returned %v", err)
	}
	if got == nil || got.Panic != "kaboom" || len(got.Stack) == 0 {
		t.Fatalf("panic not reported: %+v", got)
	}
	if b.Snapshot().Panics != 1 {
		t.Fatal("panic counter not incremented")
	}
}

func TestErrContinueAndOff(t *testing.T) {
	b := newTestBot(t)
	var calls []string
	first := b.On(nil, func(context.Context, *Bot, *Update) error { calls = append(calls, "a"); return ErrContinue })
	b.On(nil, func(context.Context, *Bot, *Update) error { calls = append(calls, "b"); return nil })
	b.On(nil, func(context.Context, *Bot, *Update) error { calls = append(calls, "c"); return nil })
	_ = b.Dispatch(context.Background(), msgUpdate("c", "1", "x"))
	if strings.Join(calls, "") != "ab" {
		t.Fatalf("calls = %v", calls)
	}
	if !b.Off(first) || b.Off(first) {
		t.Fatal("Off should succeed once")
	}
	calls = nil
	_ = b.Dispatch(context.Background(), msgUpdate("c", "2", "x"))
	if strings.Join(calls, "") != "b" {
		t.Fatalf("calls after Off = %v", calls)
	}
}

func TestDuplicateUpdatesAreDropped(t *testing.T) {
	b := newTestBot(t)
	var n atomic.Int32
	b.On(nil, func(context.Context, *Bot, *Update) error { n.Add(1); return nil })
	for i := 0; i < 3; i++ {
		_ = b.Dispatch(context.Background(), msgUpdate("c", "same", "x"))
	}
	if n.Load() != 1 || b.Snapshot().Duplicate != 2 {
		t.Fatalf("handled=%d dup=%d", n.Load(), b.Snapshot().Duplicate)
	}
}

func TestSeenSetEvictsOldest(t *testing.T) {
	s := newSeenSet(16)
	for i := 0; i < 40; i++ {
		s.add(fmt.Sprint("k", i))
	}
	if s.has("k0") || !s.has("k39") {
		t.Fatal("ring eviction is wrong")
	}
	if s.add("k39") {
		t.Fatal("duplicate accepted")
	}
}

func TestDispatcherKeepsPerKeyOrderAndRunsKeysInParallel(t *testing.T) {
	d := newDispatcher(4)
	var mu sync.Mutex
	order := map[string][]int{}
	var running, peak int32
	for i := 0; i < 20; i++ {
		for _, k := range []string{"a", "b", "c"} {
			i, k := i, k
			d.submit(k, func() {
				cur := atomic.AddInt32(&running, 1)
				for {
					p := atomic.LoadInt32(&peak)
					if cur <= p || atomic.CompareAndSwapInt32(&peak, p, cur) {
						break
					}
				}
				time.Sleep(time.Millisecond)
				atomic.AddInt32(&running, -1)
				mu.Lock()
				order[k] = append(order[k], i)
				mu.Unlock()
			})
		}
	}
	if !d.wait(5 * time.Second) {
		t.Fatal("dispatcher did not drain")
	}
	for k, v := range order {
		for i := range v {
			if v[i] != i {
				t.Fatalf("key %s out of order: %v", k, v)
			}
		}
	}
	if atomic.LoadInt32(&peak) < 2 {
		t.Fatal("different keys never ran in parallel")
	}
}

func TestParseUpdateLifecycleTypes(t *testing.T) {
	for _, typ := range []string{"StartedBot", "StoppedBot"} {
		u, err := parseUpdate(map[string]any{"type": typ, "chat_id": "c1"}, nil)
		if err != nil || u == nil || string(u.Type) != typ || u.ChatID != "c1" {
			t.Fatalf("%s: %v %+v", typ, err, u)
		}
	}
	if _, err := parseUpdate(map[string]any{"type": "Nonsense"}, nil); err == nil {
		t.Fatal("unknown type must error")
	}
}

// fakeAPI is a tiny Rubika Bot API double.
type fakeAPI struct {
	mu    sync.Mutex
	calls map[string]int
	reply func(method string, n int, w http.ResponseWriter, r *http.Request)
}

func (f *fakeAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	f.mu.Lock()
	if f.calls == nil {
		f.calls = map[string]int{}
	}
	f.calls[method]++
	n := f.calls[method]
	f.mu.Unlock()
	f.reply(method, n, w, r)
}

func okJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "OK", "data": data})
}

func TestRetriesTransientHTTPErrorsHonouringRetryAfter(t *testing.T) {
	api := &fakeAPI{reply: func(m string, n int, w http.ResponseWriter, r *http.Request) {
		if n < 3 {
			w.Header().Set("Retry-After", "0")
			http.Error(w, "slow down", http.StatusTooManyRequests)
			return
		}
		okJSON(w, map[string]any{"bot": map[string]any{"bot_id": "b1"}})
	}}
	srv := httptest.NewServer(api)
	defer srv.Close()
	b := newTestBot(t, WithBaseURL(srv.URL+"/%s/"))
	if _, err := b.makeRequest(context.Background(), "getMe", map[string]any{}); err != nil {
		t.Fatalf("expected success after retries: %v", err)
	}
	if api.calls["getMe"] != 3 {
		t.Fatalf("calls = %d", api.calls["getMe"])
	}
}

func TestDoesNotRetryPermanentErrors(t *testing.T) {
	api := &fakeAPI{reply: func(m string, n int, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"INVALID_INPUT","dev_message":"bad chat"}`))
	}}
	srv := httptest.NewServer(api)
	defer srv.Close()
	b := newTestBot(t, WithBaseURL(srv.URL+"/%s/"))
	_, err := b.makeRequest(context.Background(), "sendMessage", map[string]any{})
	if !IsStatus(err, "INVALID_INPUT") || api.calls["sendMessage"] != 1 {
		t.Fatalf("err=%v calls=%d", err, api.calls["sendMessage"])
	}
	// wrapped errors must still match
	if !IsStatus(fmt.Errorf("ctx: %w", err), "INVALID_INPUT") {
		t.Fatal("IsStatus must see through wrapping")
	}
}

func TestPollingEndToEnd(t *testing.T) {
	api := &fakeAPI{reply: func(m string, n int, w http.ResponseWriter, r *http.Request) {
		if n == 1 {
			okJSON(w, map[string]any{
				"updates": []any{map[string]any{
					"type": "NewMessage", "chat_id": "c1",
					"new_message": map[string]any{"message_id": "m1", "text": "salam"},
				}},
				"next_offset_id": "o1",
			})
			return
		}
		okJSON(w, map[string]any{"updates": []any{}, "next_offset_id": "o1"})
	}}
	srv := httptest.NewServer(api)
	defer srv.Close()

	b := newTestBot(t, WithBaseURL(srv.URL+"/%s/"), WithPollInterval(5*time.Millisecond))
	got := make(chan string, 1)
	b.On(nil, func(_ context.Context, _ *Bot, u *Update) error { got <- u.NewMessage.Text; return nil })
	started := make(chan struct{})
	b.OnStart(func(context.Context, *Bot) error { close(started); return nil })
	stopped := make(chan struct{})
	b.OnStop(func(context.Context, *Bot) error { close(stopped); return nil })

	done := make(chan error, 1)
	go func() { done <- b.Start(context.Background()) }()
	select {
	case text := <-got:
		if text != "salam" {
			t.Fatalf("text = %q", text)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("update never reached the handler")
	}
	b.Stop()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Start did not return after Stop")
	}
	<-started
	<-stopped
	if b.IsRunning() {
		t.Fatal("IsRunning must be false after Start returns")
	}
}

func TestSkipBacklogDropsFirstBatchOnly(t *testing.T) {
	mk := func(id string) map[string]any {
		return map[string]any{"type": "NewMessage", "chat_id": "c", "new_message": map[string]any{"message_id": id, "text": id}}
	}
	api := &fakeAPI{reply: func(m string, n int, w http.ResponseWriter, r *http.Request) {
		okJSON(w, map[string]any{"updates": []any{mk(fmt.Sprint("m", n))}, "next_offset_id": fmt.Sprint("o", n)})
	}}
	srv := httptest.NewServer(api)
	defer srv.Close()
	b := newTestBot(t, WithBaseURL(srv.URL+"/%s/"), WithSkipBacklog(true))
	first, err := b.Updater(context.Background(), 10, 0)
	if err != nil || len(first) != 0 {
		t.Fatalf("first batch should be skipped: %v %v", first, err)
	}
	second, err := b.Updater(context.Background(), 10, 0)
	if err != nil || len(second) != 1 {
		t.Fatalf("second batch should be delivered: %v %v", second, err)
	}
}

func TestWebhookSecretAndBodyLimit(t *testing.T) {
	b := newTestBot(t, WithWebhookSecret("s3cret"), func() Option { return Options{MaxBodyBytes: 200} }())
	var n atomic.Int32
	b.On(nil, func(context.Context, *Bot, *Update) error { n.Add(1); return nil })
	h := b.WebhookHandler()

	post := func(url, body string) int {
		req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(body))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec.Code
	}
	upd := `{"update":{"type":"NewMessage","chat_id":"c","new_message":{"message_id":"w1","text":"x"}}}`
	if c := post("/webhook/receiveUpdate", upd); c != http.StatusForbidden {
		t.Fatalf("missing secret: %d", c)
	}
	if c := post("/webhook/receiveUpdate?s=s3cret", strings.Repeat("a", 500)); c != http.StatusBadRequest {
		t.Fatalf("oversized body: %d", c)
	}
	if c := post("/webhook/receiveUpdate?s=s3cret", upd); c != http.StatusOK {
		t.Fatalf("valid request: %d", c)
	}
	b.disp.wait(2 * time.Second)
	if n.Load() != 1 {
		t.Fatalf("handled %d updates", n.Load())
	}
}

func TestBackoffIsCappedAndJittered(t *testing.T) {
	b := newTestBot(t, WithRetries(5, 100*time.Millisecond))
	b.opts.MaxBackoff = time.Second
	for attempt := 0; attempt < 64; attempt++ {
		d := b.backoffDelay(attempt)
		if d <= 0 || d > time.Second {
			t.Fatalf("attempt %d: delay %v out of range", attempt, d)
		}
	}
}

func TestOptionsStructStillWorks(t *testing.T) {
	b := New("t", Options{PollInterval: 7 * time.Millisecond, PersistOffset: false, MaxRetries: 9})
	o := b.Options()
	if o.PollInterval != 7*time.Millisecond || o.MaxRetries != 9 || o.Concurrency != DefaultConcurrency {
		t.Fatalf("unexpected options: %+v", o)
	}
}

func TestHandlerPriorityAndFallbackOrder(t *testing.T) {
	b := newTestBot(t)
	var got []string
	b.Fallback(func(c *Ctx) error { got = append(got, "fallback"); return nil })
	b.Handle(nil, func(c *Ctx) error { got = append(got, "normal"); return ErrContinue })
	b.OnPriority(PriorityPlugin, nil, func(context.Context, *Bot, *Update) error {
		got = append(got, "plugin")
		return ErrContinue
	})
	_ = b.Dispatch(context.Background(), msgUpdate("c", "1", "x"))
	if strings.Join(got, ",") != "plugin,normal,fallback" {
		t.Fatalf("order = %v", got)
	}
}
