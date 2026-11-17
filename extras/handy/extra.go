package handy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// ===== Extra handy helpers ====================
//
// More ergonomic helpers for the most common bot patterns.

// StartCommand registers a `/start` handler that sends a welcome message.
// The welcome message can include the user's name via {name} placeholder.
func StartCommand(b *bot.Bot, welcomeMessage string) {
	b.On([]bot.Rule{bot.NewCommand([]string{"start"}, []string{"/"}, false, false)},
		func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
			msg := welcomeMessage
			if upd.NewMessage != nil {
				name := upd.NewMessage.SenderID
				msg = strings.ReplaceAll(msg, "{name}", name)
			}
			_, err := bot.Reply(ctx, b, upd, msg)
			return err
		})
}

// IDCommand registers an `/id` command that returns the chat ID and
// sender ID.
func IDCommand(b *bot.Bot) {
	b.On([]bot.Rule{bot.NewCommand([]string{"id"}, []string{"/"}, false, false)},
		func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
			var senderID string
			if upd.NewMessage != nil {
				senderID = upd.NewMessage.SenderID
			}
			text := fmt.Sprintf("Chat ID: `%s`\nSender ID: `%s`", upd.ChatID, senderID)
			_, err := bot.Reply(ctx, b, upd, text)
			return err
		})
}

// TimeCommand registers a `/time` command that returns the current server
// time.
func TimeCommand(b *bot.Bot) {
	b.On([]bot.Rule{bot.NewCommand([]string{"time"}, []string{"/"}, false, false)},
		func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
			_, err := bot.Reply(ctx, b, upd, fmt.Sprintf("🕐 %s", time.Now().Format("2006-01-02 15:04:05")))
			return err
		})
}

// PingCommand registers a `/ping` command that measures response time
// and reports it back to the user.
func PingCommand(b *bot.Bot) {
	b.On([]bot.Rule{bot.NewCommand([]string{"ping"}, []string{"/"}, false, false)},
		func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
			start := time.Now()
			elapsed := time.Since(start)
			_, err := bot.Reply(ctx, b, upd, fmt.Sprintf("🏓 pong! (%s)", elapsed))
			return err
		})
}

// UserInfoCommand registers a `/whoami` command that returns the user's
// info.
func UserInfoCommand(b *bot.Bot) {
	b.On([]bot.Rule{bot.NewCommand([]string{"whoami"}, []string{"/"}, false, false)},
		func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
			if upd.NewMessage == nil {
				return nil
			}
			var sb strings.Builder
			sb.WriteString("👤 Your info:\n")
			sb.WriteString(fmt.Sprintf("• Sender ID: `%s`\n", upd.NewMessage.SenderID))
			sb.WriteString(fmt.Sprintf("• Sender Type: `%s`\n", upd.NewMessage.SenderType))
			sb.WriteString(fmt.Sprintf("• Chat ID: `%s`\n", upd.ChatID))
			_, err := bot.Reply(ctx, b, upd, sb.String())
			return err
		})
}

// ReplyToMedia registers a handler that replies to any media message
// with the file ID.
func ReplyToMedia(b *bot.Bot) {
	b.On([]bot.Rule{bot.NewText("", false)},
		func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
			if upd.NewMessage == nil || upd.NewMessage.File == nil {
				return nil
			}
			text := fmt.Sprintf("📎 File ID: `%s`", upd.NewMessage.File.FileID)
			_, err := bot.Reply(ctx, b, upd, text)
			return err
		})
}

// AdminOnly wraps a handler so it only runs for chat admins. Non-admins
// receive the rejectMessage (if non-empty).
func AdminOnly(b *bot.Bot, rejectMessage string, handler func(ctx context.Context, b *bot.Bot, upd *bot.Update) error) func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
	return func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
		if upd.NewMessage == nil {
			return nil
		}
		isAdmin, err := b.MemberIsAdmin(ctx, upd.ChatID, upd.NewMessage.SenderID)
		if err != nil {
			return err
		}
		if !isAdmin {
			if rejectMessage != "" {
				_, _ = bot.Reply(ctx, b, upd, rejectMessage)
			}
			return nil
		}
		return handler(ctx, b, upd)
	}
}

// PrivateOnly wraps a handler so it only runs in private chats.
func PrivateOnly(handler func(ctx context.Context, b *bot.Bot, upd *bot.Update) error) func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
	return func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
		if !bot.IsPrivateUpdate(upd) {
			return nil
		}
		return handler(ctx, b, upd)
	}
}

// GroupOnly wraps a handler so it only runs in groups.
func GroupOnly(handler func(ctx context.Context, b *bot.Bot, upd *bot.Update) error) func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
	return func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
		if !bot.IsGroupUpdate(upd) {
			return nil
		}
		return handler(ctx, b, upd)
	}
}

// ChannelOnly wraps a handler so it only runs in channels.
func ChannelOnly(handler func(ctx context.Context, b *bot.Bot, upd *bot.Update) error) func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
	return func(ctx context.Context, b *bot.Bot, upd *bot.Update) error {
		if !bot.IsChannelUpdate(upd) {
			return nil
		}
		return handler(ctx, b, upd)
	}
}

// SendTyping sends a "typing" chat action. This is a no-op for the bot
// API (which does not support chat activities), but is kept so code can be shared
// with the account client.
func SendTyping(ctx context.Context, b *bot.Bot, chatID string) error {
	_ = ctx
	_ = b
	_ = chatID
	return nil
}

// ForwardToChat forwards a message to another chat.
func ForwardToChat(ctx context.Context, b *bot.Bot, fromChatID, messageID, toChatID string) error {
	_, err := b.ForwardMessage(ctx, fromChatID, messageID, toChatID, false)
	return err
}

// SendKeypad sends a message with a simple inline keypad built from a
// list of (id, text) button pairs arranged in rows.
func SendKeypad(ctx context.Context, b *bot.Bot, chatID, text string, buttons [][]string) error {
	var rows []models.KeypadRow
	for _, row := range buttons {
		var btns []models.Button
		for _, b := range row {
			btns = append(btns, models.Button{
				ID:         b,
				Type:       models.ButtonTypeSimple,
				ButtonText: b,
			})
		}
		rows = append(rows, models.KeypadRow{Buttons: btns})
	}
	kp := &models.Keypad{Rows: rows}
	_, err := b.SendMessage(ctx, chatID, text, bot.SendOptions{InlineKeypad: kp})
	return err
}

// QuickReply sends a text reply to an update and returns the message ID.
func QuickReply(ctx context.Context, b *bot.Bot, upd *bot.Update, text string) (string, error) {
	mid, err := bot.Reply(ctx, b, upd, text)
	if err != nil {
		return "", err
	}
	return mid.MessageID, nil
}

// DeleteAfter sends a message and schedules its deletion after delay.
func DeleteAfter(ctx context.Context, b *bot.Bot, chatID, text string, delay time.Duration) (bot.CancelFunc, error) {
	mid, err := b.SendMessage(ctx, chatID, text)
	if err != nil {
		return nil, err
	}
	return b.AutoDeleteMessage(ctx, chatID, mid.MessageID, delay), nil
}
