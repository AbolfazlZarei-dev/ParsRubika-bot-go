package account

import (
	"strings"
	"testing"
)

func TestExtractFileNameSanitizesPathTraversal(t *testing.T) {
	cases := map[string]string{
		`attachment; filename="../../etc/passwd"`:       "passwd",
		`attachment; filename="..\..\windows\evil.exe"`: "evil.exe",
		`attachment; filename="ok.txt"`:                 "ok.txt",
		`attachment; filename="/abs/path/file.bin"`:     "file.bin",
		"attachment; filename=\"bad\x00name\"":          "badname", // real null byte stripped
		`attachment; filename="."`:                      "",
		`attachment; filename=".."`:                     "",
		`no filename here`:                              "",
	}
	for cd, want := range cases {
		got := extractFileName(cd)
		if got != want {
			t.Errorf("extractFileName(%q) = %q, want %q", cd, got, want)
		}
	}
}

func TestSafeExtFromMIME(t *testing.T) {
	cases := map[string]string{
		"image/jpeg":               "jpg",
		"video/mp4":                "mp4",
		"audio/mpeg":               "mp3",
		"application/octet-stream": "", // hyphen is not alphanumeric -> rejected
		"image/png; charset=x":     "png",
		"":                         "",
		"weird/ty\npe":             "", // control char -> rejected
	}
	for mime, want := range cases {
		got := safeExtFromMIME(mime)
		if got != want {
			t.Errorf("safeExtFromMIME(%q) = %q, want %q", mime, got, want)
		}
	}
	// The path-traversal guard: a malicious MIME must never yield an
	// extension containing a path separator or dot.
	for _, evil := range []string{"image/../../x", "image/..\\x", "image/a/b"} {
		ext := safeExtFromMIME(evil)
		if ext != "" && strings.ContainsAny(ext, "/\\.") {
			t.Errorf("safeExtFromMIME(%q) = %q leaked path characters", evil, ext)
		}
	}
}
