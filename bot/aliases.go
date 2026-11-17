package bot

import (
	"context"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/rules"
)

// Rule is an alias for rules.Rule so callers can use
// `bot.NewText(...)` / `bot.NewCommand(...)` without importing the
// rules sub-package explicitly.
type Rule = rules.Rule

// Update is an alias for models.Update so callers can refer to
// `*bot.Update` instead of `*bot/models.Update`.
type Update = models.Update

// InlineMessage is an alias for models.InlineMessage.
type InlineMessage = models.InlineMessage

// MessageID is an alias for models.MessageID.
type MessageID = models.MessageID

// Message is an alias for models.Message.
type Message = models.Message

// Chat is an alias for models.Chat.
type Chat = models.Chat

// BotInfo is an alias for models.Bot.
type BotInfo = models.Bot

// Keypad is an alias for models.Keypad.
type Keypad = models.Keypad

// Button is an alias for models.Button.
type Button = models.Button

// FileType is an alias for models.FileType.
type FileType = models.FileType

// ChatKeypadType is an alias for models.ChatKeypadType.
type ChatKeypadType = models.ChatKeypadType

// UpdateType is an alias for models.UpdateType.
type UpdateType = models.UpdateType

// UpdateEndpointType is an alias for models.UpdateEndpointType.
type UpdateEndpointType = models.UpdateEndpointType

// Metadata is an alias for models.Metadata.
type Metadata = models.Metadata

// ===== Convenience rules constructors (re-exported from rules pkg) =====

// NewText re-exports rules.NewText.
func NewText(pattern string, regex bool) *rules.Text { return rules.NewText(pattern, regex) }

// NewCommand re-exports rules.NewCommand.
func NewCommand(cmds, prefixes []string, caseInsensitive, allowNoPrefix bool) *rules.Command {
	return rules.NewCommand(cmds, prefixes, caseInsensitive, allowNoPrefix)
}

// NewChatID re-exports rules.NewChatID.
func NewChatID(ids ...string) *rules.ChatID { return rules.NewChatID(ids...) }

// NewUpdateType re-exports rules.NewUpdateType.
func NewUpdateType(types ...models.UpdateType) *rules.UpdateType {
	return rules.NewUpdateType(types...)
}

// NewButton re-exports rules.NewButton.
func NewButton(buttonID string, regex bool) *rules.Button { return rules.NewButton(buttonID, regex) }

// NewReply re-exports rules.NewReply.
func NewReply(toMessageID string) *rules.Reply { return rules.NewReply(toMessageID) }

// ===== Update convenience methods (free functions) ====================

// ReplyUpdate is the convenience function that replies to an update.
// Use it as `bot.ReplyUpdate(ctx, b, upd, "Hello!")`.
func ReplyUpdate(ctx context.Context, b *Bot, u *models.Update, msg string, opts ...SendOptions) (*models.MessageID, error) {
	o := SendOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	o.ReplyToMessageID = u.MessageID()
	return b.SendMessage(ctx, u.ChatID, msg, o)
}

// ReplyFileUpdate is the convenience function that replies with a file.
func ReplyFileUpdate(ctx context.Context, b *Bot, u *models.Update, filePath, msg string, fileType models.FileType, opts ...SendOptions) (*models.MessageID, error) {
	o := SendOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	o.ReplyToMessageID = u.MessageID()
	return b.SendFile(ctx, u.ChatID, filePath, "", msg, "", fileType, o)
}

// DeleteUpdate deletes the update's message.
func DeleteUpdate(ctx context.Context, b *Bot, u *models.Update) error {
	return b.DeleteMessage(ctx, u.ChatID, u.MessageID())
}

// EditTextUpdate edits the update's message text.
func EditTextUpdate(ctx context.Context, b *Bot, u *models.Update, newText string, opts ...EditOptions) error {
	return b.EditMessageText(ctx, u.ChatID, u.MessageID(), newText, opts...)
}

// EditKeypadUpdate edits the update's inline keypad.
func EditKeypadUpdate(ctx context.Context, b *Bot, u *models.Update, keypad *models.Keypad) error {
	return b.EditMessageKeypad(ctx, u.ChatID, u.MessageID(), keypad)
}
