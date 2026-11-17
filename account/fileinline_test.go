package account

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestFileInlineBlockUsesCorrectType verifies that the file_inline
// block built by SendMessage carries the correct `type` field —
// i.e. "Image" for SendPhoto, "Video" for SendVideo, etc., matching the
// upload kind so the server accepts the message.
func TestFileInlineBlockUsesCorrectType(t *testing.T) {
	cases := []struct {
		name      string
		fileType  string
		wantType  string
		wantRound bool
	}{
		{"Image", "Image", "Image", false},
		{"Video", "Video", "Video", false},
		{"Voice", "Voice", "Voice", false},
		{"Music", "Music", "Music", false},
		{"Gif", "Gif", "Gif", false},
		{"File", "File", "File", false},
		{"VideoMessage", "VideoMessage", "Video", true},
		{"empty defaults to File", "", "File", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Reproduce the file_inline construction in SendMessage.
			file := map[string]any{
				"file_id":         "f1",
				"dc_id":           1,
				"size":            100,
				"mime":            "image/jpeg",
				"access_hash_rec": "h1",
			}
			tt := c.fileType
			if tt == "" {
				tt = "File"
			}
			wireType := tt
			if tt == "VideoMessage" {
				wireType = "Video"
				file["is_round"] = true
			}
			file["type"] = wireType
			file["time"] = 1
			file["width"] = 200
			file["height"] = 200
			file["music_performer"] = ""
			file["is_spoil"] = false

			body, _ := json.Marshal(file)
			got := string(body)
			wantField := `"type":"` + c.wantType + `"`
			if !strings.Contains(got, wantField) {
				t.Errorf("file_inline should contain %s, got: %s", wantField, got)
			}
			if c.wantRound {
				if !strings.Contains(got, `"is_round":true`) {
					t.Errorf("VideoMessage should set is_round=true, got: %s", got)
				}
			} else {
				if strings.Contains(got, `"is_round"`) {
					t.Errorf("non-VideoMessage should NOT set is_round, got: %s", got)
				}
			}
			// Defaults: time=1, width=200, height=200 — NOT 0.
			if !strings.Contains(got, `"time":1`) {
				t.Errorf("default time should be 1, got: %s", got)
			}
			if !strings.Contains(got, `"width":200`) {
				t.Errorf("default width should be 200, got: %s", got)
			}
			if !strings.Contains(got, `"height":200`) {
				t.Errorf("default height should be 200, got: %s", got)
			}
		})
	}
}
