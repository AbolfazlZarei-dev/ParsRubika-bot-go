package ParsRubika

// نسخه: 2.0.0
// سازنده: ابوالفضل زارعی
// آدرس گیت هاب: https://github.com/Abolfazl-Zarei/ParsRubika-bot-go

import (
	"fmt"
	"strings"
)

// MessageFormat فرمت‌های مختلف پیام
type MessageFormat string

const (
	FormatPlain    MessageFormat = "plain"
	FormatMarkdown MessageFormat = "markdown"
	FormatHTML     MessageFormat = "html"
)

// FormattedMessage پیام با فرمت خاص
type FormattedMessage struct {
	Text   string        `json:"text"`
	Format MessageFormat `json:"format"`
}

// NewFormattedMessage ایجاد یک پیام با فرمت خاص
func NewFormattedMessage(text string, format MessageFormat) *FormattedMessage {
	return &FormattedMessage{
		Text:   text,
		Format: format,
	}
}

// --- توابع کمکی برای فرمت‌بندی Markdown ---

// Bold متن را ضخیم می‌کند (استفاده از ** برای خوانایی بهتر)
func Bold(text string) string {
	return fmt.Sprintf("**%s**", text)
}

// Italic متن را کج می‌کند
func Italic(text string) string {
	return fmt.Sprintf("_%s_", text)
}

// Underline متن را زیرخط‌دار می‌کند (با استفاده از تگ HTML)
// نکته: زیرخط‌دار بخش استاندارد Markdown نیست و ممکن است در همه پلتفرم‌ها پشتیبانی نشود.
func Underline(text string) string {
	return fmt.Sprintf("<u>%s</u>", text)
}

// Strikethrough متن را خط‌خورده می‌کند
func Strikethrough(text string) string {
	return fmt.Sprintf("~~%s~~", text)
}

// Code متن را به صورت کد نمایش می‌دهد
func Code(text string) string {
	return fmt.Sprintf("`%s`", text)
}

// Pre متن را به صورت بلوک کد نمایش می‌دهد
func Pre(text string) string {
	return fmt.Sprintf("```%s```", text)
}

// Spoiler متن را به صورت اسپویلر نمایش می‌دهد
func Spoiler(text string) string {
	return fmt.Sprintf("||%s||", text)
}

// Link ایجاد لینک
func Link(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

// Mention ایجاد اشاره به کاربر
func Mention(username string) string {
	return fmt.Sprintf("@%s", strings.TrimPrefix(username, "@"))
}

// Hashtag ایجاد هشتگ
func Hashtag(tag string) string {
	return fmt.Sprintf("#%s", strings.TrimPrefix(tag, "#"))
}

// Quote ایجاد نقل قول
func Quote(text string) string {
	return fmt.Sprintf("> %s", text)
}

// BlockQuote ایجاد بلوک نقل قول
func BlockQuote(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for _, line := range lines {
		result.WriteString(fmt.Sprintf("> %s\n", line))
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// OrderedList ایجاد لیست مرتب (با فرمت Markdown)
func OrderedList(items []string) string {
	var result strings.Builder

	for i, item := range items {
		result.WriteString(fmt.Sprintf("%d. %s\n", i+1, item))
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// UnorderedList ایجاد لیست نامرتب (با فرمت Markdown)
func UnorderedList(items []string) string {
	var result strings.Builder

	for _, item := range items {
		result.WriteString(fmt.Sprintf("- %s\n", item))
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// --- توابع کمکی برای ساخت کیبورد ---

// NewInlineKeyboard ایجاد یک کیبورد اینلاین
func NewInlineKeyboard(rows ...[]*InlineKeyboardButton) *InlineKeyboardMarkup {
	keyboard := &InlineKeyboardMarkup{
		InlineKeyboard: make([][]*InlineKeyboardButton, len(rows)),
	}

	for i, row := range rows {
		keyboard.InlineKeyboard[i] = make([]*InlineKeyboardButton, len(row))
		copy(keyboard.InlineKeyboard[i], row)
	}

	return keyboard
}

// NewInlineKeyboardButton ایجاد یک دکمه کیبورد اینلاین
func NewInlineKeyboardButton(text string) *InlineKeyboardButton {
	return &InlineKeyboardButton{Text: text}
}

// WithURL تنظیم URL برای دکمه
func (b *InlineKeyboardButton) WithURL(url string) *InlineKeyboardButton {
	b.URL = url
	return b
}

// WithCallbackData تنظیم داده callback برای دکمه
func (b *InlineKeyboardButton) WithCallbackData(data string) *InlineKeyboardButton {
	b.CallbackData = data
	return b
}

// WithSwitchInlineQuery تنظیم سوئیچ اینلاین کوئری برای دکمه
func (b *InlineKeyboardButton) WithSwitchInlineQuery(query string) *InlineKeyboardButton {
	b.SwitchInlineQuery = query
	return b
}

// NewReplyKeyboard ایجاد یک کیبورد پاسخ
func NewReplyKeyboard(rows ...[]*KeyboardButton) *ReplyKeyboardMarkup {
	keyboard := &ReplyKeyboardMarkup{
		Keyboard:       make([][]*KeyboardButton, len(rows)),
		ResizeKeyboard: true,
	}

	for i, row := range rows {
		keyboard.Keyboard[i] = make([]*KeyboardButton, len(row))
		copy(keyboard.Keyboard[i], row)
	}

	return keyboard
}

// NewReplyKeyboardButton ایجاد یک دکمه کیبورد پاسخ
func NewReplyKeyboardButton(text string) *KeyboardButton {
	return &KeyboardButton{Text: text}
}

// WithRequestContact تنظیم درخواست شماره تلفن برای دکمه
func (b *KeyboardButton) WithRequestContact() *KeyboardButton {
	b.RequestContact = true
	return b
}

// WithRequestLocation تنظیم درخواست موقعیت مکانی برای دکمه
func (b *KeyboardButton) WithRequestLocation() *KeyboardButton {
	b.RequestLocation = true
	return b
}

// WithResizeKeyboard تنظیم تغییر اندازه کیبورد
func (kb *ReplyKeyboardMarkup) WithResizeKeyboard(resize bool) *ReplyKeyboardMarkup {
	kb.ResizeKeyboard = resize
	return kb
}

// WithOneTimeKeyboard تنظیم کیبورد یکبار مصرف
func (kb *ReplyKeyboardMarkup) WithOneTimeKeyboard(oneTime bool) *ReplyKeyboardMarkup {
	kb.OneTimeKeyboard = oneTime
	return kb
}

// NewRemoveKeyboard ایجاد یک دستور برای حذف کیبورد
func NewRemoveKeyboard() *ReplyKeyboardRemove {
	return &ReplyKeyboardRemove{
		RemoveKeyboard: true,
	}
}
