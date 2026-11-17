package text

import (
	"strings"
	"testing"
)

// BenchmarkFormatSimple measures the cost of formatting a string with
// no markdown constructs (worst case for the regex scan).
func BenchmarkFormatSimple(b *testing.B) {
	f := New()
	text := strings.Repeat("hello world ", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Format(text)
	}
}

// BenchmarkFormatMixed measures the cost of formatting a string with
// several markdown constructs.
func BenchmarkFormatMixed(b *testing.B) {
	f := New()
	text := "**bold** __italic__ `code` ~~strike~~ --under-- ||spoiler|| [link](https://example.com)"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Format(text)
	}
}

// BenchmarkToMarkdownHTML measures the cost of the HTML → Markdown
// conversion for a typical HTML payload.
func BenchmarkToMarkdownHTML(b *testing.B) {
	f := New()
	html := "<b>bold</b> <i>italic</i> <code>mono</code> <a href=\"https://example.com\">link</a>"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.ToMarkdown(html)
	}
}
