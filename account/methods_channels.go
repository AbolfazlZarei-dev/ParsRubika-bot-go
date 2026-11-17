package account

import (
	"context"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// Channel management: info, members, admins, links and statistics.

// GetChannelInfo fetches information about a channel.
func (a *Account) GetChannelInfo(ctx context.Context, channelGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getChannelInfo", Input: map[string]any{"channel_guid": channelGUID}, Encrypt: true})
}

// AddChannel creates a new channel owned by the current user. The
// optional channelType ("Public" in maxrubika's default) is sent only
// when provided, so existing calls keep their exact old payload.
func (a *Account) AddChannel(ctx context.Context, title, description, username string, channelType ...string) (*events.Event, error) {
	in := map[string]any{"title": title, "description": description}
	if username != "" {
		in["username"] = username
	}
	if len(channelType) > 0 && channelType[0] != "" {
		in["channel_type"] = channelType[0]
	}
	return a.Request(ctx, RequestOptions{Method: "addChannel", Input: in, Encrypt: true})
}

// RemoveChannel removes (deletes) a channel.
func (a *Account) RemoveChannel(ctx context.Context, channelGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "removeChannel", Input: map[string]any{"channel_guid": channelGUID}, Encrypt: true})
}

// EditChannelInfo edits a channel's metadata. Pass an empty string for
// any field you do not want to change. The `updated_parameters` list
// is built automatically.
func (a *Account) EditChannelInfo(ctx context.Context, channelGUID, title, description, username string) (*events.Event, error) {
	in := map[string]any{"channel_guid": channelGUID}
	updated := []string{}
	if title != "" {
		in["title"] = title
		updated = append(updated, "title")
	}
	if description != "" {
		in["description"] = description
		updated = append(updated, "description")
	}
	if username != "" {
		in["username"] = username
		updated = append(updated, "username")
	}
	in["updated_parameters"] = updated
	return a.Request(ctx, RequestOptions{Method: "editChannelInfo", Input: in, Encrypt: true})
}

// CheckChannelUsername checks whether a username is free for a channel.
func (a *Account) CheckChannelUsername(ctx context.Context, username string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:     "checkChannelUsername",
		Input:      map[string]any{"username": username},
		TmpSession: true,
		Encrypt:    true,
	})
}

// UpdateChannelUsername changes the channel's username.
func (a *Account) UpdateChannelUsername(ctx context.Context, channelGUID, username string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "updateChannelUsername",
		Input:   map[string]any{"channel_guid": channelGUID, "username": username},
		Encrypt: true,
	})
}

// AddChannelMembers adds members to a channel. The API expects
// `member_guids` (not `members`).
func (a *Account) AddChannelMembers(ctx context.Context, channelGUID string, userGUIDs []string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "addChannelMembers",
		Input:   map[string]any{"channel_guid": channelGUID, "member_guids": userGUIDs},
		Encrypt: true,
	})
}

// GetChannelAllMembers returns the channel member list with an
// optional search-text filter and pagination start-id.
func (a *Account) GetChannelAllMembers(ctx context.Context, channelGUID, startID, searchText string, limit int) (*events.Event, error) {
	in := map[string]any{"channel_guid": channelGUID}
	if startID != "" {
		in["start_id"] = startID
	}
	if searchText != "" {
		in["search_text"] = searchText
	}
	if limit > 0 {
		in["limit"] = limit
	}
	return a.Request(ctx, RequestOptions{Method: "getChannelAllMembers", Input: in, Encrypt: true})
}

// GetChannelAdminMembers returns the admin member list (paged).
func (a *Account) GetChannelAdminMembers(ctx context.Context, channelGUID, startID string) (*events.Event, error) {
	in := map[string]any{"channel_guid": channelGUID}
	if startID != "" {
		in["start_id"] = startID
	}
	return a.Request(ctx, RequestOptions{Method: "getChannelAdminMembers", Input: in, Encrypt: true})
}

// GetChannelAdminAccessList returns the access rights of an admin.
// `memberGUID` is required.
func (a *Account) GetChannelAdminAccessList(ctx context.Context, channelGUID, memberGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "getChannelAdminAccessList",
		Input:   map[string]any{"channel_guid": channelGUID, "member_guid": memberGUID},
		Encrypt: true,
	})
}

// GetChannelLink returns the invite link of a channel.
func (a *Account) GetChannelLink(ctx context.Context, channelGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getChannelLink", Input: map[string]any{"channel_guid": channelGUID}, Encrypt: true})
}

// SetChannelLink generates a new invite link.
func (a *Account) SetChannelLink(ctx context.Context, channelGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "setChannelLink", Input: map[string]any{"channel_guid": channelGUID}, Encrypt: true})
}

// BanChannelMember bans a member from a channel (action = "Set").
func (a *Account) BanChannelMember(ctx context.Context, channelGUID, userGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "banChannelMember",
		Input:   map[string]any{"channel_guid": channelGUID, "member_guid": userGUID, "action": "Set"},
		Encrypt: true,
	})
}

// UnbanChannelMember unbans a channel member (action = "Unset").
func (a *Account) UnbanChannelMember(ctx context.Context, channelGUID, userGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "banChannelMember",
		Input:   map[string]any{"channel_guid": channelGUID, "member_guid": userGUID, "action": "Unset"},
		Encrypt: true,
	})
}

// ChannelPreviewByJoinLink fetches the channel preview before joining.
// The link is normalised to its trailing slug.
func (a *Account) ChannelPreviewByJoinLink(ctx context.Context, link string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:     "channelPreviewByJoinLink",
		Input:      map[string]any{"hash_link": trailingSlug(link)},
		TmpSession: true,
		Encrypt:    true,
	})
}

// JoinChannelByLink joins a channel via its invite link. The API only
// needs the `hash_link` — the channel GUID is implicit.
func (a *Account) JoinChannelByLink(ctx context.Context, link string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "joinChannelByLink",
		Input:   map[string]any{"hash_link": trailingSlug(link)},
		Encrypt: true,
	})
}

// JoinChannelAction accepts/declines a join request. `action` is
// "Join" or "Leave".
func (a *Account) JoinChannelAction(ctx context.Context, channelGUID, action string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "joinChannelAction",
		Input:   map[string]any{"channel_guid": channelGUID, "action": action},
		Encrypt: true,
	})
}

// SeenChannelMessages marks a range of channel messages as seen.
// `minID` and `maxID` define the inclusive range.
func (a *Account) SeenChannelMessages(ctx context.Context, channelGUID, minID, maxID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "seenChannelMessages",
		Input:   map[string]any{"channel_guid": channelGUID, "min_id": minID, "max_id": maxID},
		Encrypt: true,
	})
}

// GetChannelSeenCount returns the seen count for a range of messages.
func (a *Account) GetChannelSeenCount(ctx context.Context, channelGUID, minID, maxID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getChannelSeenCount", Input: map[string]any{
		"channel_guid": channelGUID,
		"min_id":       minID,
		"max_id":       maxID,
	}, Encrypt: true})
}

// GetChannelStatistics returns channel-specific statistics. The
// Rubika API exposes a dedicated `channelStatistics` method for this
// (distinct from the account-wide `getAllStatistics`), so we call it
// directly to avoid returning global stats when the caller asked for
// channel-only data.
func (a *Account) GetChannelStatistics(ctx context.Context, channelGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "channelStatistics", Input: map[string]any{
		"channel_guid": channelGUID,
	}, Encrypt: true})
}

// SetChannelAdmin promotes or demotes a channel admin.
func (a *Account) SetChannelAdmin(ctx context.Context, channelGUID, memberGUID, action string, accessList []string) (*events.Event, error) {
	in := map[string]any{
		"channel_guid": channelGUID,
		"member_guid":  memberGUID,
		"action":       action, // "SetAdmin" or "UnsetAdmin"
	}
	if len(accessList) > 0 {
		in["access_list"] = accessList
	}
	return a.Request(ctx, RequestOptions{Method: "setChannelAdmin", Input: in, Encrypt: true})
}

// PromoteChannelAdmin promotes a channel admin.
func (a *Account) PromoteChannelAdmin(ctx context.Context, channelGUID, memberGUID string, accessList []string) (*events.Event, error) {
	return a.SetChannelAdmin(ctx, channelGUID, memberGUID, "SetAdmin", accessList)
}

// DemoteChannelAdmin demotes a channel admin.
func (a *Account) DemoteChannelAdmin(ctx context.Context, channelGUID, memberGUID string) (*events.Event, error) {
	return a.SetChannelAdmin(ctx, channelGUID, memberGUID, "UnsetAdmin", nil)
}

// SetChannelType sets the channel type (Public or Private).
func (a *Account) SetChannelType(ctx context.Context, channelGUID, channelType string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "editChannelInfo", Input: map[string]any{
		"channel_guid":       channelGUID,
		"channel_type":       channelType,
		"updated_parameters": []string{"channel_type"},
	}, Encrypt: true})
}

// EditChannelTitle is a convenience wrapper for editing just the title.
func (a *Account) EditChannelTitle(ctx context.Context, channelGUID, title string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "editChannelInfo", Input: map[string]any{
		"channel_guid":       channelGUID,
		"title":              title,
		"updated_parameters": []string{"title"},
	}, Encrypt: true})
}
