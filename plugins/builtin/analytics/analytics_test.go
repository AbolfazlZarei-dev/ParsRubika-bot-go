package analytics

import (
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func TestCountsUpdatesAndCommands(t *testing.T) {
	env := plugintest.New(t)
	env.Manager.MustRegister(New())
	env.Enable("analytics")
	_ = env.Say("b0u", "u", "/start")
	_ = env.Say("b0u", "u", "/start now")
	_ = env.Say("b0u", "u", "!help@bot")
	_ = env.Say("b0u", "u", "plain")
	_ = env.Event(models.UpdateTypeStartedBot, "b0u")
	tr, _ := plugins.Service[Tracker](env.Manager, plugins.ServiceStats)
	if tr.Counter("updates.total") != 5 || tr.Counter("command.start") != 2 || tr.Counter("command.help") != 1 || tr.Counter("updates.StartedBot") != 1 {
		c, _ := tr.Snapshot()
		t.Fatalf("counters = %v", c)
	}
	tr.SetGauge("queue", 3)
	if tr.Gauge("queue") != 3 {
		t.Fatal("gauge")
	}
}
