package builtin_test

import (
	"context"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func TestRegisterAllAndEnableEverything(t *testing.T) {
	env := plugintest.New(t)
	if err := builtin.RegisterAll(env.Manager, map[string]plugins.Config{
		"access":    {"admins": []string{"root"}},
		"forcejoin": {"channels": []string{"c0x"}},
	}); err != nil {
		t.Fatal(err)
	}
	env.Respond("getChatMember", func(map[string]any) map[string]any { return map[string]any{"type": "Member"} })
	if err := env.Manager.EnableAll(context.Background()); err != nil {
		t.Fatalf("EnableAll: %v", err)
	}
	if got := len(env.Manager.Enabled()); got != len(builtin.All()) {
		t.Fatalf("enabled %d of %d: %v", got, len(builtin.All()), env.Manager.Enabled())
	}
	// A normal conversation must flow through all middlewares.
	if err := env.Say("b0u", "u", "/help"); err != nil {
		t.Fatal(err)
	}
	if len(env.Sent()) != 1 {
		t.Fatalf("sent = %v", env.Sent())
	}
	if err := env.Manager.DisableAll(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(env.Manager.Enabled()) != 0 {
		t.Fatal("DisableAll left plugins active")
	}
}

func TestForceJoinSkippedWithoutConfig(t *testing.T) {
	env := plugintest.New(t)
	_ = builtin.RegisterAll(env.Manager, nil)
	if _, ok := env.Manager.Get("forcejoin"); ok {
		t.Fatal("forcejoin must not be registered without channels")
	}
	if err := env.Manager.EnableAll(context.Background()); err != nil {
		t.Fatalf("EnableAll: %v", err)
	}
}
