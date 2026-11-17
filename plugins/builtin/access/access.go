// Package access is a plugin that controls who may talk to the bot:
// administrators, an optional allow-list and a block-list. Other
// plugins find it through the "access" service.
//
// Configuration:
//
//	admins        []string  user IDs with administrator rights
//	allowed       []string  when non-empty, only these users (and admins) get through
//	blocked       []string  user IDs that are always ignored
//	deny_message  string    reply sent to rejected users ("" = stay silent)
package access

import (
	"context"
	"sync"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

// Plugin implements plugins.Plugin and plugins.AccessService.
type Plugin struct {
	plugins.Base
	mu                       sync.RWMutex
	admins, allowed, blocked map[string]struct{}
	deny                     string
	notified                 map[string]struct{}
}

// New creates the plugin.
func New() *Plugin {
	return &Plugin{Base: plugins.NewBase(plugins.Meta{
		Name: "access", Version: "1.0.0", Priority: 10,
		Description: "Admins, allow-list and block-list",
	})}
}

func set(ids []string) map[string]struct{} {
	m := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		m[id] = struct{}{}
	}
	return m
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	cfg := h.Config()
	p.mu.Lock()
	p.admins, p.allowed, p.blocked = set(cfg.Strings("admins")), set(cfg.Strings("allowed")), set(cfg.Strings("blocked"))
	p.deny = cfg.String("deny_message", "")
	p.notified = map[string]struct{}{}
	p.mu.Unlock()
	if err := h.Provide(plugins.ServiceAccess, p); err != nil {
		return err
	}
	h.Use(func(ctx context.Context, b *bot.Bot, u *bot.Update, next func() error) error {
		uid := userOf(u)
		if uid == "" || p.permitted(uid) {
			return next()
		}
		p.reject(ctx, b, u)
		return nil
	})
	return nil
}

// userOf returns the acting user; for events without a sender the
// private chat ID doubles as the user ID.
func userOf(u *bot.Update) string {
	if id := u.SenderID(); id != "" {
		return id
	}
	if u.IsPrivateChat() {
		return u.ChatID
	}
	return ""
}

func (p *Plugin) permitted(uid string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if _, ok := p.admins[uid]; ok {
		return true
	}
	if _, bad := p.blocked[uid]; bad {
		return false
	}
	if len(p.allowed) == 0 {
		return true
	}
	_, ok := p.allowed[uid]
	return ok
}

// reject answers a denied user once (when a deny message is set).
func (p *Plugin) reject(ctx context.Context, b *bot.Bot, u *bot.Update) {
	p.mu.Lock()
	msg := p.deny
	_, done := p.notified[u.ChatID]
	if msg != "" && !done {
		p.notified[u.ChatID] = struct{}{}
	}
	p.mu.Unlock()
	if msg != "" && !done {
		_, _ = b.SendMessage(ctx, u.ChatID, msg)
	}
}

// IsAdmin implements plugins.AccessService.
func (p *Plugin) IsAdmin(id string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.admins[id]
	return ok
}

// IsBlocked implements plugins.AccessService.
func (p *Plugin) IsBlocked(id string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.blocked[id]
	return ok
}

// Block adds a user to the block-list at runtime.
func (p *Plugin) Block(id string) { p.mu.Lock(); p.blocked[id] = struct{}{}; p.mu.Unlock() }

// Unblock removes a user from the block-list.
func (p *Plugin) Unblock(id string) { p.mu.Lock(); delete(p.blocked, id); p.mu.Unlock() }

// AddAdmin grants administrator rights at runtime.
func (p *Plugin) AddAdmin(id string) { p.mu.Lock(); p.admins[id] = struct{}{}; p.mu.Unlock() }
