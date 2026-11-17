package bot

import (
	"context"
	"strconv"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/rules"
)

// BenchmarkRuleCheck measures the cost of running a typical rule chain
// against an update.
func BenchmarkRuleCheck(b *testing.B) {
	ctx := context.Background()
	upd := &models.Update{
		ChatID:     "g0GROUP",
		Type:       models.UpdateTypeNewMessage,
		NewMessage: &models.Message{Text: "/start now"},
	}
	rules := []rules.Rule{
		rules.NewCommand([]string{"start"}, []string{"/"}, false, false),
		rules.Group{},
		rules.NewText(".*", true),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range rules {
			_ = r.Check(ctx, upd)
		}
	}
}

// BenchmarkParseUpdate measures the cost of parsing a raw update
// payload into a typed *models.Update.
func BenchmarkParseUpdate(b *testing.B) {
	bot := New("dummy_token")
	raw := map[string]any{
		"type":    "NewMessage",
		"chat_id": "g0GROUP",
		"new_message": map[string]any{
			"message_id":  "m1",
			"text":        "hello world",
			"sender_id":   "u0USER",
			"sender_type": "User",
		},
		"update_time": "2024-01-01T00:00:00Z",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parseUpdate(raw, bot)
	}
}

// BenchmarkBuildRawUpdateKey measures the cost of the dedup key builder.
func BenchmarkBuildRawUpdateKey(b *testing.B) {
	raw := map[string]any{
		"type":        "NewMessage",
		"update_time": "2024",
		"new_message": map[string]any{"message_id": "m1"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buildRawUpdateKey(raw)
	}
}

func BenchmarkDispatchOneMiddleware(b *testing.B) {
	bt := New("t", WithPersistOffset(false), WithLogger(nullLogger{}), WithRateLimit(0))
	bt.Use(func(ctx context.Context, _ *Bot, _ *Update, next func() error) error { return next() })
	bt.Command("ping", func(c *Ctx) error { return nil })
	u := &models.Update{Type: models.UpdateTypeNewMessage, ChatID: "b0c", NewMessage: &models.Message{MessageID: "m", Text: "/ping x"}}
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.UpdateTime = strconv.Itoa(i) // unique key so dedup does not short-circuit
		_ = bt.Dispatch(ctx, u)
	}
}

func BenchmarkSeenSetAdd(b *testing.B) {
	s := newSeenSet(2048)
	keys := make([]string, 4096)
	for i := range keys {
		keys[i] = strconv.Itoa(i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.add(keys[i&4095])
	}
}
