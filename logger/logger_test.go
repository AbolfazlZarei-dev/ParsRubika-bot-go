package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newBuf(opts ...Option) (*Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	opts = append([]Option{WithOutput(&buf, ColorNever), WithTimeFormat("")}, opts...)
	return New(opts...), &buf
}

func TestLevelFiltering(t *testing.T) {
	l, buf := newBuf(WithLevel(LevelWarn))
	l.Info("hidden")
	l.Warn("shown")
	if strings.Contains(buf.String(), "hidden") || !strings.Contains(buf.String(), "shown") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestFieldsAndNamed(t *testing.T) {
	l, buf := newBuf()
	l.Named("bot").Named("poll").With(F("chat", "b0x")).Info("hello world", F("n", 3), Err(errors.New("boom")))
	out := buf.String()
	for _, want := range []string{"bot.poll", "hello world", "chat=b0x", "n=3", "error=boom"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in %q", want, out)
		}
	}
}

func TestQuotingAndNilError(t *testing.T) {
	l, buf := newBuf()
	var nilErr error
	l.Info("m", F("k", "a b"), Err(nilErr))
	if !strings.Contains(buf.String(), `k="a b"`) || strings.Contains(buf.String(), "error=") {
		t.Fatalf("bad output: %q", buf.String())
	}
}

func TestColorAlways(t *testing.T) {
	var buf bytes.Buffer
	l := New(WithOutput(&buf, ColorAlways))
	l.Error("bad")
	if !strings.Contains(buf.String(), "\x1b[") {
		t.Fatalf("expected ANSI codes: %q", buf.String())
	}
}

func TestJSON(t *testing.T) {
	var buf bytes.Buffer
	l := New(WithOutput(&buf, ColorNever), WithJSON())
	_ = l
	l = New(AddOutput(&buf, FormatJSON, ColorNever, LevelTrace), WithOutput(&bytes.Buffer{}, ColorNever))
	l = New(WithOutput(&bytes.Buffer{}, ColorNever), AddOutput(&buf, FormatJSON, ColorNever, LevelTrace))
	l.Info("سلام", F("id", 7), Err(errors.New("x")))
	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("invalid json %q: %v", buf.String(), err)
	}
	if m["msg"] != "سلام" || m["level"] != "info" || m["error"] != "x" || m["id"].(float64) != 7 {
		t.Fatalf("bad record: %v", m)
	}
}

func TestPrintfCompatAndStd(t *testing.T) {
	l, buf := newBuf()
	l.Printf("a %d", 1)
	l.StdLogger(LevelWarn).Println("via std")
	out := buf.String()
	if !strings.Contains(out, "a 1") || !strings.Contains(out, "via std") || !strings.Contains(out, "WARN") {
		t.Fatalf("bad output: %q", out)
	}
}

func TestRotatingFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "x.log")
	r, err := NewRotatingFile(p, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	chunk := bytes.Repeat([]byte("x"), 600<<10)
	for i := 0; i < 5; i++ {
		if _, err := r.Write(chunk); err != nil {
			t.Fatal(err)
		}
	}
	_ = r.Close()
	for _, f := range []string{p, p + ".1", p + ".2"} {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("expected %s: %v", f, err)
		}
	}
	if _, err := os.Stat(p + ".3"); err == nil {
		t.Error("too many backups kept")
	}
}

func TestSetLevelAffectsChildren(t *testing.T) {
	l, buf := newBuf(WithLevel(LevelError))
	c := l.Named("c")
	c.Info("no")
	l.SetLevel(LevelInfo)
	c.Info("yes")
	if strings.Contains(buf.String(), "no") && !strings.Contains(buf.String(), "yes") {
		t.Fatal("level change not shared")
	}
}

func TestNop(t *testing.T) { Nop().Error("nothing") }
