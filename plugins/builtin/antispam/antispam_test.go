package antispam

import (
	"testing"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func rig(t *testing.T, cfg plugins.Config) (*plugintest.Env, *Plugin, *int, *time.Time) {
	env := plugintest.New(t)
	p := New()
	clock := time.Unix(5000, 0)
	p.now = func() time.Time { return clock }
	env.Manager.MustRegister(p, cfg)
	hits := new(int)
	env.Bot.Handle(nil, func(c *bot.Ctx) error { *hits++; return nil })
	env.Enable("antispam")
	return env, p, hits, &clock
}

func TestFloodMutesAndWarnsOnce(t *testing.T) {
	env, p, hits, clock := rig(t, plugins.Config{"max_messages": 3, "window": "10s", "repeat_limit": 0, "mute_for": "1m", "warn_message": "calm"})
	for i := 0; i < 8; i++ {
		_ = env.Say("b0u", "u", string(rune('a'+i)))
	}
	if *hits != 3 || !p.IsMuted("u") {
		t.Fatalf("hits=%d muted=%v", *hits, p.IsMuted("u"))
	}
	if got := env.Sent(); len(got) != 1 || got[0] != "calm" {
		t.Fatalf("warn must be sent exactly once: %v", got)
	}
	*clock = clock.Add(2 * time.Minute)
	_ = env.Say("b0u", "u", "later")
	if *hits != 4 {
		t.Fatalf("mute did not expire, hits=%d", *hits)
	}
}

func TestRepeatedTextDetection(t *testing.T) {
	env, _, hits, _ := rig(t, plugins.Config{"max_messages": 100, "repeat_limit": 2, "warn_message": ""})
	for i := 0; i < 5; i++ {
		_ = env.Say("b0u", "u", "same")
	}
	if *hits != 2 {
		t.Fatalf("hits = %d, want 2", *hits)
	}
}

func TestUnmute(t *testing.T) {
	env, p, hits, _ := rig(t, plugins.Config{"max_messages": 1, "repeat_limit": 0, "warn_message": ""})
	_ = env.Say("b0u", "u", "1")
	_ = env.Say("b0u", "u", "2")
	p.Unmute("u")
	_ = env.Say("b0u", "u", "3")
	if *hits != 2 {
		t.Fatalf("hits = %d", *hits)
	}
}
