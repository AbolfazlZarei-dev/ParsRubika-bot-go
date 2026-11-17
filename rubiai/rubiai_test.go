package rubiai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// newTestClient builds a RubiaAI pointed at an httptest server.
func newTestClient(t *testing.T, srv *httptest.Server) *RubiaAI {
	t.Helper()
	r := New(WithAuth("test-auth-abcdefghijklmnopqrstuvwxyz12"))
	r.apiURL = srv.URL
	r.SetHTTPClient(srv.Client())
	return r
}

// handler builds an http.HandlerFunc that routes on the "method" field
// of the inbound JSON request body.
func handler(routes map[string]func(input map[string]any) (status, statusDet string, output map[string]any)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		data, _ := payload["data"].(map[string]any)
		method, _ := data["method"].(string)
		input, _ := data["input"].(map[string]any)
		route, ok := routes[method]
		if !ok {
			http.Error(w, "no route "+method, http.StatusBadRequest)
			return
		}
		st, stDet, out := route(input)
		if st == "" {
			st = "OK"
		}
		if stDet == "" {
			stDet = "OK"
		}
		body := map[string]any{
			"data": map[string]any{
				"status":     st,
				"status_det": stDet,
				"output":     out,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(body)
	}
}

func TestGetWelcomeValidatesLocation(t *testing.T) {
	r := New(WithAuth("x"))
	if _, err := r.GetWelcome(context.Background(), "bogus", ""); err == nil {
		t.Fatal("expected validation error for bad location")
	}
	srv := httptest.NewServer(handler(map[string]func(map[string]any) (string, string, map[string]any){
		"getWelcome": func(in map[string]any) (string, string, map[string]any) {
			if in["location"] != "Messenger" {
				t.Errorf("location not normalised: %v", in["location"])
			}
			return "", "", map[string]any{"welcome": "hello"}
		},
	}))
	defer srv.Close()
	r = newTestClient(t, srv)
	ev, err := r.GetWelcome(context.Background(), "messenger", "abc")
	if err != nil {
		t.Fatalf("GetWelcome: %v", err)
	}
	if ev.AnyStr("welcome") != "hello" {
		t.Fatalf("output not returned: %v", ev)
	}
}

func TestGetServiceAccess(t *testing.T) {
	srv := httptest.NewServer(handler(map[string]func(map[string]any) (string, string, map[string]any){
		"getServiceAccess": func(map[string]any) (string, string, map[string]any) {
			return "", "", map[string]any{"quota": 1000}
		},
	}))
	defer srv.Close()
	r := newTestClient(t, srv)
	ev, err := r.GetServiceAccess(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if ev.Int("quota") != 1000 {
		t.Fatalf("quota not returned: %v", ev.Map())
	}
}

func TestGetSessionsPagination(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(handler(map[string]func(map[string]any) (string, string, map[string]any){
		"getSessions": func(in map[string]any) (string, string, map[string]any) {
			n := atomic.AddInt32(&calls, 1)
			start, _ := in["start_id"].(string)
			if n == 1 {
				if start != "" {
					t.Errorf("first call start_id should be empty, got %q", start)
				}
				return "", "", map[string]any{
					"sessions":      []any{map[string]any{"session_id": "s1"}, map[string]any{"session_id": "s2"}},
					"has_continue":  true,
					"next_start_id": "cursor1",
				}
			}
			if start != "cursor1" {
				t.Errorf("second call start_id should be cursor1, got %q", start)
			}
			return "", "", map[string]any{
				"sessions":     []any{map[string]any{"session_id": "s3"}},
				"has_continue": false,
			}
		},
	}))
	defer srv.Close()
	r := newTestClient(t, srv)
	ev, err := r.GetSessions(context.Background(), true, 10, "")
	if err != nil {
		t.Fatal(err)
	}
	sessions := ev.Slice("sessions")
	if len(sessions) != 3 {
		t.Fatalf("expected 3 sessions across pages, got %d: %v", len(sessions), sessions)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Fatalf("expected 2 paginated calls, got %d", calls)
	}
}

func TestGetChatHistoryPagination(t *testing.T) {
	srv := httptest.NewServer(handler(map[string]func(map[string]any) (string, string, map[string]any){
		"getChatHistory": func(in map[string]any) (string, string, map[string]any) {
			if in["session_id"] != "sess" {
				t.Errorf("session_id not forwarded: %v", in["session_id"])
			}
			start, _ := in["start_id"].(string)
			if start == "" {
				return "", "", map[string]any{
					"messages":      []any{map[string]any{"message_id": "m1"}},
					"has_continue":  true,
					"next_start_id": "c1",
				}
			}
			return "", "", map[string]any{
				"messages":     []any{map[string]any{"message_id": "m2"}},
				"has_continue": false,
			}
		},
	}))
	defer srv.Close()
	r := newTestClient(t, srv)
	ev, err := r.GetChatHistory(context.Background(), "sess", true, 100, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(ev.Slice("messages")) != 2 {
		t.Fatalf("expected 2 messages, got %v", ev)
	}
}

func TestDeleteEditSession(t *testing.T) {
	srv := httptest.NewServer(handler(map[string]func(map[string]any) (string, string, map[string]any){
		"deleteSession": func(in map[string]any) (string, string, map[string]any) {
			if in["session_id"] != "s1" {
				t.Errorf("delete session_id: %v", in["session_id"])
			}
			return "", "", map[string]any{"ok": true}
		},
		"editSession": func(in map[string]any) (string, string, map[string]any) {
			params, _ := in["update_parameters"].([]any)
			if len(params) != 2 {
				t.Errorf("update_parameters should have 2 entries, got %v", params)
			}
			return "", "", map[string]any{"ok": true}
		},
	}))
	defer srv.Close()
	r := newTestClient(t, srv)
	if _, err := r.DeleteSession(context.Background(), "s1"); err != nil {
		t.Fatal(err)
	}
	pinned := true
	if _, err := r.EditSession(context.Background(), "s1", "New Title", &pinned); err != nil {
		t.Fatal(err)
	}
}

func TestSendFeedbackValidation(t *testing.T) {
	r := New(WithAuth("x"))
	if _, err := r.SendFeedback(context.Background(), "s", "m", "maybe"); err == nil {
		t.Fatal("expected validation error for bad feedback")
	}
	if _, err := r.StopSessionStream(context.Background(), "s", "m", "maybe"); err == nil {
		t.Fatal("expected validation error for bad stop feedback")
	}
	srv := httptest.NewServer(handler(map[string]func(map[string]any) (string, string, map[string]any){
		"sendFeedback": func(in map[string]any) (string, string, map[string]any) {
			if in["feedback"] != "Like" {
				t.Errorf("feedback not normalised: %v", in["feedback"])
			}
			return "", "", map[string]any{"ok": true}
		},
		"stopSessionStream": func(in map[string]any) (string, string, map[string]any) {
			if in["feedback"] != "Dislike" {
				t.Errorf("stop feedback not normalised: %v", in["feedback"])
			}
			return "", "", map[string]any{"ok": true}
		},
	}))
	defer srv.Close()
	r = newTestClient(t, srv)
	if _, err := r.SendFeedback(context.Background(), "s", "m", "good"); err != nil {
		t.Fatal(err)
	}
	if _, err := r.StopSessionStream(context.Background(), "s", "m", "bad"); err != nil {
		t.Fatal(err)
	}
}

func TestSpeechToTextBase64(t *testing.T) {
	srv := httptest.NewServer(handler(map[string]func(map[string]any) (string, string, map[string]any){
		"speechToText": func(in map[string]any) (string, string, map[string]any) {
			vd, _ := in["voice_data"].(string)
			if vd == "" {
				t.Error("voice_data not base64-encoded")
			}
			return "", "", map[string]any{"text": "hello world"}
		},
	}))
	defer srv.Close()
	r := newTestClient(t, srv)
	ev, err := r.SpeechToText(context.Background(), []byte("fake-audio-bytes"))
	if err != nil {
		t.Fatal(err)
	}
	if ev.AnyStr("text") != "hello world" {
		t.Fatalf("text not returned: %v", ev)
	}
	if _, err := r.SpeechToText(context.Background(), nil); err == nil {
		t.Fatal("expected error for empty voice data")
	}
}

func TestDownloadImagePolling(t *testing.T) {
	var polls int32
	srv := httptest.NewServer(handler(map[string]func(map[string]any) (string, string, map[string]any){
		"downloadImage": func(in map[string]any) (string, string, map[string]any) {
			n := atomic.AddInt32(&polls, 1)
			if n < 3 {
				return "", "", map[string]any{"status": "Pending"}
			}
			// 1x1 transparent PNG base64.
			return "", "", map[string]any{"status": "Success", "image": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII="}
		},
	}))
	defer srv.Close()
	r := newTestClient(t, srv)
	dir := t.TempDir()
	out := dir + "/img.png"
	ev, err := r.DownloadImage(context.Background(), "task-1", out, 30*time.Second, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("DownloadImage: %v", err)
	}
	if ev.AnyStr("status") != "Success" {
		t.Fatalf("expected Success, got %v", ev)
	}
	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("image not written: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("written image is empty")
	}
}

func TestErrorMapping(t *testing.T) {
	srv := httptest.NewServer(handler(map[string]func(map[string]any) (string, string, map[string]any){
		"getServiceAccess": func(map[string]any) (string, string, map[string]any) {
			return "FAILED", "TOO_REQUESTS", nil
		},
	}))
	defer srv.Close()
	r := newTestClient(t, srv)
	_, err := r.GetServiceAccess(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	// The error carries the server status_det (TOO_REQUESTS) and maps to
	// the typed sentinel ErrTooManyRequests, which errors.Is recognises.
	if !strings.Contains(err.Error(), "TOO_REQUESTS") && !strings.Contains(err.Error(), "TooRequests") {
		t.Fatalf("expected TooRequests in error, got %v", err)
	}
}

// ----- SSE sendMessage -----

func sseEvent(name string, data map[string]any) string {
	jb, _ := json.Marshal(data)
	return fmt.Sprintf("event: %s\ndata: %s\n\n", name, string(jb))
}

func TestSendMessageSSEAssembly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var payload map[string]any
		_ = json.NewDecoder(req.Body).Decode(&payload)
		data, _ := payload["data"].(map[string]any)
		if data["method"] != "sendMessage" {
			http.Error(w, "bad method", http.StatusBadRequest)
			return
		}
		input, _ := data["input"].(map[string]any)
		if input["message"] != "hi" {
			http.Error(w, "bad message", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		// Stream: create session, create message, two text deltas, update title.
		fmt.Fprint(w, sseEvent("createSession", map[string]any{
			"session": map[string]any{"session_id": "S1", "title": "New chat"},
		}))
		flusher.Flush()
		fmt.Fprint(w, sseEvent("createMessage", map[string]any{
			"message": map[string]any{"message_id": "M1", "session_id": "S1"},
		}))
		flusher.Flush()
		fmt.Fprint(w, sseEvent("updateMessage", map[string]any{
			"parameter_name": "text_content", "append_value": "Hello ",
		}))
		flusher.Flush()
		fmt.Fprint(w, sseEvent("updateMessage", map[string]any{
			"parameter_name": "text_content", "append_value": "world!",
		}))
		flusher.Flush()
		fmt.Fprint(w, sseEvent("updateSession", map[string]any{
			"parameter_name": "title", "set_value": "Renamed",
		}))
		flusher.Flush()
	}))
	defer srv.Close()
	r := newTestClient(t, srv)

	var seen []string
	res, err := r.SendMessageStream(context.Background(), SendMessageOptions{Message: "hi"}, func(event string, _ map[string]any) {
		seen = append(seen, event)
	})
	if err != nil {
		t.Fatalf("SendMessageStream: %v", err)
	}
	if res.SessionID != "S1" || res.MessageID != "M1" {
		t.Fatalf("session/message not assembled: %+v", res)
	}
	if res.TextContent != "Hello world!" {
		t.Fatalf("text not concatenated: %q", res.TextContent)
	}
	if res.Title != "Renamed" {
		t.Fatalf("title not updated: %q", res.Title)
	}
	if len(seen) < 5 {
		t.Fatalf("expected ≥5 streamed events, got %d", len(seen))
	}
}

func TestSendMessageSSEError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		fmt.Fprint(w, sseEvent("error", map[string]any{
			"data": map[string]any{"status_det": "SERVER_ERROR", "status": "ERROR"},
		}))
		flusher.Flush()
	}))
	defer srv.Close()
	r := newTestClient(t, srv)
	_, err := r.SendMessage(context.Background(), SendMessageOptions{Message: "hi"})
	if err == nil {
		t.Fatal("expected error for SERVER_ERROR event")
	}
}

func TestSendMessageValidation(t *testing.T) {
	r := New(WithAuth("x"))
	if _, err := r.SendMessage(context.Background(), SendMessageOptions{Message: "   "}); err == nil {
		t.Fatal("expected validation error for empty message")
	}
	r2 := New() // no auth
	if _, err := r2.SendMessage(context.Background(), SendMessageOptions{Message: "hi"}); err == nil {
		t.Fatal("expected error when no auth set")
	}
}

func TestRequestRequiresAuth(t *testing.T) {
	r := New() // no auth
	if _, err := r.GetServiceAccess(context.Background()); err == nil {
		t.Fatal("expected error when no auth set")
	}
}
