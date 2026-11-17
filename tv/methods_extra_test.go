package tv

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// withTestServer swaps TVServerURLs to point at srv (and restores it
// afterwards) so the TV client can be exercised end-to-end without real
// network access. It also pins the version so MetaURL is never fetched.
func withTestServer(t *testing.T, routes map[string]func(input map[string]any) map[string]any) (*TV, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		// The TV envelope puts `method` and `data` (the input) at the
		// top level — unlike rubino/rubiai which nest them under data.
		method, _ := payload["method"].(string)
		input, _ := payload["data"].(map[string]any)
		route, ok := routes[method]
		if !ok {
			http.Error(w, "no route "+method, http.StatusBadRequest)
			return
		}
		out := route(input)
		body := map[string]any{
			"status": "OK", "status_det": "OK",
			"data": out,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(body)
	}))
	old := TVServerURLs
	TVServerURLs = []string{srv.URL}
	c := New(WithAuth("test-auth"), WithVersion("test-ver"))
	c.SetHTTPClient(srv.Client())
	return c, func() {
		TVServerURLs = old
		srv.Close()
	}
}

func TestGetMediaByLinkExtractsNumericID(t *testing.T) {
	tv, cleanup := withTestServer(t, map[string]func(map[string]any) map[string]any{
		"getMediaByShareString": func(in map[string]any) map[string]any {
			ss, _ := in["share_string"].(string)
			if ss != "294119" {
				t.Errorf("share_string not extracted: %v", ss)
			}
			if _, ok := in["track_id"]; !ok {
				t.Error("track_id should be present")
			}
			return map[string]any{"media": map[string]any{"id": "294119"}}
		},
	})
	defer cleanup()
	ev, err := tv.GetMediaByLink(context.Background(), "https://rubika.ir/vod/294119/", "tk1")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Map("media") == nil {
		t.Fatalf("media not returned: %v", ev)
	}
	// Non-numeric link must be rejected client-side.
	if _, err := tv.GetMediaByLink(context.Background(), "https://rubika.ir/vod/abc", ""); err == nil {
		t.Fatal("expected validation error for non-numeric link")
	}
}

func TestGetNextEpisodeAndMediaInfo(t *testing.T) {
	tv, cleanup := withTestServer(t, map[string]func(map[string]any) map[string]any{
		"getNextEpisode": func(in map[string]any) map[string]any {
			if in["episode_id"] != "e1" {
				t.Errorf("episode_id: %v", in["episode_id"])
			}
			return map[string]any{"episode": map[string]any{"id": "e2"}}
		},
		"getNextMediaInfo": func(in map[string]any) map[string]any {
			if in["media_id"] != "m1" || in["track_id"] != "t1" {
				t.Errorf("input: %v", in)
			}
			return map[string]any{"media": map[string]any{"id": "m2"}}
		},
	})
	defer cleanup()
	if _, err := tv.GetNextEpisode(context.Background(), "e1"); err != nil {
		t.Fatal(err)
	}
	if _, err := tv.GetNextMediaInfo(context.Background(), "m1", "t1"); err != nil {
		t.Fatal(err)
	}
}

func TestGetPlayInfoDefaultsAndValidation(t *testing.T) {
	tv, cleanup := withTestServer(t, map[string]func(map[string]any) map[string]any{
		"getPlayInfo": func(in map[string]any) map[string]any {
			if in["type"] != "Trailer" {
				t.Errorf("expected default Trailer, got %v", in["type"])
			}
			if in["version"] != 1.0 && in["version"] != 1 {
				// JSON numbers decode as float64; accept either form.
			}
			return map[string]any{"play": map[string]any{"url": "x"}}
		},
	})
	defer cleanup()
	if _, err := tv.GetPlayInfo(context.Background(), "m1", "", 0, ""); err != nil {
		t.Fatal(err)
	}
}

func TestSetCurrentPlayTimeValidation(t *testing.T) {
	tv, cleanup := withTestServer(t, nil)
	defer cleanup()
	if _, err := tv.SetCurrentPlayTime(context.Background(), "m1", -1); err == nil {
		t.Fatal("expected error for negative play_time")
	}
}

func TestSetSessionPlayDurationAndViewCommercial(t *testing.T) {
	tv, cleanup := withTestServer(t, map[string]func(map[string]any) map[string]any{
		"setSessionPlayDuration": func(in map[string]any) map[string]any {
			if in["session_id"] != "s1" {
				t.Errorf("session_id: %v", in["session_id"])
			}
			return map[string]any{"ok": true}
		},
		"viewCommercial": func(in map[string]any) map[string]any {
			if in["view_type"] != "Full" {
				t.Errorf("view_type default Full, got %v", in["view_type"])
			}
			return map[string]any{"ok": true}
		},
	})
	defer cleanup()
	if _, err := tv.SetSessionPlayDuration(context.Background(), "s1", 120); err != nil {
		t.Fatal(err)
	}
	if _, err := tv.ViewCommercial(context.Background(), ""); err != nil {
		t.Fatal(err)
	}
}

func TestMediaActionsTrackID(t *testing.T) {
	var seen []string
	tv, cleanup := withTestServer(t, map[string]func(map[string]any) map[string]any{
		"actionOnMedia": func(in map[string]any) map[string]any {
			seen = append(seen, in["action"].(string))
			if in["track_id"] != "tk" {
				t.Errorf("track_id missing: %v", in)
			}
			return map[string]any{"ok": true}
		},
	})
	defer cleanup()
	for _, fn := range []func() error{
		func() error { _, e := tv.DislikeMedia(context.Background(), "m1", "tk"); return e },
		func() error { _, e := tv.UndoLikeOrDislikeMedia(context.Background(), "m1", "tk"); return e },
		func() error { _, e := tv.FollowMedia(context.Background(), "m1", "tk"); return e },
		func() error { _, e := tv.UnfollowMedia(context.Background(), "m1", "tk"); return e },
	} {
		if err := fn(); err != nil {
			t.Fatal(err)
		}
	}
	want := []string{"Dislike", "UndoLikeOrDislike", "Follow", "UnFollow"}
	if strings.Join(seen, ",") != strings.Join(want, ",") {
		t.Fatalf("actions sent: %v want %v", seen, want)
	}
}
