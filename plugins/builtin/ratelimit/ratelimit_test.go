package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/access"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func TestBurstThenThrottleThenRefill(t *testing.T) {
	env := plugintest.New(t)
	p := New()
	clock := time.Unix(1000, 0)
	p.now = func() time.Time { return clock }
	env.Manager.MustRegister(p, plugins.Config{"rate": 1.0, "burst": 3, "message": "slow down"})
	hits := 0
	env.Bot.Handle(nil, func(c *bot.Ctx) error { hits++; return nil })
	env.Enable("ratelimit")

	for i := 0; i < 6; i++ {
		_ = env.Say("b0u", "u", "x")
	}
	if hits != 3 {
		t.Fatalf("burst: hits = %d, want 3", hits)
	}
	if got := env.Sent(); len(got) != 1 {
		t.Fatalf("throttle notice must be rate-limited too, got %v", got)
	}
	clock = clock.Add(2 * time.Second) // two tokens come back
	for i := 0; i < 5; i++ {
		_ = env.Say("b0u", "u", "y")
	}
	if hits != 5 {
		t.Fatalf("after refill: hits = %d, want 5", hits)
	}
}

func TestAdminsAreExempt(t *testing.T) {
	env := plugintest.New(t)
	env.Manager.MustRegister(access.New(), plugins.Config{"admins": []string{"root"}})
	env.Manager.MustRegister(New(), plugins.Config{"burst": 1, "rate": 0.001})
	hits := 0
	env.Bot.Handle(nil, func(c *bot.Ctx) error { hits++; return nil })
	env.Enable("access", "ratelimit")
	for i := 0; i < 5; i++ {
		_ = env.Say("b0r", "root", "x")
	}
	if hits != 5 {
		t.Fatalf("admin was throttled: %d", hits)
	}
}

func TestInvalidConfigRejected(t *testing.T) {
	env := plugintest.New(t)
	env.Manager.MustRegister(New(), plugins.Config{"rate": -1.0})
	if err := env.Manager.Enable(context.Background(), "ratelimit"); err == nil { //nolint:staticcheck
		t.Fatal("negative rate must be rejected")
	}
}
