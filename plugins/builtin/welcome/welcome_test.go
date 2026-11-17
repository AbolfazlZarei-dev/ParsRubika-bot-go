package welcome

import (
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func TestGreetsOnStartedBotAndStartCommand(t *testing.T) {
	env := plugintest.New(t)
	env.Manager.MustRegister(New(), plugins.Config{"text": "hello"})
	env.Enable("welcome")
	_ = env.Event(models.UpdateTypeStartedBot, "b0u")
	_ = env.Say("b0u", "u", "/start")
	got := env.Sent()
	if len(got) != 2 || got[0] != "hello" || got[1] != "hello" {
		t.Fatalf("sent = %v", got)
	}
}
