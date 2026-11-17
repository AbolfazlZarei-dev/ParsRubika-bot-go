package account

import (
	"context"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/enums"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"strings"
	"time"
)

// Chat list, chat settings, drafts, folders and join links.

// GetChats returns the user's chats.
func (a *Account) GetChats(ctx context.Context, startID string) (*events.Event, error) {
	in := map[string]any{}
	if startID != "" {
		in["start_id"] = startID
	}
	return a.Request(ctx, RequestOptions{Method: "getChats", Input: in, Encrypt: true})
}

// GetChatsUpdates fetches the latest chat updates.
//
// The server requires a `state` parameter (Unix-seconds timestamp,
// typically `now - 150`). Without it the server replies
// `INVALID_INPUT: getChatsUpdates failed`, so the call always sends
// `state = now - 150`.
func (a *Account) GetChatsUpdates(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "getChatsUpdates",
		Input:   map[string]any{"state": time.Now().Unix() - 150},
		Encrypt: true,
	})
}

// SeenChats marks the given chat GUIDs as seen.
func (a *Account) SeenChats(ctx context.Context, chatIDs []string) (*events.Event, error) {
	if len(chatIDs) == 0 {
		return nil, nil
	}
	seen := map[string]any{}
	for _, id := range chatIDs {
		seen[id] = int64(0)
	}
	return a.Request(ctx, RequestOptions{Method: "seenChats", Input: map[string]any{"seen_list": seen}, Encrypt: true})
}

// GetAbsObjects returns the absolute object info for the given GUIDs.
func (a *Account) GetAbsObjects(ctx context.Context, guids []string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getAbsObjects", Input: map[string]any{"object_guids": guids}, Encrypt: true})
}

// SendChatActivity reports a chat activity. The API expects
// `{activity: "Typing"}` (a string), not `{activity: {type: "Typing"}}`.
func (a *Account) SendChatActivity(ctx context.Context, objectGUID string, action enums.ChatAction) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "sendChatActivity",
		Input:   map[string]any{"object_guid": objectGUID, "activity": string(action)},
		Encrypt: true,
	})
}

// SetActionChat mutes / unmutes / archives / unarchives / pins / unpins
// a chat. `action` is one of "Mute", "Unmute", "Archive", "UnArchive",
// "Pin", "Unpin".
func (a *Account) SetActionChat(ctx context.Context, objectGUID, action string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "setActionChat",
		Input:   map[string]any{"object_guid": objectGUID, "action": action},
		Encrypt: true,
	})
}

// MuteChat mutes the supplied chat.
func (a *Account) MuteChat(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.SetActionChat(ctx, objectGUID, "Mute")
}

// UnmuteChat unmutes the supplied chat.
func (a *Account) UnmuteChat(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.SetActionChat(ctx, objectGUID, "Unmute")
}

// ArchiveChat archives the supplied chat.
func (a *Account) ArchiveChat(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.SetActionChat(ctx, objectGUID, "Archive")
}

// UnarchiveChat unarchives the supplied chat.
func (a *Account) UnarchiveChat(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.SetActionChat(ctx, objectGUID, "UnArchive")
}

// PinChat pins the supplied chat to the top of the chat list.
func (a *Account) PinChat(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.SetActionChat(ctx, objectGUID, "Pin")
}

// UnpinChat unpins the supplied chat.
func (a *Account) UnpinChat(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.SetActionChat(ctx, objectGUID, "Unpin")
}

// DeleteChatHistory clears the chat history for the supplied GUID.
func (a *Account) DeleteChatHistory(ctx context.Context, objectGUID, lastMessageID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "deleteChatHistory",
		Input:   map[string]any{"object_guid": objectGUID, "last_message_id": lastMessageID},
		Encrypt: true,
	})
}

// GetLinkFromAppURL converts an `app.rubika.ir/...` URL into a deep link.
func (a *Account) GetLinkFromAppURL(ctx context.Context, appURL string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:     "getLinkFromAppUrl",
		Input:      map[string]any{"app_url": appURL},
		TmpSession: true,
		Encrypt:    true,
	})
}

// GetJoinLinks returns the join links of a chat.
func (a *Account) GetJoinLinks(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getJoinLinks", Input: map[string]any{"object_guid": objectGUID}, Encrypt: true})
}

// CreateJoinLink creates a new join link for a chat. The API uses
// `usage_limit` (not `uses_limit`) and a `request_needed` boolean.
func (a *Account) CreateJoinLink(ctx context.Context, objectGUID, title string, expirationSec, usageLimit int, requestNeeded bool) (*events.Event, error) {
	in := map[string]any{
		"object_guid":    objectGUID,
		"request_needed": requestNeeded,
		"usage_limit":    usageLimit,
	}
	if title != "" {
		in["title"] = title
	}
	if expirationSec > 0 {
		in["expire_time"] = expirationSec
	}
	return a.Request(ctx, RequestOptions{Method: "createJoinLink", Input: in, Encrypt: true})
}

// GetJoinRequests returns the pending join requests.
func (a *Account) GetJoinRequests(ctx context.Context, objectGUID, cursor string) (*events.Event, error) {
	in := map[string]any{"object_guid": objectGUID}
	if cursor != "" {
		in["cursor"] = cursor
	}
	return a.Request(ctx, RequestOptions{Method: "getJoinRequests", Input: in, Encrypt: true})
}

// ActionOnJoinRequest accepts or declines a join request. The API uses
// `{object_guid, object_type, user_guid, action}` where `object_type`
// is derived from the GUID prefix ("Group" / "Channel") and `action`
// is "Accept" or "Reject".
func (a *Account) ActionOnJoinRequest(ctx context.Context, objectGUID, userGUID, action string) (*events.Event, error) {
	objType := "Group"
	if strings.HasPrefix(objectGUID, "c0") {
		objType = "Channel"
	}
	return a.Request(ctx, RequestOptions{
		Method: "actionOnJoinRequest",
		Input: map[string]any{
			"object_guid": objectGUID,
			"object_type": objType,
			"user_guid":   userGUID,
			"action":      action,
		},
		Encrypt: true,
	})
}

// AcceptJoinRequest accepts a join request.
func (a *Account) AcceptJoinRequest(ctx context.Context, objectGUID, userGUID string) (*events.Event, error) {
	return a.ActionOnJoinRequest(ctx, objectGUID, userGUID, "Accept")
}

// RejectJoinRequest rejects a join request.
func (a *Account) RejectJoinRequest(ctx context.Context, objectGUID, userGUID string) (*events.Event, error) {
	return a.ActionOnJoinRequest(ctx, objectGUID, userGUID, "Reject")
}

// SetChatProtectedContent toggles the "protected content" (restricted
// content / no-forward) flag on a group or channel.
func (a *Account) SetChatProtectedContent(ctx context.Context, objectGUID string, enabled bool) (*events.Event, error) {
	if strings.HasPrefix(objectGUID, "c0") {
		// Channel
		return a.Request(ctx, RequestOptions{Method: "editChannelInfo", Input: map[string]any{
			"channel_guid":          objectGUID,
			"is_restricted_content": enabled,
			"updated_parameters":    []string{"is_restricted_content"},
		}, Encrypt: true})
	}
	// Group
	return a.Request(ctx, RequestOptions{Method: "editGroupInfo", Input: map[string]any{
		"group_guid":            objectGUID,
		"is_restricted_content": enabled,
		"updated_parameters":    []string{"is_restricted_content"},
	}, Encrypt: true})
}

// MuteChatFor mutes a chat for the given duration (in seconds). The
// Rubika API requires a minimum of 1800 seconds (30 minutes).
func (a *Account) MuteChatFor(ctx context.Context, objectGUID string, durationSeconds int) (*events.Event, error) {
	if durationSeconds < 1800 {
		durationSeconds = 1800
	}
	return a.Request(ctx, RequestOptions{Method: "setActionChat", Input: map[string]any{
		"object_guid": objectGUID,
		"action":      "Mute",
		"mute":        durationSeconds,
	}, Encrypt: true})
}

// DeleteNoAccessGroupChat deletes a group chat that the user cannot access.
func (a *Account) DeleteNoAccessGroupChat(ctx context.Context, groupGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "deleteNoAccessGroupChat", Input: map[string]any{"group_guid": groupGUID}, Encrypt: true})
}

// trailingSlug returns the trailing path segment of a URL or the input
// itself when it does not contain a slash. Used to normalise join
// links before sending them as `hash_link`.
func trailingSlug(link string) string {
	if link == "" {
		return ""
	}
	if i := strings.LastIndex(link, "/"); i >= 0 {
		return link[i+1:]
	}
	return link
}

// SearchChats searches for messages within a specific chat. `query`
// is the search term; `limit` caps the number of results (use 0 for
// the server default).
func (a *Account) SearchChats(ctx context.Context, objectGUID, query string, limit int) (*events.Event, error) {
	in := map[string]any{
		"object_guid": objectGUID,
		"search_text": query,
	}
	if limit > 0 {
		in["limit"] = limit
	}
	return a.Request(ctx, RequestOptions{Method: "searchInChat", Input: in, Encrypt: true})
}

// GetChatFolder returns a single chat folder by its ID. This is the
// per-folder counterpart of GetFolders.
func (a *Account) GetChatFolder(ctx context.Context, folderID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getFolder", Input: map[string]any{
		"folder_id": folderID,
	}, Encrypt: true})
}

// GetChatsByID returns info for multiple chat GUIDs at once.
func (a *Account) GetChatsByID(ctx context.Context, chatIDs []string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getChatsByID", Input: map[string]any{
		"object_guids": chatIDs,
	}, Encrypt: true})
}

// ClearDraft clears the draft of a single chat.
func (a *Account) ClearDraft(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "clearDrafts", Input: map[string]any{
		"action":      "Selected",
		"object_guid": objectGUID,
	}, Encrypt: true})
}

// ClearAllDrafts clears all chat drafts.
func (a *Account) ClearAllDrafts(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "clearDrafts", Input: map[string]any{
		"action": "All",
	}, Encrypt: true})
}

// GetAllDrafts returns all chat drafts.
func (a *Account) GetAllDrafts(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getAllDrafts", Input: map[string]any{}, Encrypt: true})
}

// SetChatUseTime sets the chat use-time hint (anti-rate-limit).
func (a *Account) SetChatUseTime(ctx context.Context, objectGUID string, useTime int) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "setChatUseTime", Input: map[string]any{
		"object_guid": objectGUID,
		"time":        useTime,
	}, Encrypt: true})
}

// SetAutoDelete sets the auto-delete timer for chats. `autoDelete` is in
// seconds (0 to disable).
func (a *Account) SetAutoDelete(ctx context.Context, objectGUIDs []string, autoDelete int) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "setAutoDelete", Input: map[string]any{
		"object_guids": objectGUIDs,
		"auto_delete":  autoDelete,
	}, Encrypt: true})
}

// GetChatAds returns chat advertisements.
func (a *Account) GetChatAds(ctx context.Context, state int) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getChatAds", Input: map[string]any{
		"state": state,
	}, Encrypt: true})
}

// ActionOnChatAds performs an action on a chat ad (View/Click/DontShow/Report).
func (a *Account) ActionOnChatAds(ctx context.Context, action, chatAdsID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "actionOnChatAds", Input: map[string]any{
		"action":      action,
		"chat_ads_id": chatAdsID,
	}, Encrypt: true})
}

// SetAskSpamAction sets the ask-spam action for a chat.
func (a *Account) SetAskSpamAction(ctx context.Context, objectGUID, action string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "setAskSpamAction", Input: map[string]any{
		"object_guid": objectGUID,
		"action":      action, // "AddToContact", "BlockUser", "Cancel" for user; "ReportAndLeave", "Cancel" for group/channel
	}, Encrypt: true})
}

// PinChatInFolder pins a chat inside a folder.
func (a *Account) PinChatInFolder(ctx context.Context, objectGUID, folderID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "setPinChatInFolder", Input: map[string]any{
		"action":      "Pin",
		"folder_id":   folderID,
		"object_guid": objectGUID,
	}, Encrypt: true})
}

// UnpinChatInFolder unpins a chat from a folder.
func (a *Account) UnpinChatInFolder(ctx context.Context, objectGUID, folderID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "setPinChatInFolder", Input: map[string]any{
		"action":      "Unpin",
		"folder_id":   folderID,
		"object_guid": objectGUID,
	}, Encrypt: true})
}

// EditJoinLink edits an existing join link.
func (a *Account) EditJoinLink(ctx context.Context, objectGUID, hashLink, title string, expirationSec, usageLimit int, requestNeeded *bool) (*events.Event, error) {
	in := map[string]any{
		"object_guid": objectGUID,
		"join_link":   trailingSlug(hashLink),
	}
	updated := []string{}
	if title != "" {
		in["title"] = title
		updated = append(updated, "title")
	}
	if expirationSec > 0 {
		in["expire_time"] = expirationSec
		updated = append(updated, "expire_time")
	}
	if usageLimit > 0 {
		in["usage_limit"] = usageLimit
		updated = append(updated, "usage_limit")
	}
	if requestNeeded != nil {
		in["request_needed"] = *requestNeeded
		updated = append(updated, "request_needed")
	}
	in["updated_parameters"] = updated
	return a.Request(ctx, RequestOptions{Method: "editJoinLink", Input: in, Encrypt: true})
}

// RevokeJoinLink revokes a join link.
func (a *Account) RevokeJoinLink(ctx context.Context, objectGUID, hashLink string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "revokeJoinLink", Input: map[string]any{
		"object_guid": objectGUID,
		"join_link":   trailingSlug(hashLink),
	}, Encrypt: true})
}

// GetFolders returns the user's chat folders. `lastState` is the
// server-time baseline (defaults to 0).
func (a *Account) GetFolders(ctx context.Context, lastState int) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getFolders", Input: map[string]any{"last_state": lastState}, Encrypt: true})
}

// GetSuggestedFolders returns the suggested folders.
func (a *Account) GetSuggestedFolders(ctx context.Context) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getSuggestedFolders", Input: map[string]any{}, Encrypt: true})
}

// DeleteFolder deletes a chat folder.
func (a *Account) DeleteFolder(ctx context.Context, folderID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "deleteFolder", Input: map[string]any{"folder_id": folderID}, Encrypt: true})
}

// AddFolder creates a new chat folder.
func (a *Account) AddFolder(ctx context.Context, name string, includeChats, excludeChats, includeChatTypes, excludeChatTypes []string, isAddToTop bool) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method: "addFolder",
		Input: map[string]any{
			"name":                 name,
			"is_add_to_top":        isAddToTop,
			"include_object_guids": includeChats,
			"exclude_object_guids": excludeChats,
			"include_chat_types":   includeChatTypes,
			"exclude_chat_types":   excludeChatTypes,
			"updated_parameters":   []string{"name", "include_chat_types", "exclude_chat_types"},
		},
		Encrypt: true,
	})
}

// ReorderFolder reorders the user's chat folders.
func (a *Account) ReorderFolder(ctx context.Context, folderIDs []string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "reorderFolder", Input: map[string]any{"folder_ids": folderIDs}, Encrypt: true})
}
