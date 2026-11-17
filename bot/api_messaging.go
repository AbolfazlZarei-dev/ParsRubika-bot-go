package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/text"
)

// SendOptions bundles the optional arguments for the various send* calls.
type SendOptions struct {
	ChatKeypad          *models.Keypad
	InlineKeypad        *models.Keypad
	DisableNotification bool
	ReplyToMessageID    string
	ChatKeypadType      models.ChatKeypadType
	ParseMode           ParseMode
	Metadata            *models.Metadata
}

// SendMessage sends a text message.
func (b *Bot) SendMessage(ctx context.Context, chatID, msg string, opts ...SendOptions) (*models.MessageID, error) {
	o := SendOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.ChatKeypadType == "" {
		o.ChatKeypadType = models.ChatKeypadNone
	}
	payload := map[string]any{
		"chat_id":              chatID,
		"text":                 msg,
		"disable_notification": o.DisableNotification,
		"chat_keypad_type":     string(o.ChatKeypadType),
	}
	if o.ReplyToMessageID != "" {
		payload["reply_to_message_id"] = o.ReplyToMessageID
	}
	if o.ChatKeypad != nil {
		payload["chat_keypad"] = o.ChatKeypad
	}
	if o.InlineKeypad != nil {
		payload["inline_keypad"] = o.InlineKeypad
	}
	if o.Metadata != nil {
		payload["metadata"] = o.Metadata
	} else if msg != "" {
		mode := o.ParseMode
		if mode == "" {
			mode = b.opts.ParseMode
		}
		if mode == "markdown" {
			md := b.formatter.Format(msg)
			mergeMd(payload, &md)
		} else if mode == "html" {
			md := b.formatter.Format(b.formatter.ToMarkdown(msg))
			mergeMd(payload, &md)
		}
	}
	data, err := b.makeRequest(ctx, "sendMessage", payload)
	if err != nil {
		return nil, err
	}
	mid := &models.MessageID{ChatID: chatID}
	if v, ok := data["message_id"]; ok {
		mid.MessageID = toString(v)
	}
	if v, ok := data["new_message_id"]; ok {
		mid.NewMessageID = toString(v)
	}
	return mid, nil
}

// SendSticker sends a sticker by sticker_id.
func (b *Bot) SendSticker(ctx context.Context, chatID, stickerID string, opts ...SendOptions) (*models.MessageID, error) {
	o := SendOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.ChatKeypadType == "" {
		o.ChatKeypadType = models.ChatKeypadNone
	}
	payload := map[string]any{
		"chat_id":              chatID,
		"sticker_id":           stickerID,
		"disable_notification": o.DisableNotification,
		"chat_keypad_type":     string(o.ChatKeypadType),
	}
	if o.ReplyToMessageID != "" {
		payload["reply_to_message_id"] = o.ReplyToMessageID
	}
	if o.ChatKeypad != nil {
		payload["chat_keypad"] = o.ChatKeypad
	}
	if o.InlineKeypad != nil {
		payload["inline_keypad"] = o.InlineKeypad
	}
	data, err := b.makeRequest(ctx, "sendSticker", payload)
	if err != nil {
		return nil, err
	}
	mid := &models.MessageID{ChatID: chatID}
	if v, ok := data["message_id"]; ok {
		mid.MessageID = toString(v)
	}
	return mid, nil
}

// SendFile sends a file. Either `filePath` or `fileID` must be set.
func (b *Bot) SendFile(ctx context.Context, chatID, filePath, fileID, msg, fileName string, fileType models.FileType, opts ...SendOptions) (*models.MessageID, error) {
	o := SendOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.ChatKeypadType == "" {
		o.ChatKeypadType = models.ChatKeypadNone
	}
	if filePath != "" {
		if fileName == "" {
			fileName = filepath.Base(filePath)
		}
		uploadURL, err := b.RequestSendFileURL(ctx, fileType)
		if err != nil {
			return nil, err
		}
		uploadedID, err := b.uploadFileTo(ctx, uploadURL, fileName, filePath)
		if err != nil {
			return nil, err
		}
		fileID = uploadedID
	}
	payload := map[string]any{
		"chat_id":              chatID,
		"file_id":              fileID,
		"text":                 msg,
		"disable_notification": o.DisableNotification,
		"chat_keypad_type":     string(o.ChatKeypadType),
	}
	if o.ReplyToMessageID != "" {
		payload["reply_to_message_id"] = o.ReplyToMessageID
	}
	if o.ChatKeypad != nil {
		payload["chat_keypad"] = o.ChatKeypad
	}
	if o.InlineKeypad != nil {
		payload["inline_keypad"] = o.InlineKeypad
	}
	if o.Metadata != nil {
		payload["metadata"] = o.Metadata
	} else if msg != "" {
		mode := o.ParseMode
		if mode == "" {
			mode = b.opts.ParseMode
		}
		if mode == "markdown" {
			md := b.formatter.Format(msg)
			mergeMd(payload, &md)
		}
	}
	data, err := b.makeRequest(ctx, "sendFile", payload)
	if err != nil {
		return nil, err
	}
	mid := &models.MessageID{ChatID: chatID, FileID: fileID}
	if v, ok := data["message_id"]; ok {
		mid.MessageID = toString(v)
	}
	return mid, nil
}

// SendPoll sends a poll.
func (b *Bot) SendPoll(ctx context.Context, chatID, question string, options []string, pollType string, allowMulti, anonymous bool, correctIdx *int, explanation, replyTo string, disableNotif bool, keypad *models.Keypad, keypadType models.ChatKeypadType) (*models.MessageID, error) {
	if pollType == "" {
		pollType = "Regular"
	}
	if keypadType == "" {
		keypadType = models.ChatKeypadNone
	}
	payload := map[string]any{
		"chat_id":                 chatID,
		"question":                question,
		"options":                 options,
		"type":                    pollType,
		"allows_multiple_answers": allowMulti,
		"is_anonymous":            anonymous,
		"correct_option_index":    correctIdx,
		"explanation":             explanation,
		"reply_to_message_id":     replyTo,
		"disable_notification":    disableNotif,
		"chat_keypad_type":        string(keypadType),
	}
	if keypad != nil {
		payload["chat_keypad"] = keypad
	}
	data, err := b.makeRequest(ctx, "sendPoll", payload)
	if err != nil {
		return nil, err
	}
	mid := &models.MessageID{ChatID: chatID}
	if v, ok := data["message_id"]; ok {
		mid.MessageID = toString(v)
	}
	return mid, nil
}

// SendLocation sends a location attachment.
func (b *Bot) SendLocation(ctx context.Context, chatID string, latitude, longitude float64, keypad *models.Keypad, disableNotif bool, replyTo string, keypadType models.ChatKeypadType) (*models.MessageID, error) {
	if keypadType == "" {
		keypadType = models.ChatKeypadNone
	}
	payload := map[string]any{
		"chat_id":              chatID,
		"latitude":             latitude,
		"longitude":            longitude,
		"disable_notification": disableNotif,
		"chat_keypad_type":     string(keypadType),
	}
	if keypad != nil {
		payload["chat_keypad"] = keypad
	}
	if replyTo != "" {
		payload["reply_to_message_id"] = replyTo
	}
	data, err := b.makeRequest(ctx, "sendLocation", payload)
	if err != nil {
		return nil, err
	}
	mid := &models.MessageID{ChatID: chatID}
	if v, ok := data["message_id"]; ok {
		mid.MessageID = toString(v)
	}
	return mid, nil
}

// SendContact sends a contact attachment.
func (b *Bot) SendContact(ctx context.Context, chatID, firstName, lastName, phone string, keypad *models.Keypad, disableNotif bool, replyTo string, keypadType models.ChatKeypadType) (*models.MessageID, error) {
	if keypadType == "" {
		keypadType = models.ChatKeypadNone
	}
	payload := map[string]any{
		"chat_id":              chatID,
		"first_name":           firstName,
		"last_name":            lastName,
		"phone_number":         phone,
		"disable_notification": disableNotif,
		"chat_keypad_type":     string(keypadType),
	}
	if keypad != nil {
		payload["chat_keypad"] = keypad
	}
	if replyTo != "" {
		payload["reply_to_message_id"] = replyTo
	}
	data, err := b.makeRequest(ctx, "sendContact", payload)
	if err != nil {
		return nil, err
	}
	mid := &models.MessageID{ChatID: chatID}
	if v, ok := data["message_id"]; ok {
		mid.MessageID = toString(v)
	}
	return mid, nil
}

// EditOptions carries optional arguments for edit operations.
type EditOptions struct {
	ParseMode ParseMode
	Metadata  *models.Metadata
}

// EditMessageText edits a message's text.
func (b *Bot) EditMessageText(ctx context.Context, chatID, messageID, msg string, opts ...EditOptions) error {
	o := EditOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	payload := map[string]any{
		"chat_id":    chatID,
		"message_id": messageID,
		"text":       msg,
	}
	if o.Metadata != nil {
		payload["metadata"] = o.Metadata
	} else if msg != "" {
		mode := o.ParseMode
		if mode == "" {
			mode = b.opts.ParseMode
		}
		if mode == "markdown" {
			md := b.formatter.Format(msg)
			mergeMd(payload, &md)
		}
	}
	_, err := b.makeRequest(ctx, "editMessageText", payload)
	return err
}

// EditMessageKeypad replaces the inline keypad of an existing message.
func (b *Bot) EditMessageKeypad(ctx context.Context, chatID, messageID string, keypad *models.Keypad) error {
	_, err := b.makeRequest(ctx, "editMessageKeypad", map[string]any{
		"chat_id":       chatID,
		"message_id":    messageID,
		"inline_keypad": keypad,
	})
	return err
}

// DeleteMessage removes a message.
func (b *Bot) DeleteMessage(ctx context.Context, chatID, messageID string) error {
	_, err := b.makeRequest(ctx, "deleteMessage", map[string]any{
		"chat_id":    chatID,
		"message_id": messageID,
	})
	return err
}

// ===== Helpers ==========================================================

func mergeMd(payload map[string]any, md *text.Result) {
	if md == nil {
		return
	}
	if md.Text != "" {
		payload["text"] = md.Text
	}
	if md.Metadata != nil {
		payload["metadata"] = md.Metadata
	}
}

func toString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		if float64(int64(x)) == x {
			return fmt.Sprintf("%d", int64(x))
		}
		return fmt.Sprintf("%f", x)
	case int:
		return fmt.Sprintf("%d", x)
	case int64:
		return fmt.Sprintf("%d", x)
	case bool:
		if x {
			return "true"
		}
		return "false"
	case nil:
		return ""
	default:
		b, _ := json.Marshal(x)
		return string(b)
	}
}

// uploadFileTo uploads a local file to the supplied upload URL using
// multipart/form-data.
func (b *Bot) uploadFileTo(ctx context.Context, uploadURL, fileName, filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	body, w := createMultipartBody(fileName, f)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := b.transport.HTTPClient().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", &APIError{Status: fmt.Sprintf("%d", resp.StatusCode), DevMessage: string(raw)}
	}
	var parsed map[string]any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", err
	}
	if d, ok := parsed["data"].(map[string]any); ok {
		if id, ok := d["file_id"].(string); ok {
			return id, nil
		}
	}
	return "", &APIError{Status: "InvalidResponse", DevMessage: "missing file_id in upload response"}
}

// createMultipartBody builds a multipart/form-data body for a single
// file. It returns a pipe reader (so the body is streamed, not buffered
// in memory) and the writer whose FormDataContentType() should be
// used as the Content-Type header.
func createMultipartBody(fileName string, f io.Reader) (reader io.Reader, w *multipart.Writer) {
	r, pw := io.Pipe()
	w = multipart.NewWriter(pw)
	go func() {
		defer pw.Close()
		fw, err := w.CreateFormFile("file", fileName)
		if err != nil {
			return
		}
		_, _ = io.Copy(fw, f)
		_ = w.Close()
	}()
	return r, w
}
