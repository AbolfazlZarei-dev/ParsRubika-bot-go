// Package router provides a simple request router for the bot API.
// It dispatches incoming updates to registered handlers based on
// the update type and optional pre-conditions.
//
// The router is an alternative to manually calling `bot.On(...)` for
// each handler — useful when the bot has dozens of handlers and you
// want to keep them in a single declarative list.
package router

import (
	"context"
	"strings"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/rules"
)

// Handler is the signature of a routed handler.
type Handler func(ctx context.Context, b *bot.Bot, upd *models.Update) error

// Route is a single routing entry: a set of rules and a handler.
type Route struct {
	Rules   []rules.Rule
	Handler Handler
}

// Router holds a list of routes and dispatches updates to the first
// matching one.
type Router struct {
	routes []Route
}

// New constructs an empty Router.
func New() *Router { return &Router{} }

// Add appends a route to the router.
func (r *Router) Add(rs []rules.Rule, h Handler) *Router {
	r.routes = append(r.routes, Route{Rules: rs, Handler: h})
	return r
}

// Command adds a route that matches a specific command.
func (r *Router) Command(cmds []string, h Handler) *Router {
	return r.Add([]rules.Rule{rules.NewCommand(cmds, []string{"/"}, false, false)}, h)
}

// Text adds a route that matches a text message.
func (r *Router) Text(text string, h Handler) *Router {
	return r.Add([]rules.Rule{rules.NewText(text, false)}, h)
}

// Regex adds a route that matches a regex against the text.
func (r *Router) Regex(pattern string, h Handler) *Router {
	return r.Add([]rules.Rule{rules.NewText(pattern, true)}, h)
}

// AnyText adds a route that matches any text message.
func (r *Router) AnyText(h Handler) *Router {
	return r.Add([]rules.Rule{rules.NewText("", false)}, h)
}

// Button adds a route that matches a button click by button_id.
func (r *Router) Button(buttonID string, h Handler) *Router {
	return r.Add([]rules.Rule{rules.NewButton(buttonID, false)}, h)
}

// Group adds a route that matches updates from groups.
func (r *Router) Group(h Handler) *Router {
	return r.Add([]rules.Rule{rules.Group{}}, h)
}

// Private adds a route that matches updates from private chats.
func (r *Router) Private(h Handler) *Router {
	return r.Add([]rules.Rule{rules.Private{}}, h)
}

// Channel adds a route that matches updates from channels.
func (r *Router) Channel(h Handler) *Router {
	return r.Add([]rules.Rule{rules.Channel{}}, h)
}

// Attach registers the router on a bot. The first matching route's
// handler is invoked; subsequent routes are skipped (first-match
// semantics).
func (r *Router) Attach(b *bot.Bot) {
	for _, route := range r.routes {
		rs := route.Rules
		h := route.Handler
		b.On(rs, h)
	}
}

// PrintRoutes returns a human-readable list of all registered
// routes (for debugging).
func (r *Router) PrintRoutes() string {
	var sb strings.Builder
	for i, route := range r.routes {
		sb.WriteString("Route ")
		sb.WriteByte(byte('0' + i))
		sb.WriteString(": ")
		for _, r := range route.Rules {
			sb.WriteString(getTypeName(r))
			sb.WriteString(" ")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// getTypeName returns a short name for a rule.
func getTypeName(r rules.Rule) string {
	if r == nil {
		return "nil"
	}
	// The Rule interface doesn't expose a name, so we use a type
	// switch for the common built-in rules.
	switch r.(type) {
	case *rules.Text:
		return "Text"
	case *rules.Command:
		return "Command"
	case *rules.ChatID:
		return "ChatID"
	case *rules.UpdateType:
		return "UpdateType"
	case *rules.Button:
		return "Button"
	case rules.Group:
		return "Group"
	case rules.Private:
		return "Private"
	case rules.Channel:
		return "Channel"
	case rules.Forward:
		return "Forward"
	case rules.File:
		return "File"
	case rules.Photo:
		return "Photo"
	case rules.Voice:
		return "Voice"
	case rules.Video:
		return "Video"
	case rules.Audio:
		return "Audio"
	case rules.Gif:
		return "Gif"
	case rules.Edited:
		return "Edited"
	case rules.StartedBot:
		return "StartedBot"
	case rules.StoppedBot:
		return "StoppedBot"
	case rules.HasAuxData:
		return "HasAuxData"
	}
	return "Rule"
}
