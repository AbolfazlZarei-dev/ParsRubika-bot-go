package enums

// ParseMode controls how text is formatted before being sent to the
// server.
type ParseMode string

const (
	// ParseModeNone sends text as-is (no metadata generation).
	ParseModeNone ParseMode = ""
	// ParseModeMarkdown runs the text through the Markdown formatter.
	ParseModeMarkdown ParseMode = "markdown"
	// ParseModeHTML first converts HTML to Markdown via ToMarkdown and
	// then runs the Markdown formatter.
	ParseModeHTML ParseMode = "html"
)

// String implements fmt.Stringer.
func (p ParseMode) String() string { return string(p) }

// Valid reports whether p is one of the supported parse modes.
func (p ParseMode) Valid() bool {
	switch p {
	case ParseModeNone, ParseModeMarkdown, ParseModeHTML:
		return true
	}
	return false
}
