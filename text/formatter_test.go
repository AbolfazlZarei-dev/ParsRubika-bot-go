package text

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFormatBold(t *testing.T) {
	f := New()
	out := f.Format("hello **world**!")
	if out.Text != "hello world!" {
		t.Fatalf("Text = %q, want %q", out.Text, "hello world!")
	}
	if out.Metadata == nil || len(out.Metadata.Parts) == 0 {
		t.Fatalf("expected at least one part, got %+v", out)
	}
	first := out.Metadata.Parts[0]
	if first.Type != "Bold" {
		t.Errorf("first part type = %q, want %q", first.Type, "Bold")
	}
	if first.FromIndex != 6 {
		t.Errorf("from_index = %d, want 6", first.FromIndex)
	}
	if first.Length != 5 {
		t.Errorf("length = %d, want 5", first.Length)
	}
}

func TestFormatMultiple(t *testing.T) {
	f := New()
	out := f.Format("**bold** and __italic__")
	if !strings.Contains(out.Text, "bold") || !strings.Contains(out.Text, "italic") {
		t.Fatalf("text lost during parsing: %q", out.Text)
	}
	if out.Metadata == nil {
		t.Fatal("expected metadata")
	}
	if len(out.Metadata.Parts) < 2 {
		t.Fatalf("expected at least 2 parts, got %d", len(out.Metadata.Parts))
	}
}

func TestFormatMono(t *testing.T) {
	f := New()
	out := f.Format("see `code` here")
	if out.Text != "see code here" {
		t.Fatalf("Text = %q, want %q", out.Text, "see code here")
	}
	if out.Metadata == nil {
		t.Fatal("expected metadata")
	}
	found := false
	for _, p := range out.Metadata.Parts {
		if p.Type == "Mono" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Mono part not found: %+v", out.Metadata)
	}
}

func TestFormatPre(t *testing.T) {
	f := New()
	out := f.Format("```go\nfmt.Println(\"hi\")\n```")
	if !strings.Contains(out.Text, "fmt.Println") {
		t.Fatalf("text missing code: %q", out.Text)
	}
	if out.Metadata == nil {
		t.Fatal("expected metadata")
	}
	var pre *Part
	for i := range out.Metadata.Parts {
		if out.Metadata.Parts[i].Type == "Pre" {
			pre = &out.Metadata.Parts[i]
			break
		}
	}
	if pre == nil {
		t.Fatal("Pre part not found")
	}
	// The formatter only treats the first line as a
	// language tag when it does NOT start with an alphanumeric
	// character (so "go\n..." is treated as code, not as lang="go").
	// Sending language="go" makes the Rubika server reject the whole
	// message with INVALID_INPUT, so the formatter behaves accordingly
	// here: language is empty when the first line is alphanumeric.
	if pre.Language != "" {
		t.Errorf("Pre language = %q, want %q (alphanumeric first line must not be treated as a language)", pre.Language, "")
	}
}

func TestFormatLink(t *testing.T) {
	f := New()
	out := f.Format("[click](https://example.com)")
	if out.Text != "click" {
		t.Fatalf("Text = %q, want %q", out.Text, "click")
	}
	if out.Metadata == nil {
		t.Fatal("expected metadata")
	}
	var link *Part
	for i := range out.Metadata.Parts {
		if out.Metadata.Parts[i].Type == "Link" {
			link = &out.Metadata.Parts[i]
			break
		}
	}
	if link == nil {
		t.Fatal("Link part not found")
	}
	if link.LinkURL != "https://example.com" {
		t.Errorf("link URL = %q, want %q", link.LinkURL, "https://example.com")
	}
}

func TestFormatSpoiler(t *testing.T) {
	f := New()
	out := f.Format("this is a ||secret||")
	if !strings.Contains(out.Text, "secret") {
		t.Fatalf("text missing content: %q", out.Text)
	}
	found := false
	for _, p := range out.Metadata.Parts {
		if p.Type == "Spoiler" {
			found = true
		}
	}
	if !found {
		t.Fatal("Spoiler part not found")
	}
}

func TestFormatStrike(t *testing.T) {
	f := New()
	out := f.Format("this is ~~struck~~")
	if !strings.Contains(out.Text, "struck") {
		t.Fatalf("text missing content: %q", out.Text)
	}
	found := false
	for _, p := range out.Metadata.Parts {
		if p.Type == "Strike" {
			found = true
		}
	}
	if !found {
		t.Fatal("Strike part not found")
	}
}

func TestFormatUnderline(t *testing.T) {
	f := New()
	// Use a string with `--` (the Underline prefix), which does not
	// collide with the Italic `__` prefix.
	out := f.Format("this is --underlined--")
	if !strings.Contains(out.Text, "underlined") {
		t.Fatalf("text missing content: %q", out.Text)
	}
	found := false
	for _, p := range out.Metadata.Parts {
		if p.Type == "Underline" {
			found = true
		}
	}
	if !found {
		t.Fatal("Underline part not found")
	}
}

func TestFormatQuote(t *testing.T) {
	f := New()
	out := f.Format("> a quoted line")
	if out.Text != "a quoted line" {
		t.Fatalf("Text = %q, want %q", out.Text, "a quoted line")
	}
	found := false
	for _, p := range out.Metadata.Parts {
		if p.Type == "Quote" {
			found = true
		}
	}
	if !found {
		t.Fatal("Quote part not found")
	}
}

// TestFormatQuoteMidText verifies that a quoted block that appears
// in the MIDDLE of a message (not at the very start) is detected
// and produces a Quote metadata part. Without the `(?m)` MULTILINE
// flag on the regex, Go's `^` only matches the start of the string,
// so mid-text quotes were silently ignored.
func TestFormatQuoteMidText(t *testing.T) {
	f := New()
	out := f.Format("📨 شما گفتید:\n\n> تست\n\n`/help` برای دستورات")
	// The Quote "> تست" should be stripped to "تست".
	if !strings.Contains(out.Text, "تست") {
		t.Errorf("output should contain 'تست', got %q", out.Text)
	}
	if strings.Contains(out.Text, "> تست") {
		t.Errorf("Quote marker '> ' should be stripped, got %q", out.Text)
	}
	if out.Metadata == nil {
		t.Fatal("expected metadata with a Quote part, got nil")
	}
	quoteFound := false
	for _, p := range out.Metadata.Parts {
		if p.Type == "Quote" {
			quoteFound = true
			if p.Length != 3 {
				t.Errorf("Quote length = %d, want 3 (for 'تست')", p.Length)
			}
		}
	}
	if !quoteFound {
		t.Errorf("expected a Quote part in metadata, got: %+v", out.Metadata)
	}
}

func TestFormatNoMetadata(t *testing.T) {
	f := New()
	out := f.Format("plain text with no formatting")
	if out.Text != "plain text with no formatting" {
		t.Fatalf("Text = %q, want %q", out.Text, "plain text with no formatting")
	}
	if out.Metadata != nil {
		t.Fatalf("expected no metadata, got %+v", out.Metadata)
	}
}

func TestResultJSON(t *testing.T) {
	f := New()
	out := f.Format("**bold**")
	js := out.JSON()
	if !strings.Contains(js, "Bold") {
		t.Fatalf("JSON missing Bold: %s", js)
	}
	var v any
	if err := json.Unmarshal([]byte(js), &v); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestToMarkdownNoHTML(t *testing.T) {
	f := New()
	out := f.ToMarkdown("plain text")
	if out != "plain text" {
		t.Fatalf("ToMarkdown(plain) = %q, want %q", out, "plain text")
	}
}

func TestToMarkdownHTML(t *testing.T) {
	f := New()
	out := f.ToMarkdown("<b>bold</b>")
	if !strings.Contains(out, "**bold**") {
		t.Fatalf("ToMarkdown(<b>bold</b>) = %q, want **bold**", out)
	}
	out = f.ToMarkdown("<i>italic</i>")
	if !strings.Contains(out, "__italic__") {
		t.Fatalf("ToMarkdown(<i>italic</i>) = %q, want __italic__", out)
	}
	out = f.ToMarkdown("<code>mono</code>")
	if !strings.Contains(out, "`mono`") {
		t.Fatalf("ToMarkdown(<code>mono</code>) = %q, want `mono`", out)
	}
}

// TestFormatPersianText verifies that the formatter handles Persian
// (and other multi-byte UTF-8) text correctly, preserving every
// character and matching the input rune-for-rune.
func TestFormatPersianText(t *testing.T) {
	f := New()

	// Single bold word in Persian.
	out := f.Format("🎉 **خوش آمدید!**")
	want := "🎉 خوش آمدید!"
	if out.Text != want {
		t.Errorf("single-line bold: Text = %q, want %q", out.Text, want)
	}
	if out.Metadata == nil || len(out.Metadata.Parts) != 1 {
		t.Fatalf("single-line bold: expected 1 part, got %+v", out.Metadata)
	}
	if out.Metadata.Parts[0].Type != "Bold" {
		t.Errorf("single-line bold: type = %q, want Bold", out.Metadata.Parts[0].Type)
	}
	if out.Metadata.Parts[0].FromIndex != 3 {
		t.Errorf("single-line bold: from_index = %d, want 3", out.Metadata.Parts[0].FromIndex)
	}
	if out.Metadata.Parts[0].Length != 10 {
		t.Errorf("single-line bold: length = %d, want 10", out.Metadata.Parts[0].Length)
	}

	// Multi-line text with multiple markdown blocks.
	out = f.Format("📋 **دستورات موجود**\n\n• `/start` — پیام خوش\u200cآمد\n• `/help`  — این لیست")
	want = "📋 دستورات موجود\n\n• /start — پیام خوش\u200cآمد\n• /help  — این لیست"
	if out.Text != want {
		t.Errorf("multi-line: Text = %q, want %q", out.Text, want)
	}
	if out.Metadata == nil || len(out.Metadata.Parts) != 3 {
		t.Fatalf("multi-line: expected 3 parts, got %+v", out.Metadata)
	}
	// Check for any byte-level corruption: the output should be
	// valid UTF-8 and contain no replacement chars or stray bytes.
	for _, r := range out.Text {
		if r == 0xFFFD {
			t.Errorf("multi-line: output contains replacement char (0xFFFD)")
		}
	}

	// Bullet list with bold + mono on each line.
	out = f.Format("• **چت:** `g0Ic2Fu0a0d535d068c9c54225a2aee5`\n• **کاربر:** `u0USER`")
	want = "• چت: g0Ic2Fu0a0d535d068c9c54225a2aee5\n• کاربر: u0USER"
	if out.Text != want {
		t.Errorf("bullet list: Text = %q, want %q", out.Text, want)
	}

	// Quote + Bold in Persian. The Quote regex matches `> hello\n`
	// (8 bytes including the trailing newline), and the content
	// extraction strips the `> ` prefix AND the trailing newline
	// (because strings.TrimSpace is applied at the end of Format).
	// The result is `hello` directly followed by the second `\n`
	// and the rest of the original text — i.e. a single `\n`.
	out = f.Format("> hello\n\n📨 **بازتاب:** hello")
	want = "hello\n📨 بازتاب: hello"
	if out.Text != want {
		t.Errorf("quote+bold: Text = %q, want %q", out.Text, want)
	}

	// ZWNJ handling — the formatter MUST NOT split
	// the ZWNJ (U+200C) inside "پیام\u200cها".
	out = f.Format("📊 **آمار پیام\u200cها:** `42`")
	want = "📊 آمار پیام\u200cها: 42"
	if out.Text != want {
		t.Errorf("zwnj: Text = %q, want %q", out.Text, want)
	}
}
