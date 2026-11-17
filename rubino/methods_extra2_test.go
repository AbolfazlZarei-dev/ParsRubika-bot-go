package rubino

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// withTestServer points the Rubino client at an httptest server that
// speaks the Rubino api_version:"0" envelope (method + data at the top
// level) and routes on the `method` field.
func withTestServer(t *testing.T, routes map[string]func(input map[string]any) map[string]any) (*Rubino, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
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
	c := New(WithAuth("test-auth"), WithProfile("my-profile-id"))
	c.apiURL = srv.URL
	c.SetHTTPClient(srv.Client())
	return c, srv.Close
}

func TestProfileCreatedAtDecodes(t *testing.T) {
	// 0x60000000 = 1610280960 = 2021-01-10T12:16:00Z
	ts := uint32(0x60000000)
	want := time.Unix(int64(ts), 0).UTC()
	id := "60000000" + strings.Repeat("a", 16) // 24 hex chars
	got, err := ProfileCreatedAt(id)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Equal(want) {
		t.Fatalf("ProfileCreatedAt = %v, want %v", got, want)
	}
	if _, err := ProfileCreatedAt("abc"); err == nil {
		t.Fatal("expected error for short profile_id")
	}
	if _, err := ProfileCreatedAt("zzzzzzzz" + strings.Repeat("a", 16)); err == nil {
		t.Fatal("expected error for non-hex profile_id")
	}
}

func TestProfileAdminGroup(t *testing.T) {
	var seen []string
	c, cleanup := withTestServer(t, map[string]func(map[string]any) map[string]any{
		"setProfileAdmin": func(in map[string]any) map[string]any {
			seen = append(seen, in["action"].(string))
			if in["profile_id"] != "my-profile-id" {
				t.Errorf("profile_id not defaulted: %v", in["profile_id"])
			}
			return map[string]any{"ok": true}
		},
		"profileAdminAction": func(in map[string]any) map[string]any {
			seen = append(seen, in["action"].(string))
			return map[string]any{"ok": true}
		},
		"getProfileAdmins": func(in map[string]any) map[string]any {
			return map[string]any{"admins": []any{}}
		},
		"getSettingAdminPermissions": func(map[string]any) map[string]any {
			return map[string]any{"permissions": []any{}}
		},
	})
	defer cleanup()
	ctx := context.Background()
	if _, err := c.SetProfileAdmin(ctx, "", "tp", "", []string{"EditProfile"}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.UnsetProfileAdmin(ctx, "", "tp"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.AcceptProfileAdmin(ctx, "", "tp"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.DenyProfileAdmin(ctx, "", "tp"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.RemoveProfileAdmin(ctx, "", "tp"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetProfileAdmins(ctx, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetSettingAdminPermissions(ctx); err != nil {
		t.Fatal(err)
	}
	want := []string{"SetAdmin", "UnsetAdmin", "Accept", "Deny", "Remove"}
	if strings.Join(seen, ",") != strings.Join(want, ",") {
		t.Fatalf("admin actions: %v want %v", seen, want)
	}
}

func TestCloseFriendsAndSettings(t *testing.T) {
	c, cleanup := withTestServer(t, map[string]func(map[string]any) map[string]any{
		"closeFriendAction": func(in map[string]any) map[string]any {
			if in["action"] != "Set" && in["action"] != "Unset" {
				t.Errorf("action: %v", in["action"])
			}
			return map[string]any{"ok": true}
		},
		"getMyCloseFriends": func(map[string]any) map[string]any {
			return map[string]any{"friends": []any{}}
		},
		"getProfileSetting": func(in map[string]any) map[string]any {
			if in["profile_id"] != "my-profile-id" {
				t.Errorf("profile_id: %v", in["profile_id"])
			}
			return map[string]any{"setting": "x"}
		},
		"setProfileSetting": func(in map[string]any) map[string]any {
			if in["profile_id"] != "my-profile-id" {
				t.Errorf("profile_id: %v", in["profile_id"])
			}
			return map[string]any{"ok": true}
		},
	})
	defer cleanup()
	ctx := context.Background()
	if _, err := c.SetCloseFriends(ctx, "", []string{"a", "b"}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.UnsetCloseFriends(ctx, "", []string{"a"}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetMyCloseFriends(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetProfileSetting(ctx, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := c.SetProfileSetting(ctx, "", map[string]any{"private": true}); err != nil {
		t.Fatal(err)
	}
}

func TestSearchGroup(t *testing.T) {
	c, cleanup := withTestServer(t, map[string]func(map[string]any) map[string]any{
		"searchPost": func(in map[string]any) map[string]any {
			if in["usage"] != "Search" {
				t.Errorf("usage default: %v", in["usage"])
			}
			return map[string]any{"posts": []any{}}
		},
		"searchHashTag": func(in map[string]any) map[string]any {
			if in["search_text"] != "funny" {
				t.Errorf("search_text: %v", in["search_text"])
			}
			return map[string]any{"posts": []any{}}
		},
		"searchLocation": func(in map[string]any) map[string]any {
			if in["lat"] != 35.7 || in["lng"] != 51.4 {
				t.Errorf("lat/lng: %v/%v", in["lat"], in["lng"])
			}
			return map[string]any{"posts": []any{}}
		},
	})
	defer cleanup()
	ctx := context.Background()
	if _, err := c.SearchPost(ctx, "cats", "", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := c.SearchHashTag(ctx, "funny", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := c.SearchLocation(ctx, "Tehran", 35.7, 51.4, ""); err != nil {
		t.Fatal(err)
	}
}

func TestPostsAndMetricsGroup(t *testing.T) {
	c, cleanup := withTestServer(t, map[string]func(map[string]any) map[string]any{
		"getPostsByIds": func(in map[string]any) map[string]any {
			if len(in["post_ids"].([]any)) != 2 {
				t.Errorf("post_ids: %v", in["post_ids"])
			}
			return map[string]any{"posts": []any{}}
		},
		"updatePost": func(in map[string]any) map[string]any {
			if in["caption"] != "new" {
				t.Errorf("caption not forwarded: %v", in["caption"])
			}
			return map[string]any{"ok": true}
		},
		"updatePostShareCount":   func(map[string]any) map[string]any { return map[string]any{"ok": true} },
		"postCollaboratorAction": func(map[string]any) map[string]any { return map[string]any{"ok": true} },
		"setPostFeedback":        func(map[string]any) map[string]any { return map[string]any{"ok": true} },
		"getFeedData":            func(map[string]any) map[string]any { return map[string]any{"feed": []any{}} },
		"viewFeedData":           func(map[string]any) map[string]any { return map[string]any{"ok": true} },
		"profileAction":          func(map[string]any) map[string]any { return map[string]any{"ok": true} },
		"setProfileFeedback":     func(map[string]any) map[string]any { return map[string]any{"ok": true} },
		"getBoostUrl":            func(map[string]any) map[string]any { return map[string]any{"url": "x"} },
		"getMetricsConfig":       func(map[string]any) map[string]any { return map[string]any{"config": "x"} },
		"trackClientMetrics":     func(map[string]any) map[string]any { return map[string]any{"ok": true} },
		"addImpressionIds":       func(map[string]any) map[string]any { return map[string]any{"ok": true} },
		"addUsageTime": func(in map[string]any) map[string]any {
			// JSON numbers decode as float64; compare accordingly.
			if v, _ := in["usage_time"].(float64); v != 30 {
				t.Errorf("usage_time: %v", in["usage_time"])
			}
			return map[string]any{"ok": true}
		},
		"settingMessageOwnerAction": func(in map[string]any) map[string]any {
			return map[string]any{"ok": in["action"]}
		},
	})
	defer cleanup()
	ctx := context.Background()
	if _, err := c.GetPostsByIDs(ctx, "", []string{"p1", "p2"}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.UpdatePost(ctx, "", "p1", map[string]any{"caption": "new"}, []string{"caption"}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.UpdatePostShareCount(ctx, "", "p1"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.PostCollaboratorAction(ctx, "", "p1", "tp", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := c.SetPostFeedback(ctx, "", "p1", "NotInterested"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetFeedData(ctx, "", "", 20); err != nil {
		t.Fatal(err)
	}
	if _, err := c.ViewFeedData(ctx, "", []string{"r1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.RemoveFollower(ctx, "", "tp"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.NotInterestedProfile(ctx, "", "tp"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetBoostURL(ctx, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetMetricsConfig(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := c.TrackClientMetrics(ctx, map[string]any{"k": "v"}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.AddImpressionIDs(ctx, "", []string{"i1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.AddUsageTime(ctx, "", 30); err != nil {
		t.Fatal(err)
	}
	if _, err := c.AcceptMessageOwner(ctx, "", "p1"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.DenyMessageOwner(ctx, "", "p1"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.RemoveMessageOwner(ctx, "", "p1"); err != nil {
		t.Fatal(err)
	}
}
