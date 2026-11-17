// Package antilink is a plugin for groups and channels: it deletes
// messages that contain links (or @mentions) sent by people who are not
// administrators of that chat.
//
// Configuration:
//
//	allowed_domains  []string  domains that are always fine
//	block_mentions   bool      also remove @username mentions (default false)
//	warn_message     string    reply sent after a removal ("" = silent)
package antilink

import (
	"context"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

var (
	linkRe    = regexp.MustCompile(`(?i)(?:https?://|www\.)[^\s]+|\b[a-z0-9-]+(?:\.[a-z0-9-]+)*\.(?:ir|com|net|org|me|io|info|xyz|ly|app|site|online)\b[^\s]*`)
	mentionRe = regexp.MustCompile(`@[A-Za-z0-9_]{3,}`)
)

// Plugin implements plugins.Plugin.
type Plugin struct {
	plugins.Base
	mu     sync.Mutex
	admins map[string]time.Time // "chat|user" -> expiry of a positive admin check
	now    func() time.Time
}

// New creates the plugin.
func New() *Plugin {
	return &Plugin{Base: plugins.NewBase(plugins.Meta{
		Name: "antilink", Version: "1.0.0", Priority: 35,
		Description: "Removes links sent by non-admins in groups",
	}), now: time.Now}
}

// DefaultConfig implements plugins.Defaulter.
func (p *Plugin) DefaultConfig() plugins.Config {
	return plugins.Config{"block_mentions": false, "warn_message": "🚫 ارسال لینک در این گروه مجاز نیست."}
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	cfg := h.Config()
	allowed, mentions, warn := cfg.Strings("allowed_domains"), cfg.Bool("block_mentions", false), cfg.String("warn_message", "")
	p.mu.Lock()
	p.admins = map[string]time.Time{}
	p.mu.Unlock()

	h.Use(func(ctx context.Context, b *bot.Bot, u *bot.Update, next func() error) error {
		if u.IsPrivateChat() || u.Message() == nil || !Contains(u.Text(), allowed, mentions) {
			return next()
		}
		uid := u.SenderID()
		if plugins.IsAdmin(h.Manager(), uid) || p.chatAdmin(ctx, b, u.ChatID, uid) {
			return next()
		}
		if err := b.DeleteMessage(ctx, u.ChatID, u.MessageID()); err != nil {
			h.Warnf("cannot delete message %s in %s: %v", u.MessageID(), u.ChatID, err)
			return next()
		}
		if warn != "" {
			_, _ = b.SendMessage(ctx, u.ChatID, warn)
		}
		return nil
	})
	return nil
}

// Contains reports whether text holds a link (outside allowedDomains)
// or, when mentions is set, an @mention.
func Contains(text string, allowedDomains []string, mentions bool) bool {
	if mentions && mentionRe.MatchString(text) {
		return true
	}
	for _, m := range linkRe.FindAllString(text, -1) {
		l := strings.ToLower(m)
		host := linkHost(l)
		ok := false
		for _, d := range allowedDomains {
			d = strings.ToLower(strings.TrimSpace(d))
			if d == "" {
				continue
			}
			// Exact host match OR a subdomain of an allowed domain.
			// strings.Contains would let "example.com.evil.com" or
			// "notexample.com" bypass the filter for allowed domain
			// "example.com"; the suffix check requires a real dot
			// boundary.
			if host == d || strings.HasSuffix(host, "."+d) {
				ok = true
				break
			}
		}
		if !ok {
			return true
		}
	}
	return false
}

// linkHost extracts the host portion of a matched link. The linkRe
// tokens may be bare domains ("example.com/path") or full URLs
// ("https://example.com/x"); url.Parse leaves the host empty for the
// bare form, in which case the leading path segment is the host.
func linkHost(l string) string {
	// Ensure there is a scheme so url.Parse populates Host.
	if !strings.Contains(l, "://") {
		l = "http://" + l
	}
	u, err := url.Parse(l)
	if err != nil {
		return strings.ToLower(strings.TrimSpace(l))
	}
	host := u.Hostname()
	if host == "" {
		host = strings.Trim(u.Path, "/")
		if i := strings.IndexByte(host, '/'); i >= 0 {
			host = host[:i]
		}
	}
	return strings.ToLower(host)
}

// chatAdmin asks Rubika whether uid administers the chat; positive
// answers are cached for five minutes. API failures count as "not admin".
func (p *Plugin) chatAdmin(ctx context.Context, b *bot.Bot, chat, uid string) bool {
	key := chat + "|" + uid
	p.mu.Lock()
	if exp, ok := p.admins[key]; ok && p.now().Before(exp) {
		p.mu.Unlock()
		return true
	}
	p.mu.Unlock()
	ok, err := b.MemberIsAdmin(ctx, chat, uid)
	if err != nil || !ok {
		return false
	}
	p.mu.Lock()
	p.admins[key] = p.now().Add(5 * time.Minute)
	p.mu.Unlock()
	return true
}
