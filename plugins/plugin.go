// Package plugins is the plugin framework of ParsRubika.
//
// A plugin is a self-contained feature (anti-spam, a /help command, a
// welcome message, ...) that plugs into a bot through a Host. The Host
// records everything the plugin registers — handlers, middlewares,
// hooks, timers, services — so the Manager can switch the plugin off
// again and leave the bot exactly as it was.
//
//	b := bot.New(token)
//	m := plugins.New(b)
//	m.MustRegister(antispam.New(), plugins.Config{"max_messages": 5})
//	m.MustRegister(help.New())
//	m.Install()            // enable everything when the bot starts
//	_ = b.Start(ctx)
//
// Writing a plugin only needs two methods:
//
//	type Hello struct{ plugins.Base }
//
//	func New() *Hello { return &Hello{plugins.NewBase(plugins.Meta{Name: "hello"})} }
//
//	func (p *Hello) Setup(ctx context.Context, h *plugins.Host) error {
//		h.Command("hello", func(c *bot.Ctx) error { _, err := c.Reply("Hi!"); return err })
//		return nil
//	}
package plugins

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// DefaultPriority is the start-up priority of plugins that do not set one.
// Plugins with a lower number are enabled earlier.
const DefaultPriority = 100

// Meta describes a plugin.
type Meta struct {
	// Name is the unique, lower-case identifier of the plugin.
	Name string
	// Version is a free-form version string (semver recommended).
	Version string
	// Description is a one-line summary shown by List and /plugins.
	Description string
	// Author credits the plugin author.
	Author string
	// Requires lists plugins that must be enabled before this one. The
	// Manager enables them automatically and disables this plugin first
	// when one of them is switched off.
	Requires []string
	// Priority orders EnableAll: lower values start first
	// (DefaultPriority when zero).
	Priority int
	// Disabled keeps the plugin registered but off until Enable is
	// called explicitly; EnableAll skips it.
	Disabled bool
}

// Plugin is the contract every plugin fulfils.
type Plugin interface {
	// Meta returns the static description of the plugin.
	Meta() Meta
	// Setup wires the plugin into the bot through h. Anything registered
	// through h is removed automatically when the plugin is disabled.
	// Returning an error aborts the activation and rolls it back.
	Setup(ctx context.Context, h *Host) error
}

// Teardowner is implemented by plugins that need to release resources of
// their own (files, connections) when they are disabled.
type Teardowner interface {
	Teardown(ctx context.Context) error
}

// Defaulter lets a plugin declare its default configuration. User
// supplied values are merged on top of it.
type Defaulter interface {
	DefaultConfig() Config
}

// Validator lets a plugin reject an invalid configuration before Setup.
type Validator interface {
	Validate(cfg Config) error
}

// Base is a convenience type to embed in plugins: it implements Meta.
type Base struct{ meta Meta }

// NewBase builds a Base from m.
func NewBase(m Meta) Base { return Base{meta: m} }

// Meta implements Plugin.
func (b Base) Meta() Meta { return b.meta }

// Errors returned by the Manager.
var (
	ErrInvalidName   = errors.New("plugins: invalid plugin name")
	ErrDuplicate     = errors.New("plugins: plugin already registered")
	ErrNotFound      = errors.New("plugins: plugin not registered")
	ErrDependency    = errors.New("plugins: dependency problem")
	ErrNoBot         = errors.New("plugins: manager has no bot")
	ErrAlreadyActive = errors.New("plugins: plugin already enabled")
)

// normalize canonicalises a plugin name.
func normalize(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

func validName(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-', r == '_', r == '.':
		default:
			return false
		}
	}
	return true
}

// SetupError wraps a failure raised while enabling a plugin.
type SetupError struct {
	Plugin string
	Err    error
}

// Error implements error.
func (e *SetupError) Error() string { return fmt.Sprintf("plugins: enable %q: %v", e.Plugin, e.Err) }

// Unwrap returns the underlying error.
func (e *SetupError) Unwrap() error { return e.Err }
