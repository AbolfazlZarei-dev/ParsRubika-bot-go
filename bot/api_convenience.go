package bot

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// SendImage sends an image file. Either `filePath` (local file to
// upload) or `fileID` (previously uploaded file ID) must be set.
// `text` is the optional caption.
func (b *Bot) SendImage(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error) {
	return b.SendFile(ctx, chatID, filePath, fileID, text, fileName, models.FileTypeImage, opts...)
}

// SendVideo sends a video file. Either `filePath` or `fileID` must be
// set. `text` is the optional caption.
func (b *Bot) SendVideo(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error) {
	return b.SendFile(ctx, chatID, filePath, fileID, text, fileName, models.FileTypeVideo, opts...)
}

// SendGif sends a GIF file. Either `filePath` or `fileID` must be set.
// `text` is the optional caption.
func (b *Bot) SendGif(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error) {
	return b.SendFile(ctx, chatID, filePath, fileID, text, fileName, models.FileTypeGif, opts...)
}

// SendAudio sends an audio file. Either `filePath` or `fileID` must be
// set. `text` is the optional caption.
func (b *Bot) SendAudio(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error) {
	return b.SendFile(ctx, chatID, filePath, fileID, text, fileName, models.FileTypeMusic, opts...)
}

// SendMusic is an alias for SendAudio (kept for backwards
// compatibility). It sends an audio file as a music message.
func (b *Bot) SendMusic(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error) {
	return b.SendAudio(ctx, chatID, filePath, fileID, text, fileName, opts...)
}

// SendVoice sends a voice message. The file name is forced to `.ogg`
// so the client renders it as a voice player.
func (b *Bot) SendVoice(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error) {
	if filePath != "" {
		base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath)) + ".ogg"
		if fileName == "" {
			fileName = base
		}
	}
	return b.SendFile(ctx, chatID, filePath, fileID, text, fileName, models.FileTypeVoice, opts...)
}

// SendMessageSimple is a one-liner for SendMessage with just text.
func (b *Bot) SendMessageSimple(ctx context.Context, chatID, text string) (*models.MessageID, error) {
	return b.SendMessage(ctx, chatID, text)
}

// SendDocument sends a generic document/file. It is an alias for SendFile
// with FileTypeFile.
func (b *Bot) SendDocument(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error) {
	return b.SendFile(ctx, chatID, filePath, fileID, text, fileName, models.FileTypeFile, opts...)
}

// SendVideoMessage sends a round video message (VideoMessage type).
// On the bot API side this is sent as a Video file; the round rendering
// is determined by the client.
func (b *Bot) SendVideoMessage(ctx context.Context, chatID, filePath, fileID, text, fileName string, opts ...SendOptions) (*models.MessageID, error) {
	return b.SendFile(ctx, chatID, filePath, fileID, text, fileName, models.FileTypeVideo, opts...)
}

// SendContactReply sends a contact as a reply to an update.
func (b *Bot) SendContactReply(ctx context.Context, upd *models.Update, firstName, lastName, phone string) (*models.MessageID, error) {
	return b.SendContact(ctx, upd.ChatID, firstName, lastName, phone, nil, false, upd.MessageID(), models.ChatKeypadNone)
}

// SendLocationReply sends a location as a reply to an update.
func (b *Bot) SendLocationReply(ctx context.Context, upd *models.Update, latitude, longitude float64) (*models.MessageID, error) {
	return b.SendLocation(ctx, upd.ChatID, latitude, longitude, nil, false, upd.MessageID(), models.ChatKeypadNone)
}

// SendPollReply sends a poll as a reply to an update.
func (b *Bot) SendPollReply(ctx context.Context, upd *models.Update, question string, options []string, pollType string) (*models.MessageID, error) {
	return b.SendPoll(ctx, upd.ChatID, question, options, pollType, false, true, nil, "", upd.MessageID(), false, nil, models.ChatKeypadNone)
}

// SendStickerReply sends a sticker as a reply to an update.
func (b *Bot) SendStickerReply(ctx context.Context, upd *models.Update, stickerID string) (*models.MessageID, error) {
	o := SendOptions{ReplyToMessageID: upd.MessageID()}
	return b.SendSticker(ctx, upd.ChatID, stickerID, o)
}

// Reply is a convenience helper that replies to a given update. It is
// a free function (because we cannot add methods to a non-local type
// from another package) — use it as `bot.Reply(ctx, b, upd, "text")`.
func Reply(ctx context.Context, b *Bot, u *models.Update, text string, opts ...SendOptions) (*models.MessageID, error) {
	o := SendOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	o.ReplyToMessageID = u.MessageID()
	return b.SendMessage(ctx, u.ChatID, text, o)
}

// ReplyFile is a convenience helper that replies with a file.
func ReplyFile(ctx context.Context, b *Bot, u *models.Update, filePath, text string, fileType models.FileType, opts ...SendOptions) (*models.MessageID, error) {
	o := SendOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	o.ReplyToMessageID = u.MessageID()
	return b.SendFile(ctx, u.ChatID, filePath, "", text, "", fileType, o)
}

// Delete is a convenience helper that deletes the update's message.
func Delete(ctx context.Context, b *Bot, u *models.Update) error {
	return b.DeleteMessage(ctx, u.ChatID, u.MessageID())
}

// EditText is a convenience helper that edits the update's message.
func EditText(ctx context.Context, b *Bot, u *models.Update, newText string, opts ...EditOptions) error {
	return b.EditMessageText(ctx, u.ChatID, u.MessageID(), newText, opts...)
}

// EditKeypad is a convenience helper that edits the update's keypad.
func EditKeypad(ctx context.Context, b *Bot, u *models.Update, keypad *models.Keypad) error {
	return b.EditMessageKeypad(ctx, u.ChatID, u.MessageID(), keypad)
}

// IsGroupUpdate reports whether the update came from a group chat.
func IsGroupUpdate(u *models.Update) bool {
	return strings.HasPrefix(u.ChatID, "g0")
}

// IsChannelUpdate reports whether the update came from a channel.
func IsChannelUpdate(u *models.Update) bool {
	return strings.HasPrefix(u.ChatID, "c0")
}

// IsPrivateUpdate reports whether the update came from a private
// (user-to-bot) chat. In Rubika, user object GUIDs start with "u0";
// "b0" is the bot GUID prefix, not a private-chat marker.
func IsPrivateUpdate(u *models.Update) bool {
	return strings.HasPrefix(u.ChatID, "u0")
}
