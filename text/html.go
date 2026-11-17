package text

import (
	"regexp"
	"strings"
)

// A minimal HTML → Markdown converter. Go's RE2 engine does not support
// backreferences, so one regex per tag is used instead of a single
// `<(\w+)>...</\1>` pattern.

var (
	htmlBold      = regexp.MustCompile(`(?is)<b[^>]*>(.*?)</\s*b>`)
	htmlStrong    = regexp.MustCompile(`(?is)<strong[^>]*>(.*?)</\s*strong>`)
	htmlItalic    = regexp.MustCompile(`(?is)<i[^>]*>(.*?)</\s*i>`)
	htmlEm        = regexp.MustCompile(`(?is)<em[^>]*>(.*?)</\s*em>`)
	htmlUnderline = regexp.MustCompile(`(?is)<u[^>]*>(.*?)</\s*u>`)
	htmlStrikeS   = regexp.MustCompile(`(?is)<s[^>]*>(.*?)</\s*s>`)
	htmlStrikeDel = regexp.MustCompile(`(?is)<del[^>]*>(.*?)</\s*del>`)
	htmlStrikeStr = regexp.MustCompile(`(?is)<strike[^>]*>(.*?)</\s*strike>`)
	htmlCode      = regexp.MustCompile(`(?is)<code[^>]*>(.*?)</\s*code>`)
	htmlPre       = regexp.MustCompile(`(?is)<pre[^>]*>(.*?)</\s*pre>`)
	htmlAnchor    = regexp.MustCompile(`(?is)<a[^>]*href=["']([^"']+)["'][^>]*>(.*?)</\s*a>`)
	htmlQuote     = regexp.MustCompile(`(?is)<blockquote[^>]*>(.*?)</\s*blockquote>`)
	htmlBr        = regexp.MustCompile(`(?is)<br\s*/?>`)
	htmlParagraph = regexp.MustCompile(`(?is)<p[^>]*>`)
	htmlClosePara = regexp.MustCompile(`(?is)</\s*p>`)
)

// toMarkdownHTML converts the most common HTML formatting tags to their
// Markdown equivalents. It is a small subset — for anything more
// elaborate users should pre-convert their HTML.
func toMarkdownHTML(html string) string {
	out := html
	out = htmlPre.ReplaceAllStringFunc(out, func(m string) string {
		s := htmlPre.FindStringSubmatch(m)
		if len(s) < 2 {
			return m
		}
		body := strings.TrimSpace(s[1])
		// Detect language from a leading <code class="language-x">.
		lang := ""
		if idx := strings.Index(body, "<code"); idx >= 0 {
			end := strings.Index(body[idx:], ">")
			if end > 0 {
				tag := body[idx : idx+end]
				if cls := strings.Index(tag, "language-"); cls >= 0 {
					rest := tag[cls+len("language-"):]
					if q := strings.IndexAny(rest, "\"' "); q > 0 {
						lang = rest[:q]
					}
				}
			}
		}
		return "\n```" + lang + "\n" + stripTags(body) + "\n```\n"
	})
	out = htmlBold.ReplaceAllString(out, "**${1}**")
	out = htmlStrong.ReplaceAllString(out, "**${1}**")
	out = htmlItalic.ReplaceAllString(out, "__${1}__")
	out = htmlEm.ReplaceAllString(out, "__${1}__")
	out = htmlUnderline.ReplaceAllString(out, "--${1}--")
	out = htmlStrikeS.ReplaceAllString(out, "~~${1}~~")
	out = htmlStrikeDel.ReplaceAllString(out, "~~${1}~~")
	out = htmlStrikeStr.ReplaceAllString(out, "~~${1}~~")
	out = htmlCode.ReplaceAllString(out, "`${1}`")
	out = htmlAnchor.ReplaceAllString(out, "[${2}](${1})")
	out = htmlQuote.ReplaceAllStringFunc(out, func(m string) string {
		s := htmlQuote.FindStringSubmatch(m)
		if len(s) < 2 {
			return m
		}
		lines := strings.Split(strings.TrimSpace(s[1]), "\n")
		for i, line := range lines {
			lines[i] = "> " + line
		}
		return strings.Join(lines, "\n")
	})
	out = htmlBr.ReplaceAllString(out, "\n")
	out = htmlParagraph.ReplaceAllString(out, "")
	out = htmlClosePara.ReplaceAllString(out, "\n\n")
	return strings.TrimSpace(out)
}

// stripTags removes every <...> tag from the supplied string.
func stripTags(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return b.String()
}
