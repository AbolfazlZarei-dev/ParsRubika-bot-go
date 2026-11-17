package stats

import (
	"strings"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/access"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func TestStatsIsAdminOnly(t *testing.T) {
	env := plugintest.New(t)
	env.Manager.MustRegister(access.New(), plugins.Config{"admins": []string{"root"}})
	env.Manager.MustRegister(New())
	env.Enable("access", "stats")
	_ = env.Say("b0u", "user", "/stats")
	if len(env.Sent()) != 0 {
		t.Fatal("non-admin got statistics")
	}
	_ = env.Say("b0r", "root", "/stats")
	got := env.Sent()
	if len(got) != 1 || !strings.Contains(got[0], "goroutine") {
		t.Fatalf("sent = %v", got)
	}
}
