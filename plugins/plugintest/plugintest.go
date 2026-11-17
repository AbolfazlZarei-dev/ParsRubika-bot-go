// Package plugintest provides a ready-made test harness for plugins: a
// bot wired to an in-process fake Rubika API that records every call.
//
//	env := plugintest.New(t)
//	env.Manager.MustRegister(myplugin.New())
//	env.Enable("myplugin")
//	env.Say("b0user", "u1", "/hello")
//	calls := env.Calls("sendMessage")
package plugintest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

// Call is one recorded API request.
type Call struct {
	Method  string
	Payload map[string]any
}

// Responder customises the answer of the fake API for a method. It
// returns the "data" object (nil means an empty object).
type Responder func(payload map[string]any) map[string]any

// Env is the test environment.
type Env struct {
	T       *testing.T
	Bot     *bot.Bot
	Manager *plugins.Manager
	Server  *httptest.Server

	mu         sync.Mutex
	calls      []Call
	responders map[string]Responder
	seq        atomic.Int64
}

type silent struct{}

func (silent) Printf(string, ...any) {}

// New builds an Env and registers its cleanup with t.
func New(t *testing.T, opts ...bot.Option) *Env {
	t.Helper()
	e := &Env{T: t, responders: map[string]Responder{}}
	e.Server = httptest.NewServer(http.HandlerFunc(e.serve))
	base := []bot.Option{
		bot.WithBaseURL(e.Server.URL + "/%s/"), bot.WithLogger(silent{}), bot.WithPersistOffset(false),
		bot.WithRateLimit(0), bot.WithRetries(0, time.Millisecond), bot.WithConcurrency(1),
	}
	e.Bot = bot.New("test-token", append(base, opts...)...)
	e.Manager = plugins.New(e.Bot)
	t.Cleanup(func() {
		_ = e.Manager.DisableAll(context.Background())
		e.Server.Close()
	})
	return e
}

func (e *Env) serve(w http.ResponseWriter, r *http.Request) {
	method := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	body, _ := io.ReadAll(r.Body)
	var payload map[string]any
	_ = json.Unmarshal(body, &payload)
	e.mu.Lock()
	e.calls = append(e.calls, Call{Method: method, Payload: payload})
	resp := e.responders[method]
	e.mu.Unlock()
	data := map[string]any{"message_id": fmt.Sprint("sent", e.seq.Add(1))}
	if resp != nil {
		if d := resp(payload); d != nil {
			data = d
		} else {
			data = map[string]any{}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "OK", "data": data})
}

// Respond installs a custom answer for an API method.
func (e *Env) Respond(method string, r Responder) {
	e.mu.Lock()
	e.responders[method] = r
	e.mu.Unlock()
}

// Calls returns the recorded calls of one method.
func (e *Env) Calls(method string) []Call {
	e.mu.Lock()
	defer e.mu.Unlock()
	var out []Call
	for _, c := range e.calls {
		if c.Method == method {
			out = append(out, c)
		}
	}
	return out
}

// Reset forgets every recorded call.
func (e *Env) Reset() {
	e.mu.Lock()
	e.calls = nil
	e.mu.Unlock()
}

// Sent returns the texts of all sendMessage calls.
func (e *Env) Sent() []string {
	var out []string
	for _, c := range e.Calls("sendMessage") {
		s, _ := c.Payload["text"].(string)
		out = append(out, s)
	}
	return out
}

// Enable activates the named plugins and fails the test on error.
func (e *Env) Enable(names ...string) {
	e.T.Helper()
	if err := e.Manager.Enable(context.Background(), names...); err != nil {
		e.T.Fatalf("enable %v: %v", names, err)
	}
}

// Say dispatches a text message from sender in chat.
func (e *Env) Say(chatID, senderID, text string) error {
	id := fmt.Sprint("m", e.seq.Add(1))
	return e.Bot.Dispatch(context.Background(), &models.Update{
		Type: models.UpdateTypeNewMessage, ChatID: chatID, UpdateTime: id,
		NewMessage: &models.Message{MessageID: id, Text: text, SenderID: senderID},
	})
}

// Press dispatches an inline button press.
func (e *Env) Press(chatID, senderID, buttonID string) error {
	id := fmt.Sprint("m", e.seq.Add(1))
	return e.Bot.Dispatch(context.Background(), &models.Update{
		Type: models.UpdateTypeNewMessage, ChatID: chatID, UpdateTime: id,
		NewMessage: &models.Message{MessageID: id, SenderID: senderID, AuxData: &models.AuxData{ButtonID: buttonID}},
	})
}

// Event dispatches a non-message update (StartedBot, StoppedBot, ...).
func (e *Env) Event(t models.UpdateType, chatID string) error {
	return e.Bot.Dispatch(context.Background(), &models.Update{Type: t, ChatID: chatID})
}
