package forms

import (
	"context"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

func TestFormBasic(t *testing.T) {
	f := New()
	var called []string
	f.Add("start", func(ctx context.Context, b *bot.Bot, upd *models.Update, answer string) (string, error) {
		called = append(called, "start:"+answer)
		return "step2", nil
	})
	f.Add("step2", func(ctx context.Context, b *bot.Bot, upd *models.Update, answer string) (string, error) {
		called = append(called, "step2:"+answer)
		return "", nil // end
	})

	if f.IsActive("chat1") {
		t.Error("chat1 should not be active before Start")
	}
	f.Start("chat1")
	if !f.IsActive("chat1") {
		t.Error("chat1 should be active after Start")
	}

	upd := &models.Update{ChatID: "chat1", NewMessage: &models.Message{Text: "hello"}}
	if err := f.Handle(context.Background(), nil, upd); err != nil {
		t.Fatalf("Handle (start): %v", err)
	}
	if !f.IsActive("chat1") {
		t.Error("chat1 should still be active after step1")
	}
	upd.NewMessage.Text = "world"
	if err := f.Handle(context.Background(), nil, upd); err != nil {
		t.Fatalf("Handle (step2): %v", err)
	}
	if f.IsActive("chat1") {
		t.Error("chat1 should not be active after last step")
	}
	if len(called) != 2 {
		t.Fatalf("called = %v, want 2 entries", called)
	}
	if called[0] != "start:hello" || called[1] != "step2:world" {
		t.Errorf("called = %v", called)
	}
}

func TestFormStop(t *testing.T) {
	f := New()
	f.Start("c1")
	if !f.IsActive("c1") {
		t.Error("c1 should be active")
	}
	f.Stop("c1")
	if f.IsActive("c1") {
		t.Error("c1 should not be active after Stop")
	}
}

func TestFormHandleInactive(t *testing.T) {
	f := New()
	upd := &models.Update{ChatID: "c1", NewMessage: &models.Message{Text: "x"}}
	if err := f.Handle(context.Background(), nil, upd); err != nil {
		t.Fatalf("Handle on inactive form should return nil, got %v", err)
	}
}
