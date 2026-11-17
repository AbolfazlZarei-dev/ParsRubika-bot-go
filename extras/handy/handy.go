// Package handy is a collection of small, ergonomic helpers that
// make writing ParsRubika bots easier. The helpers here do not fit
// neatly into a single concern — they are intentionally small
// one-liners that reduce boilerplate.
package handy

import (
	"context"
	"fmt"
	"strings"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// Echo registers a single handler that replies with the same text it
// received. Useful for testing the bot.
func Echo(b *bot.Bot) {
	b.On([]bot.Rule{bot.NewText("", false)}, func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
		_, err := bot.Reply(ctx, b, upd, upd.Text())
		return err
	})
}

// PingReply registers a single handler that replies "pong!" to any
// message that equals "ping" (case-insensitive).
func PingReply(b *bot.Bot) {
	b.On([]bot.Rule{bot.NewText("ping", false)}, func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
		_, err := bot.Reply(ctx, b, upd, "pong! 🏓")
		return err
	})
	b.On([]bot.Rule{bot.NewText("PING", false)}, func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
		_, err := bot.Reply(ctx, b, upd, "pong! 🏓")
		return err
	})
}

// CommandHelp registers a `/help` command that lists every command
// the bot knows about. The list is built from the supplied
// `(command, description)` pairs.
func CommandHelp(b *bot.Bot, commands map[string]string) {
	b.On([]bot.Rule{bot.NewCommand([]string{"help"}, []string{"/"}, false, false)},
		func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
			var sb strings.Builder
			sb.WriteString("🛟 Available commands:\n\n")
			for cmd, desc := range commands {
				sb.WriteString(fmt.Sprintf("• `/%s` — %s\n", cmd, desc))
			}
			_, err := bot.Reply(ctx, b, upd, sb.String())
			return err
		})
}

// MentionChat replies to every update with the chat's GUID. Useful
// for debugging / onboarding.
func MentionChat(b *bot.Bot) {
	b.On([]bot.Rule{bot.NewText("chatid", false)}, func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
		_, err := bot.Reply(ctx, b, upd, fmt.Sprintf("This chat's GUID is `%s`", upd.ChatID))
		return err
	})
}

// MentionUserID replies with the user's GUID on the message "userid".
func MentionUserID(b *bot.Bot) {
	b.On([]bot.Rule{bot.NewText("userid", false)}, func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
		var userID string
		if upd.NewMessage != nil {
			userID = upd.NewMessage.SenderID
		}
		_, err := bot.Reply(ctx, b, upd, fmt.Sprintf("Your GUID is `%s`", userID))
		return err
	})
}

// CopyMessage registers a handler that copies any received message
// back to the same chat. Useful for echoing stickers, files, etc.
func CopyMessage(b *bot.Bot) {
	b.On([]bot.Rule{bot.NewText("", false)}, func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
		if upd.NewMessage == nil {
			return nil
		}
		if upd.NewMessage.File != nil {
			_, err := b.SendFile(ctx, upd.ChatID, "", upd.NewMessage.File.FileID, upd.Text(), "", models.FileTypeFile)
			return err
		}
		_, err := bot.Reply(ctx, b, upd, upd.Text())
		return err
	})
}
