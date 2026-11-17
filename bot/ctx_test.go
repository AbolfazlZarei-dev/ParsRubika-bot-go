package bot

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

func TestSplitCommand(t *testing.T) {
	cases := []struct{ in, name, rest string }{
		{"/start", "start", ""},
		{"/Start@my_bot hello world", "start", "hello world"},
		{"!ban  42", "ban", " 42"},
		{"hello", "", ""},
		{"", "", ""},
	}
	for _, c := range cases {
		n, r := splitCommand(c.in)
		if n != c.name || r != c.rest {
			t.Errorf("splitCommand(%q) = %q,%q want %q,%q", c.in, n, r, c.name, c.rest)
		}
	}
}

func TestCtxArgsAndValuesFlowFromMiddleware(t *testing.T) {
	b := newTestBot(t)
	b.Use(func(ctx context.Context, bb *Bot, u *Update, next func() error) error {
		CtxFrom(ctx, bb, u).Set("user", "ali")
		return next()
	})
	var gotUser, gotCmd string
	var gotArgs []string
	b.Command("echo", func(c *Ctx) error {
		gotUser, gotCmd, gotArgs = c.GetString("user"), c.Command(), c.Args()
		return nil
	}, "say")
	_ = b.Dispatch(context.Background(), msgUpdate("b0x", "1", "!SAY a  b c"))
	if gotUser != "ali" || gotCmd != "say" || len(gotArgs) != 3 {
		t.Fatalf("user=%q cmd=%q args=%v", gotUser, gotCmd, gotArgs)
	}
}

func TestCtxReplyUsesReplyToMessageID(t *testing.T) {
	var mu sync.Mutex
	var sent map[string]any
	api := &fakeAPI{reply: func(m string, n int, w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		_ = json.Unmarshal(body, &sent)
		mu.Unlock()
		okJSON(w, map[string]any{"message_id": "r1"})
	}}
	srv := httptest.NewServer(api)
	defer srv.Close()
	b := newTestBot(t, WithBaseURL(srv.URL+"/%s/"))
	b.Message(func(c *Ctx) error { _, err := c.Reply("pong"); return err })
	_ = b.Dispatch(context.Background(), msgUpdate("b0chat", "m77", "ping"))
	mu.Lock()
	defer mu.Unlock()
	if sent["chat_id"] != "b0chat" || sent["text"] != "pong" || sent["reply_to_message_id"] != "m77" {
		t.Fatalf("payload = %v", sent)
	}
}

func TestCopyMessageRoutesByContent(t *testing.T) {
	var mu sync.Mutex
	var methods []string
	api := &fakeAPI{reply: func(m string, n int, w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		methods = append(methods, m)
		mu.Unlock()
		okJSON(w, map[string]any{"message_id": "x"})
	}}
	srv := httptest.NewServer(api)
	defer srv.Close()
	b := newTestBot(t, WithBaseURL(srv.URL+"/%s/"))
	ctx := context.Background()
	_, _ = b.CopyMessage(ctx, &models.Message{Text: "t"}, "c")
	_, _ = b.CopyMessage(ctx, &models.Message{Sticker: &models.Sticker{StickerID: "s"}}, "c")
	_, _ = b.CopyMessage(ctx, &models.Message{Location: &models.Location{Latitude: "1.5", Longitude: "2.5"}}, "c")
	if _, err := b.CopyMessage(ctx, &models.Message{}, "c"); err != ErrNothingToCopy {
		t.Fatalf("empty message: %v", err)
	}
	if _, err := b.CopyMessage(ctx, nil, "c"); err != ErrNoMessage {
		t.Fatalf("nil message: %v", err)
	}
	want := []string{"sendMessage", "sendSticker", "sendLocation"}
	for i, w := range want {
		if i >= len(methods) || methods[i] != w {
			t.Fatalf("methods = %v", methods)
		}
	}
}

func TestMeIsCached(t *testing.T) {
	api := &fakeAPI{reply: func(m string, n int, w http.ResponseWriter, r *http.Request) {
		okJSON(w, map[string]any{"bot": map[string]any{"bot_id": "b1", "bot_title": "Demo", "username": "demo_bot"}})
	}}
	srv := httptest.NewServer(api)
	defer srv.Close()
	b := newTestBot(t, WithBaseURL(srv.URL+"/%s/"))
	for i := 0; i < 3; i++ {
		if u, err := b.Username(context.Background()); err != nil || u != "demo_bot" {
			t.Fatalf("username=%q err=%v", u, err)
		}
	}
	if api.calls["getMe"] != 1 {
		t.Fatalf("getMe called %d times", api.calls["getMe"])
	}
}

func TestStateManagerCloseStopsGC(t *testing.T) {
	sm := NewStateManager(10_000_000) // 10ms
	sm.Set("u", "s")
	sm.Close()
	sm.Close() // idempotent
}
