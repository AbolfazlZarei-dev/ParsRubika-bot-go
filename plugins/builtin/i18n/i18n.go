// Package i18n is a plugin that publishes a translator ("i18n" service)
// and remembers the language of every user, switchable with /lang.
//
// Configuration:
//
//	default   string                    fallback language (default "fa")
//	messages  map[lang]map[key]string   translations
//	command   string                    command that changes the language (default "lang")
package i18n

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

// Translator is the published service.
type Translator interface {
	// T translates key for the user (their language, else the default,
	// else the key itself) and applies fmt.Sprintf when args are given.
	T(userID, key string, args ...any) string
	// Lang returns the language of a user.
	Lang(userID string) string
	// SetLang stores the language of a user.
	SetLang(userID, lang string) bool
	// Languages lists the known languages.
	Languages() []string
}

// Plugin implements plugins.Plugin and Translator.
type Plugin struct {
	plugins.Base
	mu       sync.RWMutex
	def      string
	messages map[string]map[string]string
	users    map[string]string
}

// New creates the plugin.
func New() *Plugin {
	return &Plugin{Base: plugins.NewBase(plugins.Meta{
		Name: "i18n", Version: "1.0.0", Priority: 15,
		Description: "Per-user language and translations",
	})}
}

// DefaultConfig implements plugins.Defaulter.
func (p *Plugin) DefaultConfig() plugins.Config {
	return plugins.Config{"default": "fa", "command": "lang"}
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	cfg := h.Config()
	p.mu.Lock()
	p.def = cfg.String("default", "fa")
	p.messages, p.users = map[string]map[string]string{}, map[string]string{}
	switch raw := cfg["messages"].(type) {
	case map[string]any:
		for lang, m := range raw {
			if mm, ok := m.(map[string]any); ok {
				p.messages[lang] = map[string]string{}
				for k, v := range mm {
					p.messages[lang][k] = fmt.Sprint(v)
				}
			}
		}
	case map[string]map[string]string:
		for lang, m := range raw {
			p.messages[lang] = map[string]string{}
			for k, v := range m {
				p.messages[lang][k] = v
			}
		}
	}
	p.mu.Unlock()
	if cmd := cfg.String("command", "lang"); cmd != "" {
		h.Help("/"+cmd+" <کد>", "تغییر زبان")
		h.Command(cmd, func(c *bot.Ctx) error {
			arg := strings.ToLower(strings.TrimSpace(c.ArgsText()))
			if arg == "" || !p.SetLang(c.SenderID(), arg) {
				_, err := c.Reply("زبان‌های موجود: " + strings.Join(p.Languages(), ", "))
				return err
			}
			_, err := c.Reply(p.T(c.SenderID(), "lang.changed"))
			return err
		})
	}
	return h.Provide(plugins.ServiceI18n, p)
}

// AddLanguage merges translations for lang (usable from code).
func (p *Plugin) AddLanguage(lang string, msgs map[string]string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.messages == nil {
		p.messages = map[string]map[string]string{}
	}
	if p.messages[lang] == nil {
		p.messages[lang] = map[string]string{}
	}
	for k, v := range msgs {
		p.messages[lang][k] = v
	}
}

// T implements Translator.
func (p *Plugin) T(userID, key string, args ...any) string {
	p.mu.RLock()
	lang := p.users[userID]
	if lang == "" {
		lang = p.def
	}
	s, ok := p.messages[lang][key]
	if !ok {
		s, ok = p.messages[p.def][key]
	}
	p.mu.RUnlock()
	if !ok {
		return key
	}
	if len(args) > 0 {
		return fmt.Sprintf(s, args...)
	}
	return s
}

// Lang implements Translator.
func (p *Plugin) Lang(userID string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if l := p.users[userID]; l != "" {
		return l
	}
	return p.def
}

// SetLang implements Translator; it refuses unknown languages.
func (p *Plugin) SetLang(userID, lang string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.messages[lang]; !ok {
		return false
	}
	p.users[userID] = lang
	return true
}

// Languages implements Translator.
func (p *Plugin) Languages() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]string, 0, len(p.messages))
	for l := range p.messages {
		out = append(out, l)
	}
	sort.Strings(out)
	return out
}
