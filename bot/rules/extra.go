package rules

import (
	"context"
	"regexp"
	"strings"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// ===== Text helper rules ===============================================

// TextStartsWith matches messages whose text starts with `prefix`.
type TextStartsWith struct{ Prefix string }

// Check implements Rule.
func (t TextStartsWith) Check(_ context.Context, upd *models.Update) bool {
	return strings.HasPrefix(upd.Text(), t.Prefix)
}

// TextEndsWith matches messages whose text ends with `suffix`.
type TextEndsWith struct{ Suffix string }

// Check implements Rule.
func (t TextEndsWith) Check(_ context.Context, upd *models.Update) bool {
	return strings.HasSuffix(upd.Text(), t.Suffix)
}

// TextContains matches messages whose text contains `substring`
// (case-insensitive by default).
type TextContains struct {
	Substring     string
	CaseSensitive bool
}

// Check implements Rule.
func (t TextContains) Check(_ context.Context, upd *models.Update) bool {
	text := upd.Text()
	sub := t.Substring
	if !t.CaseSensitive {
		text = strings.ToLower(text)
		sub = strings.ToLower(sub)
	}
	return strings.Contains(text, sub)
}

// TextEquals matches messages whose text exactly equals `value`.
type TextEquals struct{ Value string }

// Check implements Rule.
func (t TextEquals) Check(_ context.Context, upd *models.Update) bool {
	return upd.Text() == t.Value
}

// TextNotEquals matches messages whose text does NOT equal `value`.
type TextNotEquals struct{ Value string }

// Check implements Rule.
func (t TextNotEquals) Check(_ context.Context, upd *models.Update) bool {
	return upd.Text() != t.Value
}

// TextLength matches messages whose text length is within [min, max].
type TextLength struct{ Min, Max int }

// Check implements Rule.
func (t TextLength) Check(_ context.Context, upd *models.Update) bool {
	l := len(upd.Text())
	if t.Min > 0 && l < t.Min {
		return false
	}
	if t.Max > 0 && l > t.Max {
		return false
	}
	return true
}

// TextWordCount matches messages whose word count is within [min, max].
type TextWordCount struct{ Min, Max int }

// Check implements Rule.
func (t TextWordCount) Check(_ context.Context, upd *models.Update) bool {
	words := len(strings.Fields(upd.Text()))
	if t.Min > 0 && words < t.Min {
		return false
	}
	if t.Max > 0 && words > t.Max {
		return false
	}
	return true
}

// TextIsUpper matches messages whose text is all uppercase.
type TextIsUpper struct{}

// Check implements Rule.
func (TextIsUpper) Check(_ context.Context, upd *models.Update) bool {
	return upd.Text() != "" && upd.Text() == strings.ToUpper(upd.Text())
}

// TextIsLower matches messages whose text is all lowercase.
type TextIsLower struct{}

// Check implements Rule.
func (TextIsLower) Check(_ context.Context, upd *models.Update) bool {
	return upd.Text() != "" && upd.Text() == strings.ToLower(upd.Text())
}

// TextIsDigit matches messages whose text is all digits.
type TextIsDigit struct{}

// Check implements Rule.
func (TextIsDigit) Check(_ context.Context, upd *models.Update) bool {
	return upd.Text() != "" && allDigits(upd.Text())
}

func allDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// TextRegex matches messages whose text matches a regex.
type TextRegex struct {
	Pattern string
	re      *regexp.Regexp
}

// NewTextRegex constructs a TextRegex rule.
func NewTextRegex(pattern string) *TextRegex {
	r := &TextRegex{Pattern: pattern}
	r.re, _ = regexp.Compile(pattern)
	return r
}

// Check implements Rule.
func (t *TextRegex) Check(_ context.Context, upd *models.Update) bool {
	if t.re == nil {
		return false
	}
	return t.re.MatchString(upd.Text())
}

// ===== File helper rules ===============================================

// FileSizeGT matches messages whose file size is greater than `size` bytes.
type FileSizeGT struct{ Size int64 }

// Check implements Rule.
func (f FileSizeGT) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.File == nil {
		return false
	}
	return fileSize(upd.NewMessage.File) > f.Size
}

// FileSizeLT matches messages whose file size is less than `size` bytes.
type FileSizeLT struct{ Size int64 }

// Check implements Rule.
func (f FileSizeLT) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.File == nil {
		return false
	}
	return fileSize(upd.NewMessage.File) < f.Size
}

// FileNameContains matches messages whose file name contains `substring`.
type FileNameContains struct{ Substring string }

// Check implements Rule.
func (f FileNameContains) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.File == nil {
		return false
	}
	return strings.Contains(strings.ToLower(upd.NewMessage.File.FileName), strings.ToLower(f.Substring))
}

// FileExtension matches messages whose file name ends with `ext`.
type FileExtension struct{ Ext string }

// Check implements Rule.
func (f FileExtension) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.File == nil {
		return false
	}
	return strings.HasSuffix(strings.ToLower(upd.NewMessage.File.FileName), strings.ToLower(f.Ext))
}

// FileIDIs matches messages whose file_id equals `id`.
type FileIDIs struct{ ID string }

// Check implements Rule.
func (f FileIDIs) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.File == nil {
		return false
	}
	return upd.NewMessage.File.FileID == f.ID
}

// fileSize extracts the file size from a File model.
func fileSize(f *models.File) int64 {
	if f == nil || f.Size == "" {
		return 0
	}
	var n int64
	for _, r := range f.Size {
		if r >= '0' && r <= '9' {
			n = n*10 + int64(r-'0')
		}
	}
	return n
}

// ===== Sticker helper rules ============================================

// StickerIDIs matches messages whose sticker_id equals `id`.
type StickerIDIs struct{ ID string }

// Check implements Rule.
func (s StickerIDIs) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.Sticker == nil {
		return false
	}
	return upd.NewMessage.Sticker.StickerID == s.ID
}

// StickerEmojiIs matches messages whose sticker emoji equals `emoji`.
type StickerEmojiIs struct{ Emoji string }

// Check implements Rule.
func (s StickerEmojiIs) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.Sticker == nil {
		return false
	}
	return upd.NewMessage.Sticker.EmojiCharacter == s.Emoji
}

// ===== Chat helper rules ==============================================

// ChatTitleContains matches when the chat title contains `keyword`.
// Note: this requires the chat info to be fetched first, which the
// rule does lazily.
type ChatTitleContains struct{ Keyword string }

// Check implements Rule.
func (c ChatTitleContains) Check(_ context.Context, upd *models.Update) bool {
	// We don't have chat info in the update itself; the caller must
	// pre-fetch it and store it in ExternalData. This is a best-effort
	// check.
	if upd.ExternalData == nil {
		return false
	}
	title, _ := upd.ExternalData["chat_title"].(string)
	return strings.Contains(strings.ToLower(title), strings.ToLower(c.Keyword))
}

// ChatTypeIs matches when the update came from a chat of the given type.
// `chatType` is "user", "private", "group", "channel", or "bot".
type ChatTypeIs struct{ ChatType string }

// Check implements Rule.
func (c ChatTypeIs) Check(_ context.Context, upd *models.Update) bool {
	switch strings.ToLower(c.ChatType) {
	case "user", "private":
		return strings.HasPrefix(upd.ChatID, "u0")
	case "group":
		return strings.HasPrefix(upd.ChatID, "g0")
	case "channel":
		return strings.HasPrefix(upd.ChatID, "c0")
	case "bot":
		return strings.HasPrefix(upd.ChatID, "b0")
	}
	return false
}

// ===== Forwarded-from rules ===========================================

// ForwardedFromUser matches messages forwarded from a user.
type ForwardedFromUser struct{}

// Check implements Rule.
func (ForwardedFromUser) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.ForwardedFrom == nil {
		return false
	}
	return upd.NewMessage.ForwardedFrom.TypeFrom == "User"
}

// ForwardedFromChannel matches messages forwarded from a channel.
type ForwardedFromChannel struct{}

// Check implements Rule.
func (ForwardedFromChannel) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.ForwardedFrom == nil {
		return false
	}
	return upd.NewMessage.ForwardedFrom.TypeFrom == "Channel"
}

// ForwardedFromBot matches messages forwarded from a bot.
type ForwardedFromBot struct{}

// Check implements Rule.
func (ForwardedFromBot) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.ForwardedFrom == nil {
		return false
	}
	return upd.NewMessage.ForwardedFrom.TypeFrom == "Bot"
}

// ForwardedNoLink matches messages forwarded without a link to the
// original sender.
type ForwardedNoLink struct{}

// Check implements Rule.
func (ForwardedNoLink) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil {
		return false
	}
	return upd.NewMessage.ForwardedNoLink != ""
}

// ===== EventData rules (bot-side) ======================================

// IsEventData matches EventData updates (BotJoined, BotRemoved, etc.).
type IsEventData struct{ EventTypes []string }

// Check implements Rule.
func (e IsEventData) Check(_ context.Context, upd *models.Update) bool {
	// EventData updates have a specific type or carry event_data
	if upd.Type == "EventData" {
		if len(e.EventTypes) == 0 {
			return true
		}
		// Check event_type against the list
		for _, t := range e.EventTypes {
			if upd.ExternalData != nil {
				if et, _ := upd.ExternalData["event_type"].(string); et == t {
					return true
				}
			}
		}
		return false
	}
	return false
}

// BotJoined matches BotJoined events.
type BotJoined struct{}

// Check implements Rule.
func (BotJoined) Check(_ context.Context, upd *models.Update) bool {
	return IsEventData{EventTypes: []string{"BotJoined"}}.Check(context.Background(), upd)
}

// BotRemoved matches BotRemoved events.
type BotRemoved struct{}

// Check implements Rule.
func (BotRemoved) Check(_ context.Context, upd *models.Update) bool {
	return IsEventData{EventTypes: []string{"BotRemoved"}}.Check(context.Background(), upd)
}

// BotPermissionsChanged matches BotPermissionsChanged events.
type BotPermissionsChanged struct{}

// Check implements Rule.
func (BotPermissionsChanged) Check(_ context.Context, upd *models.Update) bool {
	return IsEventData{EventTypes: []string{"BotPermissionsChanged"}}.Check(context.Background(), upd)
}

// ===== MetadataType rule ===============================================

// MetadataType matches messages that contain metadata of the given type(s).
type MetadataType struct{ Types []string }

// Check implements Rule.
func (m MetadataType) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.Metadata == nil {
		return false
	}
	for _, part := range upd.NewMessage.Metadata.MetaDataParts {
		for _, t := range m.Types {
			if part.Type == t {
				return true
			}
		}
	}
	return false
}

// ===== Text formatting rules (IsBold, IsItalic, etc.) ==================

// IsBold matches messages that contain bold formatting.
type IsBold struct{}

// Check implements Rule.
func (IsBold) Check(_ context.Context, upd *models.Update) bool {
	return MetadataType{Types: []string{"Bold"}}.Check(context.Background(), upd)
}

// IsItalic matches messages that contain italic formatting.
type IsItalic struct{}

// Check implements Rule.
func (IsItalic) Check(_ context.Context, upd *models.Update) bool {
	return MetadataType{Types: []string{"Italic"}}.Check(context.Background(), upd)
}

// IsStrike matches messages that contain strike formatting.
type IsStrike struct{}

// Check implements Rule.
func (IsStrike) Check(_ context.Context, upd *models.Update) bool {
	return MetadataType{Types: []string{"Strike"}}.Check(context.Background(), upd)
}

// IsUnderline matches messages that contain underline formatting.
type IsUnderline struct{}

// Check implements Rule.
func (IsUnderline) Check(_ context.Context, upd *models.Update) bool {
	return MetadataType{Types: []string{"Underline"}}.Check(context.Background(), upd)
}

// IsQuote matches messages that contain a quote block.
type IsQuote struct{}

// Check implements Rule.
func (IsQuote) Check(_ context.Context, upd *models.Update) bool {
	return MetadataType{Types: []string{"Quote"}}.Check(context.Background(), upd)
}

// IsSpoiler matches messages that contain spoiler formatting.
type IsSpoiler struct{}

// Check implements Rule.
func (IsSpoiler) Check(_ context.Context, upd *models.Update) bool {
	return MetadataType{Types: []string{"Spoiler"}}.Check(context.Background(), upd)
}

// IsPre matches messages that contain a pre-formatted code block.
type IsPre struct{}

// Check implements Rule.
func (IsPre) Check(_ context.Context, upd *models.Update) bool {
	return MetadataType{Types: []string{"Pre"}}.Check(context.Background(), upd)
}

// IsMono matches messages that contain inline code formatting.
type IsMono struct{}

// Check implements Rule.
func (IsMono) Check(_ context.Context, upd *models.Update) bool {
	return MetadataType{Types: []string{"Mono"}}.Check(context.Background(), upd)
}
