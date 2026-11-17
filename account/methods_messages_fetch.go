package account

import (
	"context"
	"fmt"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"strings"
)

// Fetching, paging and searching messages.

// SearchChatMessages searches messages inside a single chat. The API
// uses `{type: "Text"|"Hashtag"}` (no `from_user_id`).
func (a *Account) SearchChatMessages(ctx context.Context, objectGUID, searchText, searchType string) (*events.Event, error) {
	if searchType == "" {
		searchType = "Text"
	}
	return a.Request(ctx, RequestOptions{
		Method:  "searchChatMessages",
		Input:   map[string]any{"object_guid": objectGUID, "search_text": searchText, "type": searchType},
		Encrypt: true,
	})
}

// GetMessages fetches a window of messages from an object chat.
func (a *Account) GetMessages(ctx context.Context, objectGUID, maxID string, limit int, sort string) (*events.Event, error) {
	if sort == "" {
		sort = "FromMax"
	}
	// Build the input conditionally — the Rubika server is picky about
	// extra/empty fields and may reject a literal "" max_id or 0 limit.
	in := map[string]any{"object_guid": objectGUID, "sort": sort}
	if maxID != "" {
		in["max_id"] = maxID
	}
	if limit > 0 {
		in["limit"] = limit
	}
	return a.Request(ctx, RequestOptions{
		Method:  "getMessages",
		Input:   in,
		Encrypt: true,
	})
}

// GetMessagesByID fetches specific messages by ID.
func (a *Account) GetMessagesByID(ctx context.Context, objectGUID string, messageIDs []string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "getMessagesByID",
		Input:   map[string]any{"object_guid": objectGUID, "message_ids": messageIDs},
		Encrypt: true,
	})
}

// GetMessagesInterval fetches a range of messages around the supplied
// middle message ID. The API uses `middle_message_id`, not `mid` /
// `from_ticket_id`.
func (a *Account) GetMessagesInterval(ctx context.Context, objectGUID, middleMessageID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "getMessagesInterval",
		Input:   map[string]any{"object_guid": objectGUID, "middle_message_id": middleMessageID},
		Encrypt: true,
	})
}

// GetMessagesUpdates returns the message updates after `state` (the
// server-side cursor; omit it to start from the newest state).
func (a *Account) GetMessagesUpdates(ctx context.Context, objectGUID string, state ...string) (*events.Event, error) {
	in := map[string]any{"object_guid": objectGUID}
	if len(state) > 0 && state[0] != "" {
		in["state"] = state[0]
	}
	return a.Request(ctx, RequestOptions{Method: "getMessagesUpdates", Input: in, Encrypt: true})
}

// GetMessageShareURL returns the share URL of a message.
func (a *Account) GetMessageShareURL(ctx context.Context, objectGUID, messageID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:     "getMessageShareUrl",
		Input:      map[string]any{"object_guid": objectGUID, "message_id": messageID},
		TmpSession: true,
		Encrypt:    true,
	})
}

// GetMessageReactions returns the reactions attached to a message.
func (a *Account) GetMessageReactions(ctx context.Context, objectGUID, messageID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "getMessageReactions",
		Input:   map[string]any{"object_guid": objectGUID, "message_id": messageID},
		Encrypt: true,
	})
}

// GetChatHistory is a convenience wrapper around GetMessages that
// fetches the most recent `limit` messages from a chat, starting from
// the given `maxID` (pass "" to start from the newest). It uses
// `sort=FromMax` so messages are returned newest-first.
func (a *Account) GetChatHistory(ctx context.Context, objectGUID, maxID string, limit int) (*events.Event, error) {
	if limit <= 0 {
		limit = 50
	}
	return a.GetMessages(ctx, objectGUID, maxID, limit, "FromMax")
}

// GetUnreadMentions returns unread mentions for a chat.
func (a *Account) GetUnreadMentions(ctx context.Context, objectGUID, maxID, minID string) (*events.Event, error) {
	in := map[string]any{"object_guid": objectGUID}
	if maxID != "" {
		in["max_id"] = maxID
	}
	if minID != "" {
		in["min_id"] = minID
	}
	return a.Request(ctx, RequestOptions{Method: "getUnreadMentions", Input: in, Encrypt: true})
}

// SearchGlobalMessages searches messages globally.
func (a *Account) SearchGlobalMessages(ctx context.Context, searchText, startID string) (*events.Event, error) {
	in := map[string]any{"search_text": searchText, "type": "Text"}
	if startID != "" {
		in["start_id"] = startID
	}
	return a.Request(ctx, RequestOptions{Method: "searchGlobalMessages", Input: in, Encrypt: true})
}

// GetPinnedMessages returns the pinned messages of a chat. It works
// for groups, channels and user (private) chats.
func (a *Account) GetPinnedMessages(ctx context.Context, objectGUID string) (*events.Event, error) {
	// Try the right info method based on the GUID prefix so user
	// (u0) chats are supported too, not just groups and channels.
	var chatInfo *events.Event
	var err error
	switch {
	case strings.HasPrefix(objectGUID, "g0"):
		chatInfo, err = a.GetGroupInfo(ctx, objectGUID)
	case strings.HasPrefix(objectGUID, "c0"):
		chatInfo, err = a.GetChannelInfo(ctx, objectGUID)
	default:
		chatInfo, err = a.GetUserInfo(ctx, objectGUID)
	}
	if err != nil {
		return nil, err
	}
	pinnedIDs := chatInfo.Slice("pinned_message_ids")
	if len(pinnedIDs) == 0 {
		return events.New(map[string]any{"messages": []any{}}), nil
	}
	// The server may return pinned_message_ids as strings or as JSON
	// numbers; coerce both to string to avoid a type-assertion panic.
	var ids []string
	for _, id := range pinnedIDs {
		ids = append(ids, anyToStr(id))
	}
	return a.GetMessagesByID(ctx, objectGUID, ids)
}

// GetScheduledHistory returns scheduled messages.
func (a *Account) GetScheduledHistory(ctx context.Context, objectGUID, startID string) (*events.Event, error) {
	in := map[string]any{"object_guid": objectGUID}
	if startID != "" {
		in["start_id"] = startID
	}
	return a.Request(ctx, RequestOptions{Method: "getScheduledHistory", Input: in, Encrypt: true})
}

// GetMessageURL returns the share URL of a message. It is an alias for
// GetMessageShareURL with a clearer name.
func (a *Account) GetMessageURL(ctx context.Context, objectGUID, messageID string) (string, error) {
	resp, err := a.GetMessageShareURL(ctx, objectGUID, messageID)
	if err != nil {
		return "", err
	}
	return resp.AnyStr("share_url", "link", "url"), nil
}

// GetChatMessages fetches messages from a chat. It is an alias for
// GetMessages.
func (a *Account) GetChatMessages(ctx context.Context, objectGUID, maxID string, limit int, sort string) (*events.Event, error) {
	return a.GetMessages(ctx, objectGUID, maxID, limit, sort)
}

// GetMessagesBatch fetches messages by IDs in batches of 20.
func (a *Account) GetMessagesBatch(ctx context.Context, objectGUID string, messageIDs []string) ([]map[string]any, error) {
	var all []map[string]any
	const batchSize = 20
	for i := 0; i < len(messageIDs); i += batchSize {
		end := i + batchSize
		if end > len(messageIDs) {
			end = len(messageIDs)
		}
		batch := messageIDs[i:end]
		resp, err := a.GetMessagesByID(ctx, objectGUID, batch)
		if err != nil {
			return all, err
		}
		msgs := resp.Slice("messages")
		for _, m := range msgs {
			if mm, ok := m.(map[string]any); ok {
				all = append(all, mm)
			}
		}
	}
	return all, nil
}

// GetAllMessages fetches all messages in a chat (up to maxPages pages).
// Set maxPages to 0 for unbounded (use with caution).
func (a *Account) GetAllMessages(ctx context.Context, objectGUID string, maxPages int) ([]map[string]any, error) {
	if maxPages <= 0 {
		maxPages = 10000 // safety bound
	}
	var all []map[string]any
	middleID := ""
	for i := 0; i < maxPages; i++ {
		resp, err := a.GetMessagesInterval(ctx, objectGUID, middleID)
		if err != nil {
			return all, err
		}
		msgs := resp.Slice("messages")
		if len(msgs) == 0 {
			break
		}
		for _, m := range msgs {
			if mm, ok := m.(map[string]any); ok {
				all = append(all, mm)
			}
		}
		// Use the last message ID as the next cursor
		if last, ok := msgs[len(msgs)-1].(map[string]any); ok {
			middleID = fmt.Sprintf("%v", last["message_id"])
		}
		if middleID == "" {
			break
		}
	}
	return all, nil
}

// GetFirstMessages fetches the first batch of messages in a chat by
// calling getMessagesInterval with an empty middleMessageID.
func (a *Account) GetFirstMessages(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.GetMessagesInterval(ctx, objectGUID, "")
}

// GetLastMessages fetches the most recent messages of a chat.
func (a *Account) GetLastMessages(ctx context.Context, objectGUID string, limit int) (*events.Event, error) {
	if limit <= 0 {
		limit = 50
	}
	return a.GetMessages(ctx, objectGUID, "", limit, "FromMax")
}

// GetScheduledMessages returns the scheduled messages of a chat.
func (a *Account) GetScheduledMessages(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.GetScheduledHistory(ctx, objectGUID, "")
}
