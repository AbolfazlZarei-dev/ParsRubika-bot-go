package text

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"
)

// MarkupType maps a markdown prefix (e.g. "**", "__", "```") to its
// wire-format type name and the regex group index that holds the inner
// content (1-indexed, like regexp submatch groups).
type MarkupType struct {
	Prefix   string
	Type     string
	GroupIdx int
}

// MarkupTypeSequence lists the markup kinds in the order in which the
// regex alternation tries to match them.
var MarkupTypeSequence = []MarkupType{
	{">", "Quote", 0},
	{"```", "Pre", 1},
	{"**", "Bold", 2},
	{"`", "Mono", 3},
	{"__", "Italic", 4},
	{"--", "Underline", 5},
	{"~~", "Strike", 6},
	{"||", "Spoiler", 7},
	{"[", "Link", 8},
}

// MentionTypes maps a one-letter URL prefix (u/g/c/b) to its mention
// type used on the wire.
var MentionTypes = map[byte]string{
	'u': "User",
	'g': "Group",
	'c': "Channel",
	'b': "Bot",
}

// markupRegex is the canonical regex used to detect every supported
// markdown construct. The `(?m)` prefix enables MULTILINE mode so
// `^` matches the start of each line (not just the start of the
// string) — this is REQUIRED for the Quote pattern `^(?:> ?...)` to
// match quoted lines that appear in the middle of a message (e.g.
// "📨 hello\n\n> quote\n\nmore text"). Without `(?m)`, Go's regex
// engine only matches `^` at the very start of the string, so
// mid-text quotes are silently ignored and no Quote metadata is
// generated.
var markupRegex = regexp.MustCompile(
	`(?m)(?:^(?:> ?[^\n]*\n?)+)|` + // Quote (group 0)
		"```([\\s\\S]*?)```|" + // Pre (group 1)
		`\*\*([^\n*]+?)\*\*|` + // Bold (group 2)
		"`([^\n`]+?)`|" + // Mono (group 3)
		`__([^\n_]+?)__|` + // Italic (group 4)
		`--([^\n-]+?)--|` + // Underline (group 5)
		`~~([^\n~]+?)~~|` + // Strike (group 6)
		`\|\|([^\n|]+?)\|\||` + // Spoiler (group 7)
		`\[([^\]]+?)\]\((\S+)\)`, // Link (group 8 + group 9 for URL)
)

// Part is the on-wire shape of one formatted region.
type Part struct {
	Type          string `json:"type"`
	FromIndex     int    `json:"from_index"`
	Length        int    `json:"length,omitempty"`
	Language      string `json:"language,omitempty"`
	LinkURL       string `json:"link_url,omitempty"`
	Link          *Link  `json:"link,omitempty"`
	MentionGUID   string `json:"mention_text_object_guid,omitempty"`
	MentionUserID string `json:"mention_text_user_id,omitempty"`
	MentionType   string `json:"mention_text_object_type,omitempty"`
}

// Link is the structured link object attached to a Link-style Part.
type Link struct {
	Type          string         `json:"type"`
	HyperlinkData *HyperlinkData `json:"hyperlink_data,omitempty"`
}

// HyperlinkData holds the URL when the link is a plain hyperlink.
type HyperlinkData struct {
	URL string `json:"url"`
}

// Parts is the outer wrapper around the meta_data_parts slice.
type Parts struct {
	Parts []Part `json:"meta_data_parts"`
}

// Result is the structured result of Format. It is what gets merged
// into the `sendMessage` payload (only Text and Metadata fields).
type Result struct {
	Text     string `json:"text"`
	Metadata *Parts `json:"metadata,omitempty"`
}

// Format converts a Markdown text into the wire-format Result.
//
// Implementation note: the server measures offsets in UTF-16 code units,
// while Go's `regexp.Regexp.FindAllStringSubmatchIndex` returns byte
// indices and Go's string slicing is byte-based — which would corrupt
// multi-byte UTF-8 sequences (e.g. Persian text). To produce exact
// offsets the formatter:
//
// 1. Convert the input to `[]rune` so we can index by codepoint.
// 2. Build a byte→rune index lookup so the regex byte-indices can be
// converted to rune indices.
// 3. Build the UTF-16 prefix array against the *original* text (so
// the `from_index` / `length` values match what the server expects
// for the original text, before any markup is stripped).
func (f Formatter) Format(text string) Result {
	parts := []Part{}
	currentRunes := []rune(text)

	// byteToRune[i] = rune index of byte position i in `text`.
	// Used to translate the regex's byte indices into rune indices
	// so we can slice `currentRunes` correctly.
	byteToRune := buildByteToRuneIndex(text)
	// utf16Prefix[i] = UTF-16 length of text[:i] (i is byte index).
	utf16Prefix := buildUTF16PrefixLengths(text)

	// Accumulated deltas across iterations.
	offset := 0     // UTF-16 length stripped so far
	charOffset := 0 // rune count stripped so far

	for _, match := range markupRegex.FindAllStringSubmatchIndex(text, -1) {
		fullStart, fullEnd := match[0], match[1]
		group := text[fullStart:fullEnd]
		adjustedStart := utf16Prefix[fullStart] - offset

		// Convert the byte-indexed start/end of the matched markup
		// into rune indices so we can slice `currentRunes`.
		runeFullStart := byteToRune[fullStart]
		runeAdjustedCharStart := runeFullStart - charOffset

		for _, mt := range MarkupTypeSequence {
			if !strings.HasPrefix(group, mt.Prefix) {
				continue
			}

			if mt.Type == "Quote" {
				lines := strings.Split(group, "\n")
				contentLines := make([]string, 0, len(lines))
				for _, line := range lines {
					if strings.HasPrefix(line, "> ") {
						contentLines = append(contentLines, line[2:])
					} else if strings.HasPrefix(line, ">") {
						contentLines = append(contentLines, line[1:])
					} else {
						contentLines = append(contentLines, line)
					}
				}
				content := strings.Join(contentLines, "\n")
				inner := f.Format(content)
				content = inner.Text
				contentLength := utf16Length(content)
				if inner.Metadata != nil {
					for _, p := range inner.Metadata.Parts {
						p.FromIndex += adjustedStart
						parts = append(parts, p)
					}
				}
				parts = append(parts, Part{
					Type:      "Quote",
					FromIndex: adjustedStart,
					Length:    contentLength,
				})

				markupLength := utf16Prefix[fullEnd] - utf16Prefix[fullStart]
				// charMarkupLength must be a RUNE count, not a byte
				// count — otherwise the next iteration's adjustedCharStart
				// will be wrong for any multi-byte text.
				charMarkupLength := byteToRune[fullEnd] - byteToRune[fullStart]
				charContentLength := len([]rune(content))

				// Splice: currentRunes[runeAdjustedCharStart:runeFullEnd-charOffset]
				// = content. runeFullEnd-charOffset is the rune index of
				// the byte just past the markup.
				runeFullEnd := byteToRune[fullEnd]
				runeSliceEnd := runeFullEnd - charOffset
				currentRunes = spliceRunes(currentRunes, runeAdjustedCharStart, runeSliceEnd, []rune(content))
				offset += markupLength - contentLength
				charOffset += charMarkupLength - charContentLength
				break
			}

			groupIdx := mt.GroupIdx
			groupStart, groupEnd := match[2*groupIdx], match[2*groupIdx+1]
			var content string
			var contentLength int
			var charContentLength int
			if groupStart == -1 || groupEnd == -1 {
				content = ""
			} else {
				content = text[groupStart:groupEnd]
				contentLength = utf16Prefix[groupEnd] - utf16Prefix[groupStart]
				// charContentLength must be a RUNE count, not a
				// byte count — otherwise charOffset accumulates
				// the wrong delta for multi-byte text and the
				// NEXT iteration's adjustedCharStart is wrong.
				charContentLength = byteToRune[groupEnd] - byteToRune[groupStart]
			}

			// For Pre and Link we don't recurse; everything else does.
			if mt.Type != "Pre" && mt.Type != "Link" {
				inner := f.Format(content)
				content = inner.Text
				contentLength = utf16Length(content)
				charContentLength = len([]rune(content))
				if inner.Metadata != nil {
					for _, p := range inner.Metadata.Parts {
						p.FromIndex += adjustedStart
						parts = append(parts, p)
					}
				}
			}

			p := Part{
				Type:      mt.Type,
				FromIndex: adjustedStart,
				Length:    contentLength,
			}

			switch mt.Type {
			case "Pre":
				// Pre-handling: split the content on the
				// first newline. If the first line is a
				// non-alphanumeric "shebang" (e.g. "#!")
				// AND there is a second line,
				// treat the first line as the language
				// and the rest as the code. Otherwise
				// (first line is alphanumeric like "go",
				// OR there is no second line) treat the
				// whole content as code with no language.
				//
				// The alphanumeric check is critical:
				// the first line is only treated as a
				// language when it does NOT start with a
				// letter or digit (so "go\nfmt.Println()" is
				// treated as code, not as lang="go").
				// Sending language="go" makes the server
				// reject the message with INVALID_INPUT.
				lines := strings.SplitN(content, "\n", 2)
				lang := ""
				if len(lines) > 1 && lines[0] != "" {
					first := strings.TrimSpace(lines[0])
					if first != "" && !isAlnum(first[0]) {
						lang = first
						content = lines[1]
						// Recompute contentLength
						// because content shrank.
						contentLength = utf16Length(content)
						charContentLength = len([]rune(content))
					}
				}
				p.Language = lang
			case "Link":
				urlStart, urlEnd := match[2*groupIdx+2], match[2*groupIdx+3]
				var url string
				if urlStart != -1 && urlEnd != -1 {
					url = text[urlStart:urlEnd]
				}
				if len(url) == 0 {
					p.LinkURL = url
					p.Link = &Link{
						Type:          "hyperlink",
						HyperlinkData: &HyperlinkData{URL: url},
					}
				} else if mentionType, ok := MentionTypes[url[0]]; ok {
					p.Type = "MentionText"
					p.MentionGUID = url
					p.MentionUserID = url
					p.MentionType = mentionType
				} else {
					p.LinkURL = url
					p.Link = &Link{
						Type:          "hyperlink",
						HyperlinkData: &HyperlinkData{URL: url},
					}
				}
			}

			parts = append(parts, p)

			markupLength := utf16Prefix[fullEnd] - utf16Prefix[fullStart]
			// charMarkupLength must be a RUNE count, not a byte
			// count — see the Quote branch above for the same comment.
			charMarkupLength := byteToRune[fullEnd] - byteToRune[fullStart]
			runeFullEnd := byteToRune[fullEnd]
			runeSliceEnd := runeFullEnd - charOffset
			currentRunes = spliceRunes(currentRunes, runeAdjustedCharStart, runeSliceEnd, []rune(content))
			offset += markupLength - contentLength
			charOffset += charMarkupLength - charContentLength
			break
		}
	}

	// TrimSpace the final text. IMPORTANT: the parts' from_index /
	// length were computed against the UNTRIMMED currentRunes, so if
	// TrimSpace removes leading/trailing whitespace, the indices may
	// be off by the amount of leading whitespace removed. We handle
	// this by clamping below.
	rawText := string(currentRunes)
	leadingWS := len(rawText) - len(strings.TrimLeft(rawText, " \t\r\n"))
	result := Result{Text: strings.TrimSpace(rawText)}
	if len(parts) > 0 {
		utf16TextLen := utf16Length(result.Text)
		// Sort parts by from_index so we can detect overlaps.
		sort.SliceStable(parts, func(i, j int) bool {
			return parts[i].FromIndex < parts[j].FromIndex
		})
		cleaned := make([]Part, 0, len(parts))
		var coveredEnd int // end of the last accepted part (exclusive)
		for _, p := range parts {
			// Shift from_index left by the amount of leading
			// whitespace that was trimmed.
			if leadingWS > 0 {
				// Convert leadingWS (byte count) to UTF-16 units.
				// For whitespace (all ASCII) this is 1:1.
				p.FromIndex -= leadingWS
				if p.FromIndex < 0 {
					p.FromIndex = 0
				}
			}
			// Drop parts with negative indices/lengths.
			if p.FromIndex < 0 || p.Length < 0 {
				continue
			}
			// Clamp length so from_index+length doesn't exceed the
			// (now-trimmed) text length.
			if p.FromIndex > utf16TextLen {
				continue
			}
			if p.FromIndex+p.Length > utf16TextLen {
				p.Length = utf16TextLen - p.FromIndex
			}
			// Drop zero-length parts.
			if p.Length == 0 {
				continue
			}
			// Drop parts that overlap with an already-accepted part.
			// The formatter's offset arithmetic can produce overlapping
			// or out-of-order parts for complex inputs (nested markup
			// + emoji + multi-line); overlapping parts make the Rubika
			// server reject the whole message with INVALID_INPUT.
			if p.FromIndex < coveredEnd {
				continue
			}
			cleaned = append(cleaned, p)
			coveredEnd = p.FromIndex + p.Length
		}
		if len(cleaned) > 0 {
			result.Metadata = &Parts{Parts: cleaned}
		}
	}
	return result
}

// spliceRunes replaces runes[start:end] with `content`.
func spliceRunes(runes []rune, start, end int, content []rune) []rune {
	if start < 0 {
		start = 0
	}
	if end > len(runes) {
		end = len(runes)
	}
	if start > end {
		start = end
	}
	// Fast path: in-place replace when lengths match.
	if len(content) == end-start {
		copy(runes[start:end], content)
		return runes
	}
	// General path: build a new slice.
	out := make([]rune, 0, len(runes)-(end-start)+len(content))
	out = append(out, runes[:start]...)
	out = append(out, content...)
	out = append(out, runes[end:]...)
	return out
}

// buildByteToRuneIndex returns a slice `idx` such that `idx[bytePos]`
// is the rune index of `bytePos` in `text`. This lets us translate
// the byte indices returned by the regex into the rune indices we
// need to slice `[]rune`.
func buildByteToRuneIndex(text string) []int {
	idx := make([]int, len(text)+1)
	runePos := 0
	for bytePos := range text {
		idx[bytePos] = runePos
		runePos++
	}
	idx[len(text)] = runePos
	return idx
}

// ToMarkdown converts the most common HTML formatting tags to their
// Markdown equivalents. If the input does not contain any HTML tag it
// is returned unchanged (after trimming whitespace).
func (f Formatter) ToMarkdown(text string) string {
	if !strings.ContainsAny(text, "<>") {
		return strings.TrimSpace(text)
	}
	return toMarkdownHTML(text)
}

// JSON returns the JSON representation of the result, useful when
// merging it into a larger payload.
func (r Result) JSON() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// buildUTF16PrefixLengths returns a slice where prefix[i] is the UTF-16
// length of text[:i]. Codepoints >0xFFFF take 2 UTF-16 code units.
func buildUTF16PrefixLengths(text string) []int {
	prefix := make([]int, len(text)+1)
	total := 0
	for i, r := range text {
		total += utf16LengthOfRune(r)
		prefix[i+1] = total
	}
	return prefix
}

func utf16LengthOfRune(r rune) int {
	if r > 0xFFFF {
		return 2
	}
	return 1
}

func utf16Length(s string) int {
	n := 0
	for _, r := range s {
		n += utf16LengthOfRune(r)
	}
	return n
}

// isAlnum reports whether the first byte of `b` is an ASCII letter or
// digit. Used by the Pre-handling logic to decide whether the first
// line of a code block is a language tag or part of the code.
func isAlnum(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}
