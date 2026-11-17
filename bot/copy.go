package bot

import (
	"context"
	"errors"
	"fmt"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// ErrNoMessage is returned by helpers that need a message the update
// does not carry.
var ErrNoMessage = errors.New("bot: update carries no message")

// ErrNothingToCopy is returned when a message has no copyable content.
var ErrNothingToCopy = errors.New("bot: message has no copyable content")

// CopyMessage re-sends the content of msg to toChatID without any
// "forwarded from" attribution. Text, files (by file_id), stickers,
// locations, contacts and polls are supported.
func (b *Bot) CopyMessage(ctx context.Context, msg *models.Message, toChatID string, opts ...SendOptions) (*models.MessageID, error) {
	if msg == nil {
		return nil, ErrNoMessage
	}
	o := SendOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	switch {
	case msg.File != nil && msg.File.FileID != "":
		return b.SendFile(ctx, toChatID, "", msg.File.FileID, msg.Text, msg.File.FileName, models.FileTypeFile, o)
	case msg.Sticker != nil && msg.Sticker.StickerID != "":
		return b.SendSticker(ctx, toChatID, msg.Sticker.StickerID, o)
	case msg.Location != nil:
		lat, lon := parseFloat(msg.Location.Latitude), parseFloat(msg.Location.Longitude)
		return b.SendLocation(ctx, toChatID, lat, lon, o.InlineKeypad, o.DisableNotification, o.ReplyToMessageID, o.ChatKeypadType)
	case msg.ContactMessage != nil:
		c := msg.ContactMessage
		return b.SendContact(ctx, toChatID, c.FirstName, c.LastName, c.PhoneNumber, o.InlineKeypad, o.DisableNotification, o.ReplyToMessageID, o.ChatKeypadType)
	case msg.Poll != nil:
		return b.SendPoll(ctx, toChatID, msg.Poll.Question, msg.Poll.Options, "Regular", false, true, nil, "", o.ReplyToMessageID, o.DisableNotification, o.InlineKeypad, o.ChatKeypadType)
	case msg.Text != "":
		if o.Metadata == nil {
			o.Metadata = msg.Metadata
		}
		return b.SendMessage(ctx, toChatID, msg.Text, o)
	}
	return nil, ErrNothingToCopy
}

func parseFloat(s string) float64 {
	var f float64
	_, _ = fmt.Sscan(s, &f)
	return f
}
