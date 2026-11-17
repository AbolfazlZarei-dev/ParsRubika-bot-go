package rules

import (
	"context"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

func TestAllRulesConstructors(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{
		ChatID:     "g0GROUP",
		Type:       models.UpdateTypeNewMessage,
		NewMessage: &models.Message{Text: "hello"},
	}
	// All the simple struct rules.
	simple := []Rule{
		Private{},
		Group{},
		Channel{},
		BotSender{},
		Forward{},
		File{},
		Sticker{},
		Location{},
		Poll{},
		Contact{},
	}
	for _, r := range simple {
		// Just ensure Check runs without panic.
		_ = r.Check(ctx, upd)
	}
}

func TestNewUpdateType(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{Type: models.UpdateTypeNewMessage}
	r := NewUpdateType(models.UpdateTypeNewMessage)
	if !r.Check(ctx, upd) {
		t.Error("UpdateType should match NewMessage")
	}
	r2 := NewUpdateType(models.UpdateTypeRemovedMessage)
	if r2.Check(ctx, upd) {
		t.Error("UpdateType should not match RemovedMessage")
	}
}

func TestNewReply(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{
		NewMessage: &models.Message{ReplyToMessageID: "m1"},
	}
	r := NewReply("m1")
	if !r.Check(ctx, upd) {
		t.Error("Reply should match m1")
	}
	r2 := NewReply("m2")
	if r2.Check(ctx, upd) {
		t.Error("Reply should not match m2")
	}
}

func TestCommandDoesNotMatchLongerWords(t *testing.T) {
	mk := func(text string) *models.Update {
		return &models.Update{Type: models.UpdateTypeNewMessage, NewMessage: &models.Message{Text: text}}
	}
	r := NewCommand([]string{"warn"}, []string{"/", "!"}, true, false)
	cases := map[string]bool{
		"/warn": true, "/warn bob": true, "!WARN bob": true, "/warn@my_bot x": true,
		"/warnings": false, "/warning": false, "/warnx": false, "warn": false, "/war": false,
	}
	for text, want := range cases {
		if got := r.Check(context.Background(), mk(text)); got != want {
			t.Errorf("Check(%q) = %v, want %v", text, got, want)
		}
	}
}
