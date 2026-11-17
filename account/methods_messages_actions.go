package account

import (
	"context"
	"fmt"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/enums"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"strconv"
	"strings"
	"time"
)

// Message actions: edit, delete, forward, pin, react, transcribe, click.

// TranscribeVoice sends a voice for transcription.
func (a *Account) TranscribeVoice(ctx context.Context, objectGUID, messageID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "transcribeVoice",
		Input:   map[string]any{"object_guid": objectGUID, "message_id": messageID},
		Encrypt: true,
	})
}

// GetTranscription fetches the result of a previous transcription.
// The API expects `{message_id, transcription_id}`.
func (a *Account) GetTranscription(ctx context.Context, messageID, transcriptionID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "getTranscription",
		Input:   map[string]any{"message_id": messageID, "transcription_id": transcriptionID},
		Encrypt: true,
	})
}

// MarkAsRead marks a chat as read by fetching the last message id and
// calling seenChats with it. The seenChats API expects a map of
// chat_guid → last_seen_message_id; passing 0 marks nothing.
func (a *Account) MarkAsRead(ctx context.Context, objectGUID string) error {
	info, err := a.GetInfo(ctx, objectGUID)
	if err != nil {
		return err
	}
	lastID := int64(0)
	if lastMsg := info.EventOf("last_message"); lastMsg != nil {
		idStr := lastMsg.AnyStr("message_id")
		if v, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			lastID = v
		}
	}
	seen := map[string]any{objectGUID: lastID}
	_, err = a.Request(ctx, RequestOptions{
		Method:  "seenChats",
		Input:   map[string]any{"seen_list": seen},
		Encrypt: true,
	})
	return err
}

// MarkAsSeen is an alias for MarkAsRead.
func (a *Account) MarkAsSeen(ctx context.Context, objectGUID string) error {
	return a.MarkAsRead(ctx, objectGUID)
}

// ClickInlineButton clicks an inline button by sending a sendMessageAPICall
// with the button's ID.
func (a *Account) ClickInlineButton(ctx context.Context, objectGUID, messageID, buttonID string) (*events.Event, error) {
	return a.SendMessageAPICall(ctx, objectGUID, messageID, buttonID)
}

// EditMessage edits the text of an existing message.
func (a *Account) EditMessage(ctx context.Context, objectGUID, messageID, msg string, parseMode enums.ParseMode) (*events.Event, error) {
	if parseMode == "" {
		parseMode = a.parseMode
	}
	input := map[string]any{
		"object_guid": objectGUID,
		"message_id":  messageID,
		"text":        strings.TrimSpace(msg),
	}
	if parseMode == enums.ParseModeHTML {
		md := a.formatter.Format(a.formatter.ToMarkdown(msg))
		mergeMetadata(input, &md)
	} else if parseMode == enums.ParseModeMarkdown {
		md := a.formatter.Format(msg)
		mergeMetadata(input, &md)
	}
	return a.Request(ctx, RequestOptions{Method: "editMessage", Input: input, Encrypt: true})
}

// DeleteMessages deletes messages. `deleteType` is "Global" (default,
// deletes for everyone) or "Local" (deletes only for the current user).
func (a *Account) DeleteMessages(ctx context.Context, objectGUID string, messageIDs []string) (*events.Event, error) {
	if len(messageIDs) == 0 {
		return nil, fmt.Errorf("no message IDs supplied")
	}
	return a.Request(ctx, RequestOptions{
		Method: "deleteMessages",
		Input: map[string]any{
			"object_guid": objectGUID,
			"message_ids": messageIDs,
			"type":        "Global", // delete for everyone (the server default)
		},
		Encrypt: true,
	})
}

// DeleteMessagesLocal deletes messages only for the current user
// (Local delete — other participants can still see them).
func (a *Account) DeleteMessagesLocal(ctx context.Context, objectGUID string, messageIDs []string) (*events.Event, error) {
	if len(messageIDs) == 0 {
		return nil, fmt.Errorf("no message IDs supplied")
	}
	return a.Request(ctx, RequestOptions{
		Method: "deleteMessages",
		Input: map[string]any{
			"object_guid": objectGUID,
			"message_ids": messageIDs,
			"type":        "Local",
		},
		Encrypt: true,
	})
}

// ForwardMessages forwards messages between chats. When `hideAuthor`
// is true, the original sender is hidden.
func (a *Account) ForwardMessages(ctx context.Context, fromObjectGUID string, messageIDs []string, toObjectGUID string, hideAuthor bool, isMute bool) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method: "forwardMessages",
		Input: map[string]any{
			"from_object_guid": fromObjectGUID,
			"message_ids":      messageIDs,
			"to_object_guid":   toObjectGUID,
			"rnd":              randInt63(),
			"drop_author":      hideAuthor,
			"is_mute":          isMute,
		},
		Encrypt: true,
	})
}

// SetPinMessage pins a message (action = "Pin").
func (a *Account) SetPinMessage(ctx context.Context, objectGUID, messageID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "setPinMessage",
		Input:   map[string]any{"object_guid": objectGUID, "message_id": messageID, "action": "Pin"},
		Encrypt: true,
	})
}

// SetUnpinMessage unpins a message (action = "Unpin").
func (a *Account) SetUnpinMessage(ctx context.Context, objectGUID, messageID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "setPinMessage",
		Input:   map[string]any{"object_guid": objectGUID, "message_id": messageID, "action": "Unpin"},
		Encrypt: true,
	})
}

// UnpinAllMessages unpins all messages in a chat.
func (a *Account) UnpinAllMessages(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "setPinMessage",
		Input:   map[string]any{"object_guid": objectGUID, "action": "UnpinAllMessages"},
		Encrypt: true,
	})
}

// AutoDeleteMessage schedules a message for deletion after `delay`.
func (a *Account) AutoDeleteMessage(ctx context.Context, objectGUID, messageID string, delay time.Duration) error {
	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return ctx.Err()
	}
	_, err := a.DeleteMessages(ctx, objectGUID, []string{messageID})
	return err
}

// AddReaction adds a reaction to a message. The API uses
// `actionOnMessageReaction` with `action: "Add"` and an int `reaction_id`.
func (a *Account) AddReaction(ctx context.Context, objectGUID, messageID string, reactionID int) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method: "actionOnMessageReaction",
		Input: map[string]any{
			"object_guid": objectGUID,
			"message_id":  messageID,
			"action":      "Add",
			"reaction_id": reactionID,
		},
		Encrypt: true,
	})
}

// RemoveReaction removes a reaction from a message.
func (a *Account) RemoveReaction(ctx context.Context, objectGUID, messageID string, reactionID int) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method: "actionOnMessageReaction",
		Input: map[string]any{
			"object_guid": objectGUID,
			"message_id":  messageID,
			"action":      "Remove",
			"reaction_id": reactionID,
		},
		Encrypt: true,
	})
}

// ClickMessageURL clicks a link inside a message.
func (a *Account) ClickMessageURL(ctx context.Context, objectGUID, messageID, linkURL string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "clickMessageUrl", Input: map[string]any{
		"object_guid": objectGUID,
		"message_id":  messageID,
		"link_url":    linkURL,
	}, Encrypt: true})
}

// ClickKeyPadButton clicks a chat keypad button by its ID.
func (a *Account) ClickKeyPadButton(ctx context.Context, objectGUID, messageID, buttonID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "sendMessage", Input: map[string]any{
		"object_guid": objectGUID,
		"message_id":  messageID,
		"aux_data":    map[string]any{"button_id": buttonID},
		"rnd":         randInt63(),
	}, Encrypt: true})
}

// SendMessageAPICall clicks an inline button on the user-side by
// sending a message with the button's aux_data.
func (a *Account) SendMessageAPICall(ctx context.Context, objectGUID, messageID, buttonID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "sendMessageAPICall", Input: map[string]any{
		"object_guid": objectGUID,
		"message_id":  messageID,
		"aux_data":    map[string]any{"button_id": buttonID},
	}, Encrypt: true})
}

// FeedbackTranscription sends feedback on a voice transcription.
func (a *Account) FeedbackTranscription(ctx context.Context, messageID, transcriptionID, feedbackType, feedbackText string) (*events.Event, error) {
	in := map[string]any{
		"message_id":       messageID,
		"transcription_id": transcriptionID,
		"feedback_type":    feedbackType, // "OK" or "Incorrect"
	}
	if feedbackText != "" {
		in["feedback_text"] = feedbackText
	}
	return a.Request(ctx, RequestOptions{Method: "feedbackTranscription", Input: in, Encrypt: true})
}

// SendNowScheduledMessage sends a scheduled message immediately.
func (a *Account) SendNowScheduledMessage(ctx context.Context, objectGUID, messageID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "sendNowScheduledMessage", Input: map[string]any{
		"object_guid": objectGUID,
		"message_id":  messageID,
	}, Encrypt: true})
}
