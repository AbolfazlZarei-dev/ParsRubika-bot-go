package bot

import (
	"context"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/rules"
)

func TestRulesText(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{ChatID: "g0GROUP", Type: models.UpdateTypeNewMessage, NewMessage: &models.Message{Text: "hello"}}
	cases := []struct {
		name string
		r    rules.Rule
		want bool
	}{
		{"Text matches 'hello'", rules.NewText("hello", false), true},
		{"Text matches 'world'", rules.NewText("world", false), false},
		{"Text regex .*ell.*", rules.NewText(".*ell.*", true), true},
		{"Text regex ^xyz$", rules.NewText("^xyz$", true), false},
		{"Any text", rules.NewText("", false), true},
		{"Group", rules.Group{}, true},
		{"Private", rules.Private{}, false},
		{"Channel", rules.Channel{}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.r.Check(ctx, upd); got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestRuleCombinators(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{ChatID: "g0GROUP", Type: models.UpdateTypeNewMessage, NewMessage: &models.Message{Text: "hello"}}
	cases := []struct {
		name string
		r    rules.Rule
		want bool
	}{
		{"And(Text, Group)", rules.And{rules.NewText("hello", false), rules.Group{}}, true},
		{"And(Text, Private)", rules.And{rules.NewText("hello", false), rules.Private{}}, false},
		{"Or(Group, Private)", rules.Or{rules.Group{}, rules.Private{}}, true},
		{"Or(Channel, Private)", rules.Or{rules.Channel{}, rules.Private{}}, false},
		{"Not(Group)", rules.Not{Inner: rules.Group{}}, false},
		{"Not(Private)", rules.Not{Inner: rules.Private{}}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.r.Check(ctx, upd); got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestRulesCommand(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{
		ChatID:     "b0USER",
		Type:       models.UpdateTypeNewMessage,
		NewMessage: &models.Message{Text: "/start now"},
	}
	c := rules.NewCommand([]string{"start"}, []string{"/"}, false, false)
	if !c.Check(ctx, upd) {
		t.Errorf("expected /start to match")
	}
	upd.NewMessage.Text = "/help"
	if c.Check(ctx, upd) {
		t.Errorf("did not expect /help to match /start filter")
	}
}

func TestRulesChatID(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{ChatID: "b0USER"}
	c := rules.NewChatID("b0USER", "b0ADMIN")
	if !c.Check(ctx, upd) {
		t.Error("expected chat ID match")
	}
	upd.ChatID = "b0OTHER"
	if c.Check(ctx, upd) {
		t.Error("did not expect chat ID match")
	}
}

func TestRulesForward(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{
		Type:       models.UpdateTypeNewMessage,
		NewMessage: &models.Message{ForwardedFrom: &models.ForwardedFrom{TypeFrom: "User"}},
	}
	fwd := rules.Forward{}
	if !fwd.Check(ctx, upd) {
		t.Error("Forward rule should match forwarded messages")
	}
}

func TestRulesButton(t *testing.T) {
	ctx := context.Background()
	upd := &models.Update{
		Type:       models.UpdateTypeNewMessage,
		NewMessage: &models.Message{AuxData: &models.AuxData{ButtonID: "btn_123"}},
	}
	if !rules.NewButton("btn_123", false).Check(ctx, upd) {
		t.Error("Button rule should match btn_123")
	}
	if rules.NewButton("btn_456", false).Check(ctx, upd) {
		t.Error("Button rule should not match btn_456")
	}
	if !rules.NewButton("btn_.*", true).Check(ctx, upd) {
		t.Error("Button regex rule should match btn_123")
	}
}

func TestParseUpdate(t *testing.T) {
	b := New("dummy_token")
	raw := map[string]any{
		"type":    "NewMessage",
		"chat_id": "g0GROUP",
		"new_message": map[string]any{
			"message_id": "m1",
			"text":       "hello",
		},
		"update_time": "2024-01-01T00:00:00Z",
	}
	upd, err := parseUpdate(raw, b)
	if err != nil {
		t.Fatalf("parseUpdate: %v", err)
	}
	if upd.Type != models.UpdateTypeNewMessage {
		t.Errorf("Type = %q, want NewMessage", upd.Type)
	}
	if upd.ChatID != "g0GROUP" {
		t.Errorf("ChatID = %q", upd.ChatID)
	}
	if upd.NewMessage == nil || upd.NewMessage.Text != "hello" {
		t.Fatalf("NewMessage = %+v", upd.NewMessage)
	}
	if upd.MessageID() != "m1" {
		t.Errorf("MessageID = %q, want m1", upd.MessageID())
	}
	if upd.Text() != "hello" {
		t.Errorf("Text = %q, want hello", upd.Text())
	}
}

func TestBuildRawUpdateKey(t *testing.T) {
	cases := []struct {
		name string
		in   map[string]any
		want string
	}{
		{"empty", map[string]any{}, ""},
		{"removed_message", map[string]any{"type": "RemovedMessage", "removed_message_id": "m1"}, "RemovedMessage:m1"},
		{
			"new_message",
			map[string]any{
				"type":        "NewMessage",
				"update_time": "2024",
				"new_message": map[string]any{"message_id": "m1"},
			},
			"NewMessage:m1:2024",
		},
		{
			"updated_message",
			map[string]any{
				"type":            "UpdatedMessage",
				"updated_message": map[string]any{"message_id": "m2", "time": "T1"},
			},
			"UpdatedMessage:m2:T1",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := buildRawUpdateKey(c.in); got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestAPIError(t *testing.T) {
	e := NewAPIError("InvalidAuth", "bad key")
	if e.Error() != "bot: InvalidAuth: bad key" {
		t.Errorf("Error() = %q", e.Error())
	}
	if !IsStatus(e, "InvalidAuth") {
		t.Error("IsStatus should return true")
	}
	if IsStatus(e, "Other") {
		t.Error("IsStatus should return false for non-matching status")
	}
	if IsStatus(nil, "InvalidAuth") {
		t.Error("IsStatus should return false for nil error")
	}
}
