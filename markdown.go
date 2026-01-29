package ParsRubika

// نسخه: 2.5.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/AbolfazlZarei-dev/ParsRubika-bot-go

import (
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// Style نوع استایل‌های مختلف متنی (Markdown و HTML) را مشخص می‌کند
type Style string

const (
	// استایل‌های Markdown
	StylePre        Style = "Pre"
	StyleLink       Style = "Link"
	StyleMention    Style = "Mention"
	StyleCodeInline Style = "CodeInline"
	StyleSpoiler    Style = "Spoiler"
	StyleBold       Style = "Bold"
	StyleStrike     Style = "Strike"
	StyleUnderline  Style = "Underline"
	StyleItalic     Style = "Italic"
	StyleQuote      Style = "Quote"

	// استایل‌های HTML
	StylePreHTML        Style = "PreHTML"
	StyleHTMLLink       Style = "HTMLLink"
	StyleMentionHTML    Style = "MentionHTML"
	StyleCodeInlineHTML Style = "CodeInlineHTML"
	StyleSpoilerHTML    Style = "SpoilerHTML"
	StyleBoldHTML       Style = "BoldHTML"
	StyleItalicHTML     Style = "ItalicHTML"
	StyleStrikeHTML     Style = "StrikeHTML"
	StyleUnderlineHTML  Style = "UnderlineHTML"
)

// Match اطلاعات یک تطابق پیدا شده در متن را نگه می‌دارد
type Match struct {
	Start     int    `json:"start"`           // محل شروع در متن
	End       int    `json:"end"`             // محل پایان در متن
	Style     Style  `json:"style"`           // نوع استایل
	Content   string `json:"content"`         // محتوای استخراج‌شده
	FullMatch string `json:"full_match"`      // متن کامل شامل علائم
	Extra     string `json:"extra,omitempty"` // اطلاعات اضافی (مثل لینک یا شناسه کاربر)
	Priority  int    `json:"priority"`        // اولویت پردازش
}

// MarkdownMetadata ساختار متادیتای نهایی برای ارسال به API
type MarkdownMetadata struct {
	Type              string `json:"type"`
	FromIndex         int    `json:"from_index"`
	Length            int    `json:"length"`
	LinkURL           string `json:"link_url,omitempty"`
	MentionTextUserID string `json:"mention_text_user_id,omitempty"`
}

// MarkdownParseResult خروجی نهایی پردازش متن
type MarkdownParseResult struct {
	Metadata []MarkdownMetadata `json:"metadata"`
	Text     string             `json:"text"`
}

// الگوهای Markdown
var MarkdownPatterns = map[Style][]string{
	StylePre:        {"```(?:[^\\n]*\\n)?([\\s\\S]*?)```"},
	StyleLink:       {"\\[([^\\]]+?)\\]\\((https?://[^)]+)\\)"},
	StyleMention:    {"\\(([^)]+?)\\)\\[([^\\]]+?)\\]"},
	StyleCodeInline: {"`([^`]+?)`"},
	StyleSpoiler:    {"\\|\\|([^|]+?)\\|\\|"},
	StyleBold:       {"\\*\\*([^\\*]+?)\\*\\*"},
	StyleStrike:     {"~~([^~]+?)~~"},
	StyleUnderline:  {"__([^_]+?)__"},
	StyleItalic:     {"__([^_]+?)__"},
	StyleQuote:      {"××([^×]+?)××"},
}

// الگوهای HTML
var HTMLPatterns = map[Style][]string{
	StylePreHTML:        {"<pre>([\\s\\S]*?)</pre>"},
	StyleHTMLLink:       {`<a\s+href="([^"]+?)">([^<]+?)</a>`},
	StyleMentionHTML:    {`<mention\s+objectId="([^"]+?)">([^<]+?)</mention>`},
	StyleCodeInlineHTML: {"<code>([^<]+?)</code>"},
	StyleSpoilerHTML:    {`<span\s+class="tg-spoiler">([^<]+?)</span>`},
	StyleBoldHTML:       {"<b>([^<]+?)</b>", "<strong>([^<]+?)</strong>"},
	StyleItalicHTML:     {"<i>([^<]+?)</i>", "<em>([^<]+?)</em>"},
	StyleStrikeHTML:     {"<s>([^<]+?)</s>", "<del>([^<]+?)</del>"},
	StyleUnderlineHTML:  {"<u>([^<]+?)</u>"},
}

// اولویت پردازش استایل‌های Markdown
var MarkdownPriority = map[Style]int{
	StylePre:        0,
	StyleLink:       1,
	StyleMention:    1,
	StyleCodeInline: 2,
	StyleSpoiler:    3,
	StyleBold:       4,
	StyleStrike:     5,
	StyleUnderline:  6,
	StyleItalic:     7,
	StyleQuote:      8,
}

// اولویت پردازش استایل‌های HTML
var HTMLPriority = map[Style]int{
	StylePreHTML:        0,
	StyleHTMLLink:       1,
	StyleMentionHTML:    1,
	StyleCodeInlineHTML: 2,
	StyleSpoilerHTML:    3,
	StyleBoldHTML:       4,
	StyleItalicHTML:     5,
	StyleStrikeHTML:     6,
	StyleUnderlineHTML:  7,
}

// collectMatches تمام تطابق‌های موجود در متن را استخراج می‌کند
func collectMatches(text string, patterns map[Style][]string, priority map[Style]int) []Match {
	var matches []Match

	for style, pats := range patterns {
		for _, pat := range pats {
			re := regexp.MustCompile(pat)
			allMatches := re.FindAllStringSubmatchIndex(text, -1)

			for _, m := range allMatches {
				start, end := m[0], m[1]
				var groups []string

				for i := 2; i < len(m); i += 2 {
					if m[i] != -1 && m[i+1] != -1 {
						groups = append(groups, text[m[i]:m[i+1]])
					}
				}

				content := ""
				extra := ""

				switch style {
				case StyleLink:
					if len(groups) >= 2 {
						content = groups[0]
						extra = groups[1]
					}
				case StyleMention:
					if len(groups) >= 2 {
						extra = groups[0]   // شناسه کاربر
						content = groups[1] // متن نمایشی
					}
				case StyleHTMLLink, StyleMentionHTML:
					if len(groups) >= 2 {
						extra = groups[0]
						content = groups[1]
					}
				default:
					if len(groups) >= 1 {
						content = groups[0]
					} else {
						content = text[start:end]
					}
				}

				if style == StylePre || style == StylePreHTML {
					content = strings.Trim(content, "\n")
				}

				matches = append(matches, Match{
					Start:     start,
					End:       end,
					Style:     style,
					Content:   content,
					FullMatch: text[start:end],
					Extra:     extra,
					Priority:  priority[style],
				})
			}
		}
	}

	return matches
}

// allowMatch بررسی می‌کند که یک تطابق جدید با تطابق‌های قبلی تداخل نامعتبر نداشته باشد
func allowMatch(chosen []Match, candidate Match) bool {
	s, e := candidate.Start, candidate.End

	for _, c := range chosen {
		os, oe := c.Start, c.End
		// اگر یکی کاملاً داخل دیگری باشد، مجاز است
		if (s >= os && e <= oe) || (os >= s && oe <= e) {
			continue
		}
		// اگر هم‌پوشانی ناقص وجود داشته باشد، غیرمجاز است
		if !(e <= os || s >= oe) {
			return false
		}
	}

	return true
}

// pickMatchesAllowingNested تطابق‌های نهایی را با اجازه تو در تو بودن انتخاب می‌کند
func pickMatchesAllowingNested(matches []Match) []Match {
	// مرتب‌سازی بر اساس اولویت و محل شروع
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Priority != matches[j].Priority {
			return matches[i].Priority < matches[j].Priority
		}
		return matches[i].Start < matches[j].Start
	})

	var chosen []Match
	for _, m := range matches {
		if allowMatch(chosen, m) {
			chosen = append(chosen, m)
		}
	}

	// مرتب‌سازی نهایی بر اساس محل شروع
	sort.Slice(chosen, func(i, j int) bool {
		return chosen[i].Start < chosen[j].Start
	})

	return chosen
}

// ParseMarkdown متن Markdown را پردازش کرده و متادیتا تولید می‌کند
func ParseMarkdown(text string) MarkdownParseResult {
	if text == "" {
		return MarkdownParseResult{Metadata: []MarkdownMetadata{}, Text: ""}
	}

	allMatches := collectMatches(text, MarkdownPatterns, MarkdownPriority)
	chosen := pickMatchesAllowingNested(allMatches)

	var metadata []MarkdownMetadata
	var outParts []string
	last := 0
	outLength := 0

	for _, m := range chosen {
		if last < m.Start {
			plainPart := text[last:m.Start]
			outParts = append(outParts, plainPart)
			outLength += utf8.RuneCountInString(plainPart)
		}

		currentIndex := outLength
		content := m.Content
		outParts = append(outParts, content)
		length := utf8.RuneCountInString(content)
		outLength += length

		st := m.Style

		switch st {
		case StyleLink:
			metadata = append(metadata, MarkdownMetadata{
				Type:      "Link",
				FromIndex: currentIndex,
				Length:    length,
				LinkURL:   m.Extra,
			})
		case StyleMention:
			metadata = append(metadata, MarkdownMetadata{
				Type:              "MentionText",
				FromIndex:         currentIndex,
				Length:            length,
				MentionTextUserID: m.Extra,
			})
		case StyleCodeInline:
			metadata = append(metadata, MarkdownMetadata{
				Type:      "Mono",
				FromIndex: currentIndex,
				Length:    length,
			})
		case StylePre:
			metadata = append(metadata, MarkdownMetadata{
				Type:      "Pre",
				FromIndex: currentIndex,
				Length:    length,
			})
		case StyleQuote:
			metadata = append(metadata, MarkdownMetadata{
				Type:      "Quote",
				FromIndex: currentIndex,
				Length:    length,
			})
		default:
			mapTypes := map[Style]string{
				StyleBold:      "Bold",
				StyleItalic:    "Italic",
				StyleUnderline: "Underline",
				StyleStrike:    "Strike",
				StyleSpoiler:   "Spoiler",
			}
			metaType := mapTypes[st]
			if metaType == "" {
				metaType = string(st)
			}
			metadata = append(metadata, MarkdownMetadata{
				Type:      metaType,
				FromIndex: currentIndex,
				Length:    length,
			})
		}

		last = m.End
	}

	if last < len(text) {
		outParts = append(outParts, text[last:])
	}

	return MarkdownParseResult{
		Metadata: metadata,
		Text:     strings.Join(outParts, ""),
	}
}

// ParseHTML متن HTML را پردازش کرده و متادیتا تولید می‌کند
func ParseHTML(text string) MarkdownParseResult {
	if text == "" {
		return MarkdownParseResult{Metadata: []MarkdownMetadata{}, Text: ""}
	}

	allMatches := collectMatches(text, HTMLPatterns, HTMLPriority)
	chosen := pickMatchesAllowingNested(allMatches)

	var metadata []MarkdownMetadata
	var outParts []string
	last := 0
	outLength := 0

	for _, m := range chosen {
		if last < m.Start {
			plainPart := text[last:m.Start]
			outParts = append(outParts, plainPart)
			outLength += utf8.RuneCountInString(plainPart)
		}

		currentIndex := outLength
		content := m.Content
		outParts = append(outParts, content)
		length := utf8.RuneCountInString(content)
		outLength += length

		st := m.Style

		switch st {
		case StyleHTMLLink:
			url := m.Extra
			if strings.HasPrefix(url, "rubika://") {
				metadata = append(metadata, MarkdownMetadata{
					Type:              "MentionText",
					FromIndex:         currentIndex,
					Length:            length,
					MentionTextUserID: strings.Replace(url, "rubika://", "", 1),
				})
			} else {
				metadata = append(metadata, MarkdownMetadata{
					Type:      "Link",
					FromIndex: currentIndex,
					Length:    length,
					LinkURL:   url,
				})
			}
		case StyleMentionHTML:
			metadata = append(metadata, MarkdownMetadata{
				Type:              "MentionText",
				FromIndex:         currentIndex,
				Length:            length,
				MentionTextUserID: m.Extra,
			})
		case StyleCodeInlineHTML:
			metadata = append(metadata, MarkdownMetadata{
				Type:      "Mono",
				FromIndex: currentIndex,
				Length:    length,
			})
		case StylePreHTML:
			metadata = append(metadata, MarkdownMetadata{
				Type:      "Pre",
				FromIndex: currentIndex,
				Length:    length,
			})
		case StyleSpoilerHTML:
			metadata = append(metadata, MarkdownMetadata{
				Type:      "Spoiler",
				FromIndex: currentIndex,
				Length:    length,
			})
		default:
			mapHTML := map[Style]string{
				StyleBoldHTML:      "Bold",
				StyleItalicHTML:    "Italic",
				StyleUnderlineHTML: "Underline",
				StyleStrikeHTML:    "Strike",
			}
			metaType := mapHTML[st]
			if metaType == "" {
				metaType = string(st)
			}
			metadata = append(metadata, MarkdownMetadata{
				Type:      metaType,
				FromIndex: currentIndex,
				Length:    length,
			})
		}

		last = m.End
	}

	if last < len(text) {
		outParts = append(outParts, text[last:])
	}

	return MarkdownParseResult{
		Metadata: metadata,
		Text:     strings.Join(outParts, ""),
	}
}

// TextParser یک Wrapper ساده برای استفاده از پارسرهاست
type TextParser struct{}

// CheckMarkdown متن Markdown را بررسی و پردازش می‌کند
func (tp *TextParser) CheckMarkdown(text string) (MarkdownParseResult, error) {
	return ParseMarkdown(text), nil
}

// CheckHTML متن HTML را بررسی و پردازش می‌کند
func (tp *TextParser) CheckHTML(text string) (MarkdownParseResult, error) {
	return ParseHTML(text), nil
}
