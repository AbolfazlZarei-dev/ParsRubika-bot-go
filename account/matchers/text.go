package matchers

import (
	"context"
	"regexp"
	"strings"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// Text matches the message text. When `pattern` is empty it matches
// any text; when `regex` is true the pattern is treated as a Go regex.
type Text struct {
	pattern string
	regex   bool
	re      *regexp.Regexp
	kind    string
}

// NewText constructs a Text matcher. `kind` controls which update type
// it accepts (defaults to "MessageUpdates").
func NewText(pattern string, regex bool, kind string) *Text {
	if kind == "" {
		kind = "MessageUpdates"
	}
	t := &Text{pattern: pattern, regex: regex, kind: kind}
	if regex && pattern != "" {
		t.re, _ = regexp.Compile(pattern)
	}
	return t
}

// Kind implements Matcher.
func (t *Text) Kind() string { return t.kind }

// Match implements Matcher.
func (t *Text) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	txt := ev.Text()
	if txt == "" {
		return false
	}
	if t.pattern == "" {
		return true
	}
	if t.regex {
		return t.re != nil && t.re.MatchString(txt)
	}
	return txt == t.pattern
}

// Command matches a bot command (e.g. `/start`, `!help`).
type Command struct {
	commands        []string
	prefixes        []string
	caseInsensitive bool
	allowNoPrefix   bool
	kind            string
}

// NewCommand constructs a Command matcher.
func NewCommand(cmds, prefixes []string, caseInsensitive, allowNoPrefix bool, kind string) *Command {
	if len(prefixes) == 0 {
		prefixes = []string{"/"}
	}
	if kind == "" {
		kind = "MessageUpdates"
	}
	c := &Command{
		commands:        append([]string{}, cmds...),
		prefixes:        prefixes,
		caseInsensitive: caseInsensitive,
		allowNoPrefix:   allowNoPrefix,
		kind:            kind,
	}
	if caseInsensitive {
		for i, s := range c.commands {
			c.commands[i] = strings.ToLower(s)
		}
	}
	return c
}

// Kind implements Matcher.
func (c *Command) Kind() string { return c.kind }

// Match implements Matcher.
func (c *Command) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	t := ev.Text()
	if t == "" {
		return false
	}
	parts := strings.Fields(t)
	if len(parts) == 0 {
		return false
	}
	cmd := parts[0]
	cmp := cmd
	if c.caseInsensitive {
		cmp = strings.ToLower(cmd)
	}
	for _, cc := range c.commands {
		for _, p := range c.prefixes {
			if cmp == p+cc || strings.HasPrefix(cmp, p+cc) {
				return true
			}
		}
		if c.allowNoPrefix && (cmp == cc || strings.HasPrefix(cmp, cc)) {
			return true
		}
	}
	return false
}

// ObjectGUID matches by chat GUID.
type ObjectGUID struct {
	guids []string
	kind  string
}

// NewObjectGUID constructs an ObjectGUID matcher.
func NewObjectGUID(kind string, guids ...string) *ObjectGUID {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &ObjectGUID{guids: guids, kind: kind}
}

// Kind implements Matcher.
func (o *ObjectGUID) Kind() string { return o.kind }

// Match implements Matcher.
func (o *ObjectGUID) Match(_ context.Context, raw map[string]any) bool {
	g := asEvent(raw).ObjectGUID()
	for _, want := range o.guids {
		if g == want {
			return true
		}
	}
	return false
}

// AuthorGUID matches by author GUID.
type AuthorGUID struct {
	guids []string
	kind  string
}

// NewAuthorGUID constructs an AuthorGUID matcher.
func NewAuthorGUID(kind string, guids ...string) *AuthorGUID {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &AuthorGUID{guids: guids, kind: kind}
}

// Kind implements Matcher.
func (a *AuthorGUID) Kind() string { return a.kind }

// Match implements Matcher.
func (a *AuthorGUID) Match(_ context.Context, raw map[string]any) bool {
	g := asEvent(raw).AuthorGUID()
	for _, want := range a.guids {
		if g == want {
			return true
		}
	}
	return false
}

// IsGroup matches when the update came from a group.
type IsGroup struct{ kind string }

// NewIsGroup constructs an IsGroup matcher.
func NewIsGroup(kind string) *IsGroup {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsGroup{kind: kind}
}

// Kind implements Matcher.
func (g *IsGroup) Kind() string { return g.kind }

// Match implements Matcher.
func (g *IsGroup) Match(_ context.Context, raw map[string]any) bool {
	return asEvent(raw).IsGroup()
}

// IsPrivate matches when the update came from a private chat.
type IsPrivate struct{ kind string }

// NewIsPrivate constructs an IsPrivate matcher.
func NewIsPrivate(kind string) *IsPrivate {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsPrivate{kind: kind}
}

// Kind implements Matcher.
func (p *IsPrivate) Kind() string { return p.kind }

// Match implements Matcher.
func (p *IsPrivate) Match(_ context.Context, raw map[string]any) bool {
	return asEvent(raw).IsPrivate()
}

// IsChannel matches when the update came from a channel.
type IsChannel struct{ kind string }

// NewIsChannel constructs an IsChannel matcher.
func NewIsChannel(kind string) *IsChannel {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsChannel{kind: kind}
}

// Kind implements Matcher.
func (c *IsChannel) Kind() string { return c.kind }

// Match implements Matcher.
func (c *IsChannel) Match(_ context.Context, raw map[string]any) bool {
	return asEvent(raw).IsChannel()
}

// IsText matches text messages (any non-empty text).
type IsText struct{ kind string }

// NewIsText constructs an IsText matcher.
func NewIsText(kind string) *IsText {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsText{kind: kind}
}

// Kind implements Matcher.
func (t *IsText) Kind() string { return t.kind }

// Match implements Matcher.
func (t *IsText) Match(_ context.Context, raw map[string]any) bool {
	return asEvent(raw).Text() != ""
}

// IsFile matches messages with a file_inline of type "File".
type IsFile struct{ kind string }

// NewIsFile constructs an IsFile matcher.
func NewIsFile(kind string) *IsFile {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsFile{kind: kind}
}

// Kind implements Matcher.
func (f *IsFile) Kind() string { return f.kind }

// Match implements Matcher.
func (f *IsFile) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	return ev.IsFileInline() && ev.FileInline() != nil && ev.FileInline().Type() == "File"
}

// IsPhoto matches messages with a file_inline of type "Image".
type IsPhoto struct{ kind string }

// NewIsPhoto constructs an IsPhoto matcher.
func NewIsPhoto(kind string) *IsPhoto {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsPhoto{kind: kind}
}

// Kind implements Matcher.
func (p *IsPhoto) Kind() string { return p.kind }

// Match implements Matcher.
func (p *IsPhoto) Match(_ context.Context, raw map[string]any) bool {
	ev := asEvent(raw)
	return ev.IsFileInline() && ev.FileInline() != nil && ev.FileInline().Type() == "Image"
}

// IsSticker matches sticker messages.
type IsSticker struct{ kind string }

// NewIsSticker constructs an IsSticker matcher.
func NewIsSticker(kind string) *IsSticker {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsSticker{kind: kind}
}

// Kind implements Matcher.
func (s *IsSticker) Kind() string { return s.kind }

// Match implements Matcher.
func (s *IsSticker) Match(_ context.Context, raw map[string]any) bool {
	return asEvent(raw).Sticker() != nil
}

// IsForward matches forwarded messages.
type IsForward struct{ kind string }

// NewIsForward constructs an IsForward matcher.
func NewIsForward(kind string) *IsForward {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsForward{kind: kind}
}

// Kind implements Matcher.
func (f *IsForward) Kind() string { return f.kind }

// Match implements Matcher.
func (f *IsForward) Match(_ context.Context, raw map[string]any) bool {
	return asEvent(raw).IsForward()
}

// IsEdited matches edited messages.
type IsEdited struct{ kind string }

// NewIsEdited constructs an IsEdited matcher.
func NewIsEdited(kind string) *IsEdited {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsEdited{kind: kind}
}

// Kind implements Matcher.
func (e *IsEdited) Kind() string { return e.kind }

// Match implements Matcher.
func (e *IsEdited) Match(_ context.Context, raw map[string]any) bool {
	return asEvent(raw).IsEdited()
}

// IsMe matches updates where the author is the current user. The
// current user's GUID must be supplied at construction time.
type IsMe struct {
	selfGUID string
	kind     string
}

// NewIsMe constructs an IsMe matcher.
func NewIsMe(selfGUID, kind string) *IsMe {
	if kind == "" {
		kind = "MessageUpdates"
	}
	return &IsMe{selfGUID: selfGUID, kind: kind}
}

// Kind implements Matcher.
func (m *IsMe) Kind() string { return m.kind }

// Match implements Matcher.
func (m *IsMe) Match(_ context.Context, raw map[string]any) bool {
	return asEvent(raw).IsMe(m.selfGUID)
}

// Ensure imports are used.
var _ = events.New
