package account

import (
	"strings"
)

// Markdown text formatting helpers.

// Bold wraps text in **bold** markers.
func Bold(text string) string { return "**" + text + "**" }

// Italic wraps text in __italic__ markers.
func Italic(text string) string { return "__" + text + "__" }

// Underline wraps text in --underline-- markers.
func Underline(text string) string { return "--" + text + "--" }

// Strike wraps text in ~~strike~~ markers.
func Strike(text string) string { return "~~" + text + "~~" }

// Spoiler wraps text in ||spoiler|| markers.
func Spoiler(text string) string { return "||" + text + "||" }

// Code wraps text in `code` markers (inline monospace).
func Code(text string) string { return "`" + text + "`" }

// Pre wraps text in ```pre``` markers (block monospace).
func Pre(text string) string { return "```" + text + "```" }

// Hyperlink creates a [text](url) markdown link.
func Hyperlink(text, url string) string { return "[" + text + "](" + url + ")" }

// Mention creates a mention link for a GUID. The URL prefix is inferred
// from the GUID: u0→User, g0→Group, c0→Channel, b0→Bot.
func Mention(text, guid string) string {
	prefix := "u"
	switch {
	case strings.HasPrefix(guid, "g0"):
		prefix = "g"
	case strings.HasPrefix(guid, "c0"):
		prefix = "c"
	case strings.HasPrefix(guid, "b0"):
		prefix = "b"
	}
	return Hyperlink(text, prefix+guid)
}

// Quote creates a >quote block.
func Quote(text string) string { return "> " + text }
