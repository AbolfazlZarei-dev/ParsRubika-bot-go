package account

import (
	"context"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// Group management: info, members, admins, links and restrictions.

// EditSlowMode sets the slow-mode interval for a group. `slowMode`
// is in seconds (0 to disable).
func (a *Account) EditSlowMode(ctx context.Context, groupGUID string, slowMode int) (*events.Event, error) {
	return a.EditGroupInfoFull(ctx, groupGUID, "", "", map[string]any{"slow_mode": slowMode}, []string{"slow_mode"})
}

// EditGroupTimer is an alias for EditSlowMode.
func (a *Account) EditGroupTimer(ctx context.Context, groupGUID string, timer int) (*events.Event, error) {
	return a.EditSlowMode(ctx, groupGUID, timer)
}

// LockGroup removes SendMessages from the default access list.
func (a *Account) LockGroup(ctx context.Context, groupGUID string) (*events.Event, error) {
	res, err := a.GetGroupDefaultAccess(ctx, groupGUID)
	if err != nil {
		return nil, err
	}
	access := res.Slice("access_list")
	var newAccess []string
	for _, item := range access {
		if s, ok := item.(string); ok && s != "SendMessages" {
			newAccess = append(newAccess, s)
		}
	}
	return a.SetGroupDefaultAccess(ctx, groupGUID, newAccess)
}

// UnlockGroup adds SendMessages to the default access list.
func (a *Account) UnlockGroup(ctx context.Context, groupGUID string) (*events.Event, error) {
	res, err := a.GetGroupDefaultAccess(ctx, groupGUID)
	if err != nil {
		return nil, err
	}
	access := res.Slice("access_list")
	var newAccess []string
	hasSend := false
	for _, item := range access {
		if s, ok := item.(string); ok {
			newAccess = append(newAccess, s)
			if s == "SendMessages" {
				hasSend = true
			}
		}
	}
	if !hasSend {
		newAccess = append(newAccess, "SendMessages")
	}
	return a.SetGroupDefaultAccess(ctx, groupGUID, newAccess)
}

// GetGroupMessageReadParticipants returns who has read a specific message.
func (a *Account) GetGroupMessageReadParticipants(ctx context.Context, groupGUID, messageID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getGroupMessageReadParticipants", Input: map[string]any{
		"group_guid": groupGUID,
		"message_id": messageID,
	}, Encrypt: true})
}

// EditGroupInfoFull is the full-featured version of EditGroupInfo that
// accepts extra parameters (slow_mode, event_messages, reactions, etc.).
func (a *Account) EditGroupInfoFull(ctx context.Context, groupGUID, title, description string, extraValues map[string]any, updatedParameters []string) (*events.Event, error) {
	in := map[string]any{"group_guid": groupGUID}
	if title != "" {
		in["title"] = title
	}
	if description != "" {
		in["description"] = description
	}
	for k, v := range extraValues {
		in[k] = v
	}
	in["updated_parameters"] = updatedParameters
	return a.Request(ctx, RequestOptions{Method: "editGroupInfo", Input: in, Encrypt: true})
}

// GetGroupInfo fetches information about a group.
func (a *Account) GetGroupInfo(ctx context.Context, groupGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getGroupInfo", Input: map[string]any{"group_guid": groupGUID}, Encrypt: true})
}

// AddGroup creates a new group.
func (a *Account) AddGroup(ctx context.Context, title, description string, memberGUIDs []string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "addGroup",
		Input:   map[string]any{"title": title, "description": description, "member_guids": memberGUIDs},
		Encrypt: true,
	})
}

// RemoveGroup deletes a group owned by the current user.
func (a *Account) RemoveGroup(ctx context.Context, groupGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "removeGroup", Input: map[string]any{"group_guid": groupGUID}, Encrypt: true})
}

// EditGroupInfo edits a group's metadata. Pass an empty string for any
// field you do not want to change. The `updatedParameters` list is built
// automatically from the non-empty fields.
func (a *Account) EditGroupInfo(ctx context.Context, groupGUID, title, description string) (*events.Event, error) {
	in := map[string]any{"group_guid": groupGUID}
	updated := []string{}
	if title != "" {
		in["title"] = title
		updated = append(updated, "title")
	}
	if description != "" {
		in["description"] = description
		updated = append(updated, "description")
	}
	in["updated_parameters"] = updated
	return a.Request(ctx, RequestOptions{Method: "editGroupInfo", Input: in, Encrypt: true})
}

// EditGroupTitle is a convenience wrapper that only updates the title.
func (a *Account) EditGroupTitle(ctx context.Context, groupGUID, title string) (*events.Event, error) {
	return a.EditGroupInfo(ctx, groupGUID, title, "")
}

// EditGroupDescription is a convenience wrapper for the description only.
func (a *Account) EditGroupDescription(ctx context.Context, groupGUID, description string) (*events.Event, error) {
	return a.EditGroupInfo(ctx, groupGUID, "", description)
}

// JoinGroup joins a group via its invite link. The link is normalised
// to just the trailing slug before being sent as `hash_link`.
func (a *Account) JoinGroup(ctx context.Context, link string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "joinGroup",
		Input:   map[string]any{"hash_link": trailingSlug(link)},
		Encrypt: true,
	})
}

// LeaveGroup leaves a group.
func (a *Account) LeaveGroup(ctx context.Context, groupGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "leaveGroup", Input: map[string]any{"group_guid": groupGUID}, Encrypt: true})
}

// AddGroupMembers adds members to a group.
func (a *Account) AddGroupMembers(ctx context.Context, groupGUID string, memberGUIDs []string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "addGroupMembers",
		Input:   map[string]any{"group_guid": groupGUID, "member_guids": memberGUIDs},
		Encrypt: true,
	})
}

// BanGroupMember bans a group member (action = "Set").
func (a *Account) BanGroupMember(ctx context.Context, groupGUID, userGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "banGroupMember",
		Input:   map[string]any{"group_guid": groupGUID, "member_guid": userGUID, "action": "Set"},
		Encrypt: true,
	})
}

// UnbanGroupMember unbans a group member (action = "Unset").
func (a *Account) UnbanGroupMember(ctx context.Context, groupGUID, userGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "banGroupMember",
		Input:   map[string]any{"group_guid": groupGUID, "member_guid": userGUID, "action": "Unset"},
		Encrypt: true,
	})
}

// GetGroupAllMembers returns the group member list with an optional
// search-text filter and a pagination start-id.
func (a *Account) GetGroupAllMembers(ctx context.Context, groupGUID, startID, searchText string) (*events.Event, error) {
	in := map[string]any{"group_guid": groupGUID}
	if startID != "" {
		in["start_id"] = startID
	}
	if searchText != "" {
		in["search_text"] = searchText
	}
	return a.Request(ctx, RequestOptions{Method: "getGroupAllMembers", Input: in, Encrypt: true})
}

// GetGroupAdminMembers returns the admin member list (paged).
func (a *Account) GetGroupAdminMembers(ctx context.Context, groupGUID, startID string) (*events.Event, error) {
	in := map[string]any{"group_guid": groupGUID}
	if startID != "" {
		in["start_id"] = startID
	}
	return a.Request(ctx, RequestOptions{Method: "getGroupAdminMembers", Input: in, Encrypt: true})
}

// GetGroupAdminAccessList returns the access rights of an admin.
func (a *Account) GetGroupAdminAccessList(ctx context.Context, groupGUID, memberGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "getGroupAdminAccessList",
		Input:   map[string]any{"group_guid": groupGUID, "member_guid": memberGUID},
		Encrypt: true,
	})
}

// SetGroupAdmin promotes/demotes a group admin. `action` is "SetAdmin"
// or "UnsetAdmin". `accessList` is the list of permissions (e.g.
// ["BanMember", "PinMessages"]).
func (a *Account) SetGroupAdmin(ctx context.Context, groupGUID, userGUID, action string, accessList []string) (*events.Event, error) {
	in := map[string]any{
		"group_guid":  groupGUID,
		"member_guid": userGUID,
		"action":      action,
	}
	if len(accessList) > 0 {
		in["access_list"] = accessList
	}
	return a.Request(ctx, RequestOptions{Method: "setGroupAdmin", Input: in, Encrypt: true})
}

// PromoteGroupAdmin is a convenience wrapper that promotes a member.
func (a *Account) PromoteGroupAdmin(ctx context.Context, groupGUID, userGUID string, accessList []string) (*events.Event, error) {
	return a.SetGroupAdmin(ctx, groupGUID, userGUID, "SetAdmin", accessList)
}

// DemoteGroupAdmin is a convenience wrapper that demotes a member.
func (a *Account) DemoteGroupAdmin(ctx context.Context, groupGUID, userGUID string) (*events.Event, error) {
	return a.SetGroupAdmin(ctx, groupGUID, userGUID, "UnsetAdmin", nil)
}

// GetGroupDefaultAccess returns the default access rights of a group.
func (a *Account) GetGroupDefaultAccess(ctx context.Context, groupGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getGroupDefaultAccess", Input: map[string]any{"group_guid": groupGUID}, Encrypt: true})
}

// SetGroupDefaultAccess updates the default access rights of a group.
func (a *Account) SetGroupDefaultAccess(ctx context.Context, groupGUID string, accessList []string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:  "setGroupDefaultAccess",
		Input:   map[string]any{"group_guid": groupGUID, "access_list": accessList},
		Encrypt: true,
	})
}

// GetGroupLink returns the group's invite link.
func (a *Account) GetGroupLink(ctx context.Context, groupGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getGroupLink", Input: map[string]any{"group_guid": groupGUID}, Encrypt: true})
}

// SetGroupLink regenerates the group's invite link.
func (a *Account) SetGroupLink(ctx context.Context, groupGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "setGroupLink", Input: map[string]any{"group_guid": groupGUID}, Encrypt: true})
}

// GroupPreviewByJoinLink fetches the preview before joining via a link.
// The link is normalised to its trailing slug before being sent as
// `hash_link`.
func (a *Account) GroupPreviewByJoinLink(ctx context.Context, link string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:     "groupPreviewByJoinLink",
		Input:      map[string]any{"hash_link": trailingSlug(link)},
		TmpSession: true,
		Encrypt:    true,
	})
}

// GetGroupOnlineCount returns the number of online members.
func (a *Account) GetGroupOnlineCount(ctx context.Context, groupGUID string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{Method: "getGroupOnlineCount", Input: map[string]any{"group_guid": groupGUID}, Encrypt: true})
}

// GetGroupMentionList returns the mention list of a group, optionally
// filtered by `searchMention`.
func (a *Account) GetGroupMentionList(ctx context.Context, groupGUID, searchMention string) (*events.Event, error) {
	in := map[string]any{"group_guid": groupGUID}
	if searchMention != "" {
		in["search_mention"] = searchMention
	}
	return a.Request(ctx, RequestOptions{Method: "getGroupMentionList", Input: in, Encrypt: true})
}

// GetBannedGroupMembers returns the list of banned members (paged).
func (a *Account) GetBannedGroupMembers(ctx context.Context, groupGUID, startID string) (*events.Event, error) {
	in := map[string]any{"group_guid": groupGUID}
	if startID != "" {
		in["start_id"] = startID
	}
	return a.Request(ctx, RequestOptions{Method: "getBannedGroupMembers", Input: in, Encrypt: true})
}
