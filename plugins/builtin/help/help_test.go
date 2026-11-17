package help

import (
	"strings"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func TestHelpListsDocumentedCommandsAndBeatsFallback(t *testing.T) {
	env := plugintest.New(t)
	env.Manager.MustRegister(New(), plugins.Config{"header": "H"})
	// An application-level catch-all registered BEFORE the plugin starts
	// must not swallow /help.
	env.Bot.Fallback(func(c *bot.Ctx) error { _, err := c.Reply("fallback"); return err })
	env.Enable("help")
	if err := env.Say("b0u", "u", "/help"); err != nil {
		t.Fatal(err)
	}
	got := env.Sent()
	if len(got) != 1 || !strings.HasPrefix(got[0], "H") || !strings.Contains(got[0], "/help") {
		t.Fatalf("reply = %v", got)
	}
}
