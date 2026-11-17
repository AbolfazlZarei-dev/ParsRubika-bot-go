package rules

import (
	"context"
	"regexp"
	"strings"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// Text matches the message text either by exact string or by regex.
// When `Pattern` is empty it matches any text message.
type Text struct {
	Pattern string
	Regex   bool
	re      *regexp.Regexp
}

// NewText constructs a Text rule.
func NewText(pattern string, regex bool) *Text {
	f := &Text{Pattern: pattern, Regex: regex}
	if regex && pattern != "" {
		f.re, _ = regexp.Compile(pattern)
	}
	return f
}

// Check implements Rule.
func (f *Text) Check(_ context.Context, upd *models.Update) bool {
	t := upd.Text()
	if t == "" {
		return false
	}
	if f.Pattern == "" {
		return true
	}
	if f.Regex {
		return f.re != nil && f.re.MatchString(t)
	}
	return t == f.Pattern
}

// Command matches a bot command (e.g. `/start`, `!help`).
type Command struct {
	Commands      []string
	Prefixes      []string
	CaseSensitive bool
	AllowNoPrefix bool
}

// NewCommand constructs a Command rule.
func NewCommand(cmds, prefixes []string, caseInsensitive, allowNoPrefix bool) *Command {
	if len(prefixes) == 0 {
		prefixes = []string{"/"}
	}
	c := &Command{
		Commands:      append([]string{}, cmds...),
		Prefixes:      prefixes,
		CaseSensitive: !caseInsensitive,
		AllowNoPrefix: allowNoPrefix,
	}
	if caseInsensitive {
		for i, s := range c.Commands {
			c.Commands[i] = strings.ToLower(s)
		}
	}
	return c
}

// sameCommand reports whether word is exactly want, or want addressed to
// a specific bot ("/start@my_bot"). Longer words such as "/startup" or
// "/warnings" must not match "/start" or "/warn".
func sameCommand(word, want string) bool {
	if word == want {
		return true
	}
	return strings.HasPrefix(word, want+"@")
}

// Check implements Rule.
func (c *Command) Check(_ context.Context, upd *models.Update) bool {
	t := upd.Text()
	if t == "" {
		return false
	}
	parts := strings.Fields(t)
	if len(parts) == 0 {
		return false
	}
	cmd := parts[0]
	cmp := cmd
	if !c.CaseSensitive {
		cmp = strings.ToLower(cmd)
	}
	for _, cc := range c.Commands {
		for _, p := range c.Prefixes {
			if sameCommand(cmp, p+cc) {
				return true
			}
		}
		if c.AllowNoPrefix && sameCommand(cmp, cc) {
			return true
		}
	}
	return false
}
