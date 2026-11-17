package access

import (
	"context"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func setup(t *testing.T, cfg plugins.Config) (*plugintest.Env, *int) {
	env := plugintest.New(t)
	env.Manager.MustRegister(New(), cfg)
	hits := new(int)
	env.Bot.Handle(nil, func(c *bot.Ctx) error { *hits++; return nil })
	env.Enable("access")
	return env, hits
}

func TestBlockedUserNeverReachesHandlers(t *testing.T) {
	env, hits := setup(t, plugins.Config{"blocked": []string{"bad"}})
	_ = env.Say("u0bad", "bad", "hi")
	_ = env.Say("u0ok", "ok", "hi")
	if *hits != 1 {
		t.Fatalf("handler ran %d times, want 1", *hits)
	}
}

func TestAllowListAndAdminBypass(t *testing.T) {
	env, hits := setup(t, plugins.Config{"allowed": []string{"vip"}, "admins": []string{"root"}, "deny_message": "no"})
	_ = env.Say("u0x", "x", "1")
	_ = env.Say("u0x", "x", "2")
	_ = env.Say("u0vip", "vip", "3")
	_ = env.Say("u0root", "root", "4")
	if *hits != 2 {
		t.Fatalf("hits = %d", *hits)
	}
	if got := env.Sent(); len(got) != 1 || got[0] != "no" {
		t.Fatalf("deny message must be sent exactly once, got %v", got)
	}
	if !plugins.IsAdmin(env.Manager, "root") || plugins.IsAdmin(env.Manager, "vip") {
		t.Fatal("IsAdmin service wrong")
	}
}

func TestStartedBotUsesChatAsUser(t *testing.T) {
	env, hits := setup(t, plugins.Config{"blocked": []string{"u0bad"}})
	_ = env.Event(models.UpdateTypeStartedBot, "u0bad")
	_ = env.Event(models.UpdateTypeStartedBot, "u0good")
	if *hits != 1 {
		t.Fatalf("hits = %d", *hits)
	}
	_ = context.Background()
}
