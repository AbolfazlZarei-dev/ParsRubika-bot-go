package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// GetMe returns the bot identity.
func (b *Bot) GetMe(ctx context.Context) (*models.Bot, error) {
	data, err := b.makeRequest(ctx, "getMe", map[string]any{})
	if err != nil {
		return nil, err
	}
	botRaw, ok := data["bot"].(map[string]any)
	if !ok {
		return nil, &APIError{Status: "InvalidResponse", DevMessage: "missing 'bot' field"}
	}
	raw, _ := json.Marshal(botRaw)
	var bot models.Bot
	if err := json.Unmarshal(raw, &bot); err != nil {
		return nil, err
	}
	return &bot, nil
}

// SetCommands registers the bot's command list.
func (b *Bot) SetCommands(ctx context.Context, cmds []models.BotCommand) error {
	_, err := b.makeRequest(ctx, "setCommands", map[string]any{"bot_commands": cmds})
	return err
}

// UpdateBotEndpoints registers or unregisters the webhook endpoints.
func (b *Bot) UpdateBotEndpoints(ctx context.Context, url, endpointType string) (map[string]any, error) {
	return b.makeRequest(ctx, "updateBotEndpoints", map[string]any{
		"url":  url,
		"type": endpointType,
	})
}

// GetChat returns the chat object for `chatID`.
func (b *Bot) GetChat(ctx context.Context, chatID string) (*models.Chat, error) {
	data, err := b.makeRequest(ctx, "getChat", map[string]any{"chat_id": chatID})
	if err != nil {
		return nil, err
	}
	chatRaw, ok := data["chat"].(map[string]any)
	if !ok {
		return nil, &APIError{Status: "InvalidResponse", DevMessage: "missing 'chat' field"}
	}
	raw, _ := json.Marshal(chatRaw)
	var c models.Chat
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// ForwardMessage forwards a message to another chat.
func (b *Bot) ForwardMessage(ctx context.Context, fromChatID, messageID, toChatID string, disableNotif bool) (*models.MessageID, error) {
	data, err := b.makeRequest(ctx, "forwardMessage", map[string]any{
		"from_chat_id":         fromChatID,
		"message_id":           messageID,
		"to_chat_id":           toChatID,
		"disable_notification": disableNotif,
	})
	if err != nil {
		return nil, err
	}
	mid := &models.MessageID{ChatID: toChatID}
	if v, ok := data["message_id"]; ok {
		mid.MessageID = toString(v)
	}
	return mid, nil
}

// GetChatAdministrators returns the admin list of a chat.
func (b *Bot) GetChatAdministrators(ctx context.Context, chatID string) ([]map[string]any, error) {
	data, err := b.makeRequest(ctx, "getChatAdministrators", map[string]any{"chat_id": chatID})
	if err != nil {
		return nil, err
	}
	if admins, ok := data["administrators"].([]any); ok {
		out := make([]map[string]any, 0, len(admins))
		for _, a := range admins {
			if m, ok := a.(map[string]any); ok {
				out = append(out, m)
			}
		}
		return out, nil
	}
	return nil, nil
}

// GetChatMemberCount returns the member count of a chat.
func (b *Bot) GetChatMemberCount(ctx context.Context, chatID string) (int, error) {
	data, err := b.makeRequest(ctx, "getChatMemberCount", map[string]any{"chat_id": chatID})
	if err != nil {
		return 0, err
	}
	switch n := data["count"].(type) {
	case float64:
		return int(n), nil
	case int:
		return n, nil
	case int64:
		return int(n), nil
	}
	return 0, &APIError{Status: "InvalidResponse", DevMessage: "missing 'count' field"}
}

// BanChatMember bans a member.
func (b *Bot) BanChatMember(ctx context.Context, chatID, userID string) error {
	_, err := b.makeRequest(ctx, "banChatMember", map[string]any{
		"chat_id": chatID,
		"user_id": userID,
	})
	return err
}

// UnbanChatMember revokes a ban.
func (b *Bot) UnbanChatMember(ctx context.Context, chatID, userID string) error {
	_, err := b.makeRequest(ctx, "unbanChatMember", map[string]any{
		"chat_id": chatID,
		"user_id": userID,
	})
	return err
}

// GetChatMember fetches a single chat member.
func (b *Bot) GetChatMember(ctx context.Context, chatID, userID string) (map[string]any, error) {
	return b.makeRequest(ctx, "getChatMember", map[string]any{
		"chat_id": chatID,
		"user_id": userID,
	})
}

// PinChatMessage pins a message.
func (b *Bot) PinChatMessage(ctx context.Context, chatID, messageID string, disableNotif bool) error {
	_, err := b.makeRequest(ctx, "pinChatMessage", map[string]any{
		"chat_id":              chatID,
		"message_id":           messageID,
		"disable_notification": disableNotif,
	})
	return err
}

// UnpinChatMessage unpins a message.
func (b *Bot) UnpinChatMessage(ctx context.Context, chatID, messageID string) error {
	_, err := b.makeRequest(ctx, "unpinChatMessage", map[string]any{
		"chat_id":    chatID,
		"message_id": messageID,
	})
	return err
}

// UnpinAllChatMessages unpins all messages in a chat. The Rubika
// API only accepts a chat_id for this operation.
func (b *Bot) UnpinAllChatMessages(ctx context.Context, chatID string) error {
	_, err := b.makeRequest(ctx, "unpinAllChatMessages", map[string]any{
		"chat_id": chatID,
	})
	return err
}

// EditChatKeypad edits the chat keypad.
func (b *Bot) EditChatKeypad(ctx context.Context, chatID string, keypadType models.ChatKeypadType, keypad *models.Keypad) error {
	payload := map[string]any{
		"chat_id":          chatID,
		"chat_keypad_type": string(keypadType),
	}
	if keypad != nil {
		payload["chat_keypad"] = keypad
	}
	_, err := b.makeRequest(ctx, "editChatKeypad", payload)
	return err
}

// RequestSendFileURL requests an upload URL for the supplied file type.
func (b *Bot) RequestSendFileURL(ctx context.Context, fileType models.FileType) (string, error) {
	data, err := b.makeRequest(ctx, "requestSendFile", map[string]any{"type": string(fileType)})
	if err != nil {
		return "", err
	}
	url, _ := data["upload_url"].(string)
	return url, nil
}

// GetFile returns the download URL of a previously uploaded file.
func (b *Bot) GetFile(ctx context.Context, fileID string) (string, error) {
	data, err := b.makeRequest(ctx, "getFile", map[string]any{"file_id": fileID})
	if err != nil {
		return "", err
	}
	url, _ := data["download_url"].(string)
	return url, nil
}

// DownloadFile downloads a file by file_id. When `saveAs` is empty the
// raw bytes are returned. The optional `progress` callback is invoked
// after each chunk.
func (b *Bot) DownloadFile(ctx context.Context, fileID, saveAs string, asBytes bool, chunkSize int, progress func(downloaded, total int64)) (any, error) {
	url, err := b.GetFile(ctx, fileID)
	if err != nil {
		return nil, err
	}
	if url == "" {
		return nil, &APIError{Status: "InvalidResponse", DevMessage: "empty download_url"}
	}
	return downloadFromURL(ctx, b, url, saveAs, asBytes, chunkSize, progress)
}

// downloadFromURL is a small indirection so the webhook file can re-use
// the same download path.
func downloadFromURL(ctx context.Context, b *Bot, url, saveAs string, asBytes bool, chunkSize int, progress func(downloaded, total int64)) (any, error) {
	req, err := newGetRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	resp, err := b.transport.HTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, &APIError{Status: fmt.Sprintf("%d", resp.StatusCode), DevMessage: "download failed"}
	}
	if chunkSize <= 0 {
		chunkSize = 65536
	}
	total := resp.ContentLength
	if asBytes {
		// Guard against ContentLength == -1 (chunked encoding),
		// which would make make([]byte, 0, -1) panic.
		var buf []byte
		if total > 0 {
			buf = make([]byte, 0, total)
		} else {
			buf = make([]byte, 0)
		}
		chunk := make([]byte, chunkSize)
		var downloaded int64
		for {
			n, err := resp.Body.Read(chunk)
			if n > 0 {
				buf = append(buf, chunk[:n]...)
				downloaded += int64(n)
				if progress != nil {
					progress(downloaded, total)
				}
			}
			if err != nil {
				// Distinguish a clean EOF (ends the stream) from a real I/O error;
				// a non-EOF error is surfaced rather than returning
				// truncated data as success.
				if err != io.EOF {
					return nil, &APIError{Status: "DownloadFailed", DevMessage: err.Error()}
				}
				break
			}
		}
		return buf, nil
	}
	if saveAs == "" {
		return nil, &APIError{Status: "InvalidInput", DevMessage: "saveAs required when asBytes is false"}
	}
	f, err := createFile(saveAs)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	chunk := make([]byte, chunkSize)
	var downloaded int64
	for {
		n, err := resp.Body.Read(chunk)
		if n > 0 {
			if _, werr := f.Write(chunk[:n]); werr != nil {
				return nil, werr
			}
			downloaded += int64(n)
			if progress != nil {
				progress(downloaded, total)
			}
		}
		if err != nil {
			// Distinguish a clean EOF (ends the stream) from a real I/O error;
			// a non-EOF error is surfaced rather than writing
			// truncated data to disk as if successful.
			if err != io.EOF {
				return nil, &APIError{Status: "DownloadFailed", DevMessage: err.Error()}
			}
			break
		}
	}
	return saveAs, nil
}
