package rules

import (
	"context"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

func TestPhotoRule(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{
		Type:       models.UpdateTypeNewMessage,
		NewMessage: &models.Message{File: &models.File{FileName: "test.jpg"}},
	}
	p := Photo{}
	if !p.Check(ctx, upd) {
		t.Error("Photo should match .jpg")
	}
	upd.NewMessage.File.FileName = "test.txt"
	if p.Check(ctx, upd) {
		t.Error("Photo should not match .txt")
	}
}

func TestVoiceRule(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{
		Type:       models.UpdateTypeNewMessage,
		NewMessage: &models.Message{File: &models.File{FileName: "test.ogg"}},
	}
	v := Voice{}
	if !v.Check(ctx, upd) {
		t.Error("Voice should match .ogg")
	}
}

func TestVideoRule(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{
		Type:       models.UpdateTypeNewMessage,
		NewMessage: &models.Message{File: &models.File{FileName: "test.mp4"}},
	}
	v := Video{}
	if !v.Check(ctx, upd) {
		t.Error("Video should match .mp4")
	}
}

func TestGifRule(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{
		Type:       models.UpdateTypeNewMessage,
		NewMessage: &models.Message{File: &models.File{FileName: "test.gif"}},
	}
	g := Gif{}
	if !g.Check(ctx, upd) {
		t.Error("Gif should match .gif")
	}
}

func TestEditedRule(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{Type: models.UpdateTypeUpdatedMessage}
	e := Edited{}
	if !e.Check(ctx, upd) {
		t.Error("Edited should match UpdatedMessage")
	}
	upd.Type = models.UpdateTypeNewMessage
	if e.Check(ctx, upd) {
		t.Error("Edited should not match NewMessage")
	}
}

func TestSenderTypeRule(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{
		Type:       models.UpdateTypeNewMessage,
		NewMessage: &models.Message{SenderType: "User"},
	}
	r := NewSenderType("User")
	if !r.Check(ctx, upd) {
		t.Error("SenderType User should match")
	}
	r2 := NewSenderType("Bot")
	if r2.Check(ctx, upd) {
		t.Error("SenderType Bot should not match User")
	}
}

func TestHasAuxDataRule(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{
		Type:       models.UpdateTypeNewMessage,
		NewMessage: &models.Message{AuxData: &models.AuxData{ButtonID: "btn_1"}},
	}
	h := HasAuxData{}
	if !h.Check(ctx, upd) {
		t.Error("HasAuxData should match")
	}
	upd.NewMessage.AuxData = nil
	if h.Check(ctx, upd) {
		t.Error("HasAuxData should not match")
	}
}

func TestStartedBotRule(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{Type: models.UpdateTypeStartedBot}
	s := StartedBot{}
	if !s.Check(ctx, upd) {
		t.Error("StartedBot should match")
	}
}

func TestStoppedBotRule(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{Type: models.UpdateTypeStoppedBot}
	s := StoppedBot{}
	if !s.Check(ctx, upd) {
		t.Error("StoppedBot should match")
	}
}
