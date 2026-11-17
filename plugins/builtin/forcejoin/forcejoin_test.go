package forcejoin

import (
	"context"
	"errors"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func rig(t *testing.T) (*plugintest.Env, *int) {
	env := plugintest.New(t)
	env.Manager.MustRegister(New(), plugins.Config{"channels": []string{"c0chan"}, "links": []string{"rubika.ir/chan"}})
	hits := new(int)
	env.Bot.Handle(nil, func(c *bot.Ctx) error { *hits++; return nil })
	env.Enable("forcejoin")
	return env, hits
}

func TestNotJoinedIsBlockedWithLinks(t *testing.T) {
	env, hits := rig(t)
	env.Respond("getChatMember", func(map[string]any) map[string]any { return map[string]any{} })
	_ = env.Say("u0u", "u0u", "hi")
	if *hits != 0 {
		t.Fatal("non-member reached handlers")
	}
	if got := env.Sent(); len(got) != 1 || !contains(got[0], "rubika.ir/chan") {
		t.Fatalf("reminder = %v", got)
	}
}

func TestJoinedPassesAndIsCached(t *testing.T) {
	env, hits := rig(t)
	env.Respond("getChatMember", func(map[string]any) map[string]any { return map[string]any{"type": "Member"} })
	for i := 0; i < 3; i++ {
		_ = env.Say("u0u", "u0u", "hi")
	}
	if *hits != 3 {
		t.Fatalf("hits = %d", *hits)
	}
	if n := len(env.Calls("getChatMember")); n != 1 {
		t.Fatalf("membership must be cached, API called %d times", n)
	}
}

func TestAPIFailureFailsOpen(t *testing.T) {
	env, hits := rig(t)
	env.Server.Config.Handler = failing{}
	_ = env.Say("u0u", "u0u", "hi")
	if *hits != 1 {
		t.Fatal("an outage must not lock users out")
	}
}

func TestConfigWithoutChannelsIsRejected(t *testing.T) {
	env := plugintest.New(t)
	env.Manager.MustRegister(New())
	if err := env.Manager.Enable(context.Background(), "forcejoin"); err == nil || !errors.Is(err, errNoChannels) {
		t.Fatalf("err = %v", err)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
