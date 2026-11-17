package users

import (
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest"
)

func TestTracksChatsAndSkipsStoppedOnes(t *testing.T) {
	env := plugintest.New(t)
	env.Manager.MustRegister(New())
	env.Enable("users")
	_ = env.Say("b0a", "a", "hi")
	_ = env.Say("b0a", "a", "again")
	_ = env.Say("b0b", "b", "yo")
	_ = env.Event(models.UpdateTypeStoppedBot, "b0b")
	dir, _ := plugins.Service[plugins.UserDirectory](env.Manager, plugins.ServiceUsers)
	if got := dir.Chats(); len(got) != 1 || got[0] != "b0a" {
		t.Fatalf("chats = %v", got)
	}
	_ = env.Event(models.UpdateTypeStartedBot, "b0b")
	if dir.Count() != 2 {
		t.Fatalf("restart must reactivate the chat, count=%d", dir.Count())
	}
	p, _ := env.Manager.Get("users")
	if r, _ := p.(*Plugin).Get("b0a"); r.Messages != 2 {
		t.Fatalf("messages = %d", r.Messages)
	}
}
