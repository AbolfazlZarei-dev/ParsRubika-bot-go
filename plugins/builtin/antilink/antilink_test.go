package antilink

import (
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func TestContains(t *testing.T) {
	cases := []struct {
		text     string
		allowed  []string
		mentions bool
		want     bool
	}{
		{"hello there", nil, false, false},
		{"see https://spam.example.com/x", nil, false, true},
		{"join rubika.ir/mychannel now", nil, false, true},
		{"docs at www.good.org", []string{"good.org"}, false, false},
		{"ping @someone", nil, false, false},
		{"ping @someone", nil, true, true},
		{"سلام به همه", nil, false, false},
		{"version 1.2.3 released", nil, false, false},
	}
	for _, c := range cases {
		if got := Contains(c.text, c.allowed, c.mentions); got != c.want {
			t.Errorf("Contains(%q) = %v, want %v", c.text, got, c.want)
		}
	}
}

func TestDeletesLinksOfNonAdminsOnly(t *testing.T) {
	env := plugintest.New(t)
	env.Manager.MustRegister(New(), plugins.Config{"warn_message": "no links"})
	hits := 0
	env.Bot.Handle(nil, func(c *bot.Ctx) error { hits++; return nil })
	env.Enable("antilink")
	env.Respond("getChatMember", func(p map[string]any) map[string]any {
		if p["user_id"] == "boss" {
			return map[string]any{"type": "Admin"}
		}
		return map[string]any{"type": "Member"}
	})
	_ = env.Say("g0grp", "member", "visit https://spam.test")
	_ = env.Say("g0grp", "boss", "visit https://ok.test")
	_ = env.Say("g0grp", "member", "plain text")
	_ = env.Say("u0priv", "member", "https://private.chat/allowed")
	if len(env.Calls("deleteMessage")) != 1 {
		t.Fatalf("deleted %d messages, want 1", len(env.Calls("deleteMessage")))
	}
	if hits != 3 {
		t.Fatalf("handler hits = %d, want 3", hits)
	}
	if got := env.Sent(); len(got) != 1 || got[0] != "no links" {
		t.Fatalf("warnings = %v", got)
	}
}
