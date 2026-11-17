package account

import (
	"context"
	"fmt"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"strings"
)

// Deprecated aliases kept for backward compatibility; prefer the canonical methods.

// GetChatInfo is a dispatcher that fetches chat info by GUID type.
// It is an alias for GetInfo with a clearer name.
func (a *Account) GetChatInfo(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.GetInfo(ctx, objectGUID)
}

// GetInfoByLink resolves a join link to a chat preview. It dispatches
// to channelPreviewByJoinLink or groupPreviewByJoinLink based on the
// URL path.
func (a *Account) GetInfoByLink(ctx context.Context, link string) (*events.Event, error) {
	if strings.Contains(link, "/joing/") {
		return a.GroupPreviewByJoinLink(ctx, link)
	}
	if strings.Contains(link, "/joinc/") {
		return a.ChannelPreviewByJoinLink(ctx, link)
	}
	// Default to group preview
	return a.GroupPreviewByJoinLink(ctx, link)
}

// GetInfoByUsername resolves a @username to an object.
func (a *Account) GetInfoByUsername(ctx context.Context, username string) (*events.Event, error) {
	username = strings.TrimPrefix(username, "@")
	return a.GetObjectByUsername(ctx, username)
}

// JoinChat joins a group or channel by link or GUID. It is a dispatcher
// alias for Join.
func (a *Account) JoinChat(ctx context.Context, input string) (*events.Event, error) {
	return a.Join(ctx, input)
}

// LeaveChatByGUID leaves a chat by GUID. It is a dispatcher alias for
// LeaveChat.
func (a *Account) LeaveChatByGUID(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.LeaveChat(ctx, objectGUID)
}

// JoinChannelByGUID joins a channel by GUID (action=Join).
func (a *Account) JoinChannelByGUID(ctx context.Context, channelGUID string) (*events.Event, error) {
	return a.JoinChannelAction(ctx, channelGUID, "Join")
}

// LeaveChannel leaves a channel (action=Leave).
func (a *Account) LeaveChannel(ctx context.Context, channelGUID string) (*events.Event, error) {
	return a.JoinChannelAction(ctx, channelGUID, "Leave")
}

// ArchiveChannel archives a channel (action=Archive).
func (a *Account) ArchiveChannel(ctx context.Context, channelGUID string) (*events.Event, error) {
	return a.JoinChannelAction(ctx, channelGUID, "Archive")
}

// GetChannelAdmins is an alias for GetChannelAdminMembers.
func (a *Account) GetChannelAdmins(ctx context.Context, channelGUID, startID string) (*events.Event, error) {
	return a.GetChannelAdminMembers(ctx, channelGUID, startID)
}

// GetChannelMembers is an alias for GetChannelAllMembers.
func (a *Account) GetChannelMembers(ctx context.Context, channelGUID, startID, searchText string, limit int) (*events.Event, error) {
	return a.GetChannelAllMembers(ctx, channelGUID, startID, searchText, limit)
}

// GetGroupAdmins is an alias for GetGroupAdminMembers.
func (a *Account) GetGroupAdmins(ctx context.Context, groupGUID, startID string) (*events.Event, error) {
	return a.GetGroupAdminMembers(ctx, groupGUID, startID)
}

// GetGroupMembers is an alias for GetGroupAllMembers.
func (a *Account) GetGroupMembers(ctx context.Context, groupGUID, startID, searchText string) (*events.Event, error) {
	return a.GetGroupAllMembers(ctx, groupGUID, startID, searchText)
}

// GetBannedMembers fetches banned members from a group or channel.
func (a *Account) GetBannedMembers(ctx context.Context, objectGUID, startID string) (*events.Event, error) {
	if strings.HasPrefix(objectGUID, "c0") {
		// Channels use getBannedChannelMembers in some API versions;
		// the standard getBannedGroupMembers works for both.
		return a.GetBannedGroupMembers(ctx, objectGUID, startID)
	}
	return a.GetBannedGroupMembers(ctx, objectGUID, startID)
}

// SetAdmin promotes a user to admin in a group or channel (dispatcher).
func (a *Account) SetAdmin(ctx context.Context, objectGUID, memberGUID string, accessList []string) (*events.Event, error) {
	if strings.HasPrefix(objectGUID, "c0") {
		return a.SetChannelAdmin(ctx, objectGUID, memberGUID, "SetAdmin", accessList)
	}
	return a.SetGroupAdmin(ctx, objectGUID, memberGUID, "SetAdmin", accessList)
}

// UnsetAdmin demotes an admin in a group or channel (dispatcher).
func (a *Account) UnsetAdmin(ctx context.Context, objectGUID, memberGUID string) (*events.Event, error) {
	if strings.HasPrefix(objectGUID, "c0") {
		return a.SetChannelAdmin(ctx, objectGUID, memberGUID, "UnsetAdmin", nil)
	}
	return a.SetGroupAdmin(ctx, objectGUID, memberGUID, "UnsetAdmin", nil)
}

// GetAdminAccessList fetches the admin access list for a member.
func (a *Account) GetAdminAccessList(ctx context.Context, objectGUID, memberGUID string) (*events.Event, error) {
	if strings.HasPrefix(objectGUID, "c0") {
		return a.GetChannelAdminAccessList(ctx, objectGUID, memberGUID)
	}
	return a.GetGroupAdminAccessList(ctx, objectGUID, memberGUID)
}

// ChangeOwner requests ownership transfer of a group or channel.
// This is an alias for RequestChangeOwner.
func (a *Account) ChangeOwner(ctx context.Context, objectGUID, newOwnerGUID string) (*events.Event, error) {
	return a.RequestChangeOwner(ctx, objectGUID, newOwnerGUID)
}

// CancelChangeOwner cancels a pending ownership transfer.
// This is an alias for the existing method.
func (a *Account) CancelChangeOwnerOp(ctx context.Context, objectGUID string) (*events.Event, error) {
	return a.CancelChangeOwner(ctx, objectGUID)
}

// SearchGlobalChats searches all chats globally. Alias for
// SearchGlobalObjects.
func (a *Account) SearchGlobalChats(ctx context.Context, query string) (*events.Event, error) {
	return a.SearchGlobalObjects(ctx, query)
}

// SearchGlobalUsers searches globally for users.
func (a *Account) SearchGlobalUsers(ctx context.Context, query string) (*events.Event, error) {
	resp, err := a.SearchGlobalObjects(ctx, query)
	if err != nil {
		return nil, err
	}
	// Filter to only users
	users := resp.Slice("users")
	out := events.New(map[string]any{"users": users})
	return out, nil
}

// SearchGlobalChannels searches globally for channels.
func (a *Account) SearchGlobalChannels(ctx context.Context, query string) (*events.Event, error) {
	resp, err := a.SearchGlobalObjects(ctx, query)
	if err != nil {
		return nil, err
	}
	channels := resp.Slice("channels")
	out := events.New(map[string]any{"channels": channels})
	return out, nil
}

// SearchGlobalGroups searches globally for groups.
func (a *Account) SearchGlobalGroups(ctx context.Context, query string) (*events.Event, error) {
	resp, err := a.SearchGlobalObjects(ctx, query)
	if err != nil {
		return nil, err
	}
	groups := resp.Slice("groups")
	out := events.New(map[string]any{"groups": groups})
	return out, nil
}

// StartBotByGUID starts a bot by sending /start. Alias for StartBot.
func (a *Account) StartBotByGUID(ctx context.Context, botGUID string) (*events.Event, error) {
	return a.StartBot(ctx, botGUID)
}

// StopBotByGUID stops a bot. Alias for StopBot.
func (a *Account) StopBotByGUID(ctx context.Context, botGUID string) (*events.Event, error) {
	return a.StopBot(ctx, botGUID)
}

// GetBotInfoByGUID fetches bot info. Alias for GetBotInfo.
func (a *Account) GetBotInfoByGUID(ctx context.Context, botGUID string) (*events.Event, error) {
	return a.GetBotInfo(ctx, botGUID)
}

// GetInlineBotsList fetches info for multiple inline bots.
func (a *Account) GetInlineBotsList(ctx context.Context, botGUIDs []string) (*events.Event, error) {
	return a.GetInlineBots(ctx, botGUIDs)
}

// GetServiceInfoByGUID fetches service chat info.
func (a *Account) GetServiceInfoByGUID(ctx context.Context, serviceGUID string) (*events.Event, error) {
	return a.GetServiceInfo(ctx, serviceGUID)
}

// DeleteServiceChatByGUID deletes a service chat.
func (a *Account) DeleteServiceChatByGUID(ctx context.Context, serviceGUID, lastMessageID string) (*events.Event, error) {
	return a.DeleteServiceChat(ctx, serviceGUID, lastMessageID)
}

// DeleteBotChatByGUID deletes a bot chat.
func (a *Account) DeleteBotChatByGUID(ctx context.Context, botGUID, lastMessageID string) (*events.Event, error) {
	return a.DeleteBotChat(ctx, botGUID, lastMessageID)
}

// GetInfo returns generic info for the supplied object — it dispatches
// to getUserInfo / getGroupInfo / getChannelInfo / getBotInfo /
// getServiceInfo based on the GUID prefix, or `getObjectByUsername`
// when the input is a username. We check prefixes directly (rather than
// relying on events.GUIDType) because GUIDType returns "User" for
// anything that is not a group or channel — including bots (b0) and
// service chats (s0), which would be mis-dispatched to getUserInfo.
func (a *Account) GetInfo(ctx context.Context, objectGUID string) (*events.Event, error) {
	switch {
	case strings.HasPrefix(objectGUID, "c0"):
		return a.GetChannelInfo(ctx, objectGUID)
	case strings.HasPrefix(objectGUID, "g0"):
		return a.GetGroupInfo(ctx, objectGUID)
	case strings.HasPrefix(objectGUID, "b0"):
		return a.GetBotInfo(ctx, objectGUID)
	case strings.HasPrefix(objectGUID, "s0"):
		return a.GetServiceInfo(ctx, objectGUID)
	case strings.HasPrefix(objectGUID, "u0"):
		return a.GetUserInfo(ctx, objectGUID)
	default:
		// Otherwise treat as a username.
		return a.GetObjectByUsername(ctx, objectGUID)
	}
}

// Join joins an object. If the input looks like a link, it dispatches
// to `joinGroup` / `joinChannelByLink`. Otherwise (GUID), it dispatches
// to `joinChannelAction(channelGUID, "Join")`.
func (a *Account) Join(ctx context.Context, input string) (*events.Event, error) {
	if strings.HasPrefix(input, "http") || strings.Contains(input, "/") {
		switch events.GUIDType(input) {
		case "Channel":
			return a.JoinChannelByLink(ctx, input)
		default:
			return a.JoinGroup(ctx, input)
		}
	}
	if strings.HasPrefix(input, "c0") {
		return a.JoinChannelAction(ctx, input, "Join")
	}
	return a.JoinGroup(ctx, input)
}

// LeaveChat leaves a chat — it dispatches to `joinChannelAction(channel,
// "Remove")` (channel) or `leaveGroup` (group).
func (a *Account) LeaveChat(ctx context.Context, objectGUID string) (*events.Event, error) {
	if strings.HasPrefix(objectGUID, "c0") {
		return a.JoinChannelAction(ctx, objectGUID, "Remove")
	}
	if strings.HasPrefix(objectGUID, "g0") {
		return a.LeaveGroup(ctx, objectGUID)
	}
	return nil, nil
}

// BanMember bans a member from a chat — it dispatches based on GUID
// prefix to `banGroupMember` / `banChannelMember`.
func (a *Account) BanMember(ctx context.Context, objectGUID, memberGUID string) (*events.Event, error) {
	if strings.HasPrefix(objectGUID, "c0") {
		return a.BanChannelMember(ctx, objectGUID, memberGUID)
	}
	return a.BanGroupMember(ctx, objectGUID, memberGUID)
}

// UnbanMember unbans a member from a chat.
func (a *Account) UnbanMember(ctx context.Context, objectGUID, memberGUID string) (*events.Event, error) {
	if strings.HasPrefix(objectGUID, "c0") {
		return a.UnbanChannelMember(ctx, objectGUID, memberGUID)
	}
	return a.UnbanGroupMember(ctx, objectGUID, memberGUID)
}

// GetMembers returns the full member list of a chat (group/channel).
func (a *Account) GetMembers(ctx context.Context, objectGUID string) ([]map[string]any, error) {
	var out []map[string]any
	startID := ""
	for {
		var res *events.Event
		var err error
		switch events.GUIDType(objectGUID) {
		case "Channel":
			res, err = a.GetChannelAllMembers(ctx, objectGUID, startID, "", 100)
		case "Group":
			res, err = a.GetGroupAllMembers(ctx, objectGUID, startID, "")
		default:
			return nil, fmt.Errorf("account: cannot list members for user chats")
		}
		if err != nil {
			return nil, err
		}
		for _, item := range res.Slice("in_chat") {
			if m, ok := item.(map[string]any); ok {
				out = append(out, m)
			}
		}
		next := res.Str("next_start_id")
		if next == "" || next == startID {
			break
		}
		startID = next
	}
	return out, nil
}
