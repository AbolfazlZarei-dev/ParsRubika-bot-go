// Package help is a plugin that answers /help with the commands that
// the other plugins documented through Host.Help.
//
// Configuration:
//
//	commands  []string  command names that trigger the help (default ["help"])
//	header    string    text above the list
//	footer    string    text below the list
package help

import (
	"context"
	"strings"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

// Plugin implements plugins.Plugin.
type Plugin struct{ plugins.Base }

// New creates the plugin.
func New() *Plugin {
	return &Plugin{plugins.NewBase(plugins.Meta{
		Name: "help", Version: "1.0.0", Priority: 90,
		Description: "Automatic /help generated from plugin documentation",
	})}
}

// DefaultConfig implements plugins.Defaulter.
func (p *Plugin) DefaultConfig() plugins.Config {
	return plugins.Config{"commands": []string{"help"}, "header": "📖 راهنمای ربات", "footer": ""}
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	cfg := h.Config()
	names := cfg.Strings("commands")
	if len(names) == 0 {
		names = []string{"help"}
	}
	h.Help("/"+names[0], "نمایش همین راهنما")
	h.Command(names[0], func(c *bot.Ctx) error {
		_, err := c.Reply(Render(h.Manager(), cfg.String("header", ""), cfg.String("footer", "")))
		return err
	}, names[1:]...)
	return nil
}

// Render builds the help text from the documented commands of all
// active plugins.
func Render(m *plugins.Manager, header, footer string) string {
	var sb strings.Builder
	if header != "" {
		sb.WriteString(header)
		sb.WriteString("\n\n")
	}
	entries := m.HelpEntries()
	if len(entries) == 0 {
		sb.WriteString("دستوری ثبت نشده است.")
	}
	for _, e := range entries {
		sb.WriteString(e.Command)
		if e.Description != "" {
			sb.WriteString(" — ")
			sb.WriteString(e.Description)
		}
		sb.WriteByte('\n')
	}
	if footer != "" {
		sb.WriteString("\n")
		sb.WriteString(footer)
	}
	return strings.TrimRight(sb.String(), "\n")
}
